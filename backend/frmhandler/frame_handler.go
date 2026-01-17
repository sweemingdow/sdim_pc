package frmhandler

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"sdim_pc/backend/api/convapi"
	"sdim_pc/backend/chat"
	"sdim_pc/backend/client/frm"
	"sdim_pc/backend/mylog"
	"sdim_pc/backend/preinld"
	"sdim_pc/backend/utils/parser/json"
	"sync/atomic"
)

type FrameHandler struct {
	frmCh  <-chan *frm.Frame
	stopCh chan struct{}
	cm     *chat.ConvManager
	pool   *ants.Pool
	closed atomic.Bool
	ca     *convapi.ConvApi
}

func NewFrameHandler(
	frmCh <-chan *frm.Frame,
	cm *chat.ConvManager,
	ca *convapi.ConvApi,
) *FrameHandler {
	pool, _ := ants.NewPool(
		8,
		ants.WithMaxBlockingTasks(128),
		ants.WithPreAlloc(true),
		ants.WithNonblocking(true),
	)

	fh := &FrameHandler{
		frmCh:  frmCh,
		stopCh: make(chan struct{}),
		cm:     cm,
		ca:     ca,
		pool:   pool,
	}

	go fh.receiveFrame()

	return fh
}

func (fh *FrameHandler) StopReceive(_ context.Context) {
	if !fh.closed.CompareAndSwap(false, true) {
		return
	}
	close(fh.stopCh)
}

func (fh *FrameHandler) Cleanup(ctx context.Context, frames []*frm.Frame) {
	for _, frame := range frames {
		fh.handleFrame(frame)
	}
}

func (fh *FrameHandler) receiveFrame() {
	for {
		select {
		case <-fh.stopCh:
			return
		case frame, ok := <-fh.frmCh:
			if !ok {
				return
			}

			mylog.GetLogger().Trace().Msgf("reveive frame data, frameDesc:%s, frameData:%+v", frm.FrameType2desc(frame.Header.Ftype), frame)

			fh.submitTask(func() {
				fh.handleFrame(frame)
			})
		}
	}
}

func (fh *FrameHandler) handleFrame(frame *frm.Frame) {
	lg := mylog.GetLogger()

	ft := frame.Header.Ftype
	// 发送的回帧
	if ft == frm.SendAck {
		var saf preinld.SendFrameAck
		err := json.Parse(frame.Payload.Body, &saf)
		if err != nil {
			lg.Error().Stack().Err(err).Msg("parse send frame ack body failed")
			return
		}

		lg.Trace().Msgf("parse send frame ack body success, data=%+v", saf.Data)

		// 模拟发送耗时
		//time.Sleep(2500 * time.Millisecond)

		if !preinld.IsOk(saf.ErrCode) {
			lg.Warn().Msgf("send frame ack errCode is not ok, errCode=%d, errDesc:%s", saf.ErrCode, saf.ErrDesc)
			items, idx, ok := fh.cm.UpdateMsgWhenSentFailed(saf)
			if ok {
				fh.cm.EmitConvListUpdateEvent(items, idx)
			}
			return
		}

		items, idx, ok := fh.cm.UpdateMsgWhenSentSuccess(saf.Data)
		if ok {
			fh.cm.EmitConvListUpdateEvent(items, idx)
		}
	} else if ft == frm.Forward { // 转发给自己的消息帧
		var ffb preinld.ForwardFrameBody
		err := json.Parse(frame.Payload.Body, &ffb)
		if err != nil {
			lg.Error().Stack().Err(err).Msg("parse forward frame body failed")
			return
		}

		// todo 接受到转发消息回帧给服务端

		// 修改会话 & 消息
		items, idx, ok := fh.cm.InsertMsgAfterReceived(&ffb)
		if ok {
			fh.cm.EmitConvListUpdateEvent(items, idx)
		}
	} else if ft == frm.ConvUpdate {
		var cuf preinld.ConvUpdateFrame
		err := json.Parse(frame.Payload.Body, &cuf)
		if err != nil {
			lg.Error().Stack().Err(err).Msg("parse conv update frame failed")
			return
		}

		// 修改会话 & 消息
		items, idx, ok := fh.cm.UpdateWhenConvUpdate(&cuf)
		if ok {
			fh.cm.EmitConvListUpdateEvent(items, idx)
		}
	}
}

func (fh *FrameHandler) submitTask(run func()) {
	if err := fh.pool.Submit(run); err == ants.ErrPoolOverload {
		mylog.GetLogger().Error().Err(err).Msg("too many tasks, check pls")
	}
}
