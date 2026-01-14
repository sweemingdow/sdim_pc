package client

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"sdim_pc/backend/client/frm"
	"sdim_pc/backend/mylog"
	preinld "sdim_pc/backend/preinld"
	"sdim_pc/backend/user"
	"sdim_pc/backend/utils"
	"sync"
	"sync/atomic"
	"time"
)

var (
	NotConnectedErr = errors.New("not connected")
	WasClosedErr    = errors.New("client was closed")
)

type ClientEventCallback interface {
	// 连接中
	OnConnecting()

	// 连接成功
	OnConnected()

	// 连接失败
	OnConnectFailed(retryTimes int)

	// 服务端主动断连
	OnDisconnected(errCode string, reason string)
}

// Client 客户端
type Client struct {
	conn        net.Conn
	addr        string
	connected   atomic.Bool
	connecting  atomic.Bool
	closed      atomic.Bool
	pongTimed   time.Time
	ctime       time.Time
	utime       time.Time
	mu          sync.Mutex
	stopChan    chan struct{}
	interrupted atomic.Bool
	activeDis   atomic.Bool // 主动断连
	frameChan   chan *frm.Frame
	eventCb     ClientEventCallback
}

func NewClient(addr string, eventCb ClientEventCallback) (*Client, error) {
	return &Client{
		addr:      addr,
		stopChan:  make(chan struct{}),
		frameChan: make(chan *frm.Frame, 100),
		eventCb:   eventCb,
	}, nil
}

// Connect 连接到服务器，并在1秒内等待ConnAck，否则超时断连
func (c *Client) Connect(uid string, cType uint8) error {
	c.mu.Lock()
	conn := c.conn
	if conn != nil {
		_ = conn.Close()
	}
	user.Reset()
	c.mu.Unlock()

	c.connecting.Store(true)

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.connected.Store(true)
	c.mu.Unlock()

	// 发送 conn 帧
	err = c.sendConnFrame(uid, cType)
	if err != nil {
		_ = c.conn.Close()
		return fmt.Errorf("send conn frame failed: %w", err)
	}

	// 设置读取超时为1秒
	err = c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		_ = c.conn.Close()
		return fmt.Errorf("send conn frame failed: %w", err)
	}

	// 读取 conn_ack 帧（同步读，带超时）
	frame, err := c.readFrame()

	if err != nil {
		_ = c.conn.Close()
		// 判断是否是超时错误
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("timeout waiting for ConnAck (1s)")
		}
		return fmt.Errorf("read conn ack failed: %w", err)
	}

	// 恢复读取 deadline（可选，因后续由 readLoop 控制）
	err = c.conn.SetReadDeadline(time.Time{})

	if frame.Header.Ftype != frm.ConnAck {
		_ = c.conn.Close()
		return fmt.Errorf("unexpected frame type: %d, expected ConnAck", frame.Header.Ftype)
	}

	// 解析 conn_ack
	var caf preinld.ConnAckFrame
	err = json.Unmarshal(frame.Payload.Body, &caf)
	if err != nil {
		_ = c.conn.Close()
		return fmt.Errorf("parse conn ack body failed: %w", err)
	}

	if !preinld.IsOk(caf.ErrCode) {
		return fmt.Errorf("conn ack frame not ok, errCode=%d, errDesc=%s", caf.ErrCode, caf.ErrDesc)
	}

	c.mu.Lock()
	c.conn = conn
	c.activeDis.Store(false)
	c.connecting.Store(false)
	c.ctime = time.Now()

	user.Replace(user.UserInfo{
		Uid:        uid,
		SignKey:    caf.SignKey,
		ClientType: cType,
		CTime:      time.Now(),
	})

	c.mu.Unlock()

	c.loopStart()

	return nil
}

func (c *Client) GetFrameChan() <-chan *frm.Frame {
	return c.frameChan
}

func (c *Client) Disconnect() error {
	c.activeDis.Store(true)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.connected.Store(false)
	err := c.conn.Close()
	c.conn = nil

	return err
}

func (c *Client) Stop(ctx context.Context) ([]*frm.Frame, error) {
	if !c.closed.CompareAndSwap(false, true) {
		return nil, nil
	}

	c.connected.Store(false)
	c.interrupted.Store(true)

	close(c.stopChan)

	done := make(chan []*frm.Frame)
	go func() {
		c.mu.Lock()
		conn := c.conn
		if conn != nil {
			_ = conn.Close()
		}
		c.mu.Unlock()

		close(c.frameChan)

		var frames []*frm.Frame
		// 遍历未处理完的frames
		for frame := range c.frameChan {
			frames = append(frames, frame)
		}

		done <- frames
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case frames := <-done:
		return frames, nil
	}
}

func (c *Client) SendMsgFrame(msd preinld.MsgSendData) error {
	if err := c.nextIfCan(); err != nil {
		return err
	}

	sender := user.GetUserInfo().Uid

	mylog.GetLogger().Debug().Msgf("sender=%s, receiver=%s, chatType=%d, msgContent:%+v, start send msg frame",
		sender, msd.Receiver, msd.ChatType, *msd.MsgContent)

	sf := preinld.SendFrame{
		Sender:         sender,
		Receiver:       msd.Receiver,
		ChatType:       msd.ChatType,
		SendMills:      time.Now().UnixMilli(),
		Sign:           "",
		Ttl:            msd.Ttl,
		ClientUniqueId: msd.ClientId,
		MsgContent:     msd.MsgContent,
	}

	bodies, err := json.Marshal(sf)
	if err != nil {
		return err
	}

	reqId := utils.Uuid()

	frame := c.buildFrame(reqId, bodies, frm.Send)

	frmBuf, err := c.encodeFrame(frame)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(frmBuf)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) loopStart() {
	go c.readLoop()

	go c.heartbeatLoop()
}

func (c *Client) readLoop() {
	for !c.interrupted.Load() {
		select {
		case <-c.stopChan:
			return
		default:

		}

		if !c.connected.Load() {
			return
		}

		lg := mylog.GetLogger()

		frame, err := c.readFrame()

		if err != nil {
			lg.Error().Stack().Err(err).Msg("read frame failed, will reconnecting")

			che, ok := decodeClientHandleError(err)

			if ok {
				if che.ShouldReconnect() {
					// 发送ping失败触发重连
					c.doReconnect()
				}
			}

			return
		}

		if c.interrupted.Load() {
			return
		}

		// 服务端断开连接
		if frame.Header.Ftype == frm.Disconnect {
			c.connected.Store(false)
		} else if frame.Header.Ftype == frm.Pong {
			c.mu.Lock()
			c.pongTimed = time.Now()
			c.mu.Unlock()
		}

		select {
		case <-c.stopChan:
			return
		case c.frameChan <- frame: // 写入frame处理chan中
		}
	}
}

func (c *Client) heartbeatLoop() {
	ticker := time.NewTicker(preinld.HeartbeatInterval)
	defer ticker.Stop()

	for !c.interrupted.Load() {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			if !c.connected.Load() {
				return
			}

			if c.interrupted.Load() {
				return
			}

			if err := c.sendPingFrame(); err != nil {
				mylog.GetLogger().Error().Stack().Err(err).Msg("send ping frame failed, will reconnecting")

				che, ok := decodeClientHandleError(err)

				if ok {
					if che.ShouldReconnect() {
						// 发送ping失败触发重连
						c.doReconnect()
					}
				}
			}
		}
	}
}

func (c *Client) buildFrame(reqId [frm.ReqIdSize]byte, bodies []byte, ft frm.FrameType) *frm.Frame {
	return &frm.Frame{
		Header: frm.FrameHeader{
			Magic:      frm.MagicNumber,
			Version:    frm.Version,
			Ftype:      ft,
			PayloadLen: uint32(1 + frm.ReqIdSize + len(bodies)),
			CheckSum:   0,
		},
		Payload: frm.Payload{
			PayloadProtocol: frm.JsonPayload,
			ReqId:           reqId,
			Body:            bodies,
		},
	}
}

// 编码帧
func (c *Client) encodeFrame(frame *frm.Frame) ([]byte, error) {
	// payload 数据长度
	pdLen := 1 + frm.ReqIdSize + len(frame.Payload.Body)

	if pdLen > frm.MaxPayloadSize {
		return nil, fmt.Errorf("frame too large")
	}

	pdBuf := make([]byte, pdLen)

	pdBuf[0] = byte(frame.Payload.PayloadProtocol)

	copy(pdBuf[1:1+frm.ReqIdSize], frame.Payload.ReqId[:])
	if len(frame.Payload.Body) > 0 {
		copy(pdBuf[1+frm.ReqIdSize:], frame.Payload.Body)
	}

	checksum := crc32.ChecksumIEEE(pdBuf)

	totalLen := frm.HeaderFrameSize + pdLen
	frameBuf := make([]byte, totalLen)
	binary.BigEndian.PutUint16(frameBuf[0:2], frm.MagicNumber)
	frameBuf[2] = frm.Version
	frameBuf[3] = uint8(frame.Header.Ftype)
	binary.BigEndian.PutUint32(frameBuf[4:8], uint32(pdLen))
	binary.BigEndian.PutUint32(frameBuf[8:12], checksum)
	copy(frameBuf[frm.HeaderFrameSize:], pdBuf)

	return frameBuf, nil
}

func (c *Client) readFrame() (*frm.Frame, error) {
	var frame frm.Frame
	headerBuf := make([]byte, frm.HeaderFrameSize)

	read := 0
	for read < frm.HeaderFrameSize {
		n, err := c.conn.Read(headerBuf[read:])
		if err != nil {
			return nil, newClientHandleError(fmt.Errorf("read header failed at byte %d: %w", read, err), err != io.EOF)
		}
		read += n

		if err = c.nextIfCan(); err != nil {
			return nil, newClientHandleError(err, false)
		}
	}

	frame.Header.Magic = binary.BigEndian.Uint16(headerBuf[0:2])
	frame.Header.Version = headerBuf[2]
	frame.Header.Ftype = frm.FrameType(headerBuf[3])
	frame.Header.PayloadLen = binary.BigEndian.Uint32(headerBuf[4:8])
	frame.Header.CheckSum = binary.BigEndian.Uint32(headerBuf[8:12])

	if frame.Header.Magic != frm.MagicNumber {
		return nil, newClientHandleError(fmt.Errorf("invalid magic: 0x%04X", frame.Header.Magic), false)
	}
	if frame.Header.Version != frm.Version {
		return nil, newClientHandleError(fmt.Errorf("unsupported version: %d", frame.Header.Version), false)
	}

	if frame.Header.PayloadLen == 0 {
		return &frame, nil
	}
	if frame.Header.PayloadLen > frm.MaxPayloadSize {
		return nil, newClientHandleError(fmt.Errorf("payload too large: %d", frame.Header.PayloadLen), false)
	}

	pdBuf := make([]byte, frame.Header.PayloadLen)
	read = 0
	for read < int(frame.Header.PayloadLen) {
		n, err := c.conn.Read(pdBuf[read:])
		if err != nil {
			return nil, newClientHandleError(fmt.Errorf("read payload failed: %w", err), err != io.EOF)
		}
		read += n

		if err = c.nextIfCan(); err != nil {
			return nil, newClientHandleError(err, false)
		}
	}

	checksum := crc32.ChecksumIEEE(pdBuf)
	if frame.Header.CheckSum != checksum {
		return nil, newClientHandleError(fmt.Errorf("checksum mismatch: got 0x%08X, want 0x%08X",
			checksum, frame.Header.CheckSum), false)
	}

	// 解析 payload：所有帧都包含 PayloadProtocol + ReqId + Body
	// 如果服务器对 Ack 帧也按此格式返回，则兼容
	if len(pdBuf) < 1+frm.ReqIdSize {
		return nil, newClientHandleError(fmt.Errorf("payload too short: %d", len(pdBuf)), false)
	}

	frame.Payload.PayloadProtocol = frm.PayloadProtocolType(pdBuf[0])
	copy(frame.Payload.ReqId[:], pdBuf[1:1+frm.ReqIdSize])
	if len(pdBuf) > 1+frm.ReqIdSize {
		frame.Payload.Body = pdBuf[1+frm.ReqIdSize:]
	}

	return &frame, nil
}

func (c *Client) nextIfCan() error {
	if c.isClosed() {
		return WasClosedErr
	}

	if !c.isConnected() {
		return NotConnectedErr
	}

	return nil
}

func (c *Client) isConnected() bool {
	return c.connected.Load()
}

func (c *Client) isClosed() bool {
	return c.closed.Load()
}

// 发送连接帧
func (c *Client) sendConnFrame(uid string, cType uint8) error {
	mylog.GetLogger().Debug().Msgf("uid=%s, cType=%d, start send conn frame", uid, cType)

	cf := preinld.ConnFrame{
		Uid:     uid,
		CType:   cType,
		TsMills: time.Now().UnixMilli(),
	}
	bodies, err := json.Marshal(cf)
	if err != nil {
		return err
	}

	reqId := utils.Uuid()

	scf := c.buildFrame(reqId, bodies, frm.Conn)

	frmBuf, err := c.encodeFrame(scf)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(frmBuf)

	return err
}

// 发送ping帧
func (c *Client) sendPingFrame() error {
	if err := c.nextIfCan(); err != nil {
		return newClientHandleError(err, false)
	}

	mylog.GetLogger().Trace().Msgf("uid=%s, start send ping frame", user.GetUid())

	frame := c.buildFrame(utils.Uuid(), []byte("ping"), frm.Ping)

	frmBuf, err := c.encodeFrame(frame)
	if err != nil {
		return newClientHandleError(err, false)
	}

	_, err = c.conn.Write(frmBuf)
	if err != nil {
		return newClientHandleError(err, err != io.EOF)
	}

	return nil
}

func (c *Client) doReconnect() {
	//if !c.connecting.CompareAndSwap(false, true) {
	//	return
	//}
	//
	//if c.activeDis.Load() {
	//	return
	//}
	//
	//c.interrupted.Store(true)
	//
	//c.connected.Store(false)
	//
	//if c.isClosed() {
	//	return
	//}
	//
	//c.mu.Lock()
	//cn := c.conn
	//c.mu.Unlock()
	//
	//if cn != nil {
	//	_ = cn.Close()
	//}
	//
	//go func() {
	//	defer func() {
	//		c.connecting.Store(false)
	//	}()
	//
	//	lg := mylog.GetLogger()
	//
	//	for times := 1; times <= preinld.ReconnectMaxRetryTimes; times++ {
	//		lg.Info().Msgf("client reconnecting, retryTimes=%d", times)
	//	}
	//}()
}
