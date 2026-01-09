package frm

// 数据帧类型
type FrameType uint8

const (
	Ping       FrameType = 1 // ping(c2s)
	Pong       FrameType = 2 // pong(s2c)
	Conn       FrameType = 3 // 连接(c2s)
	ConnAck    FrameType = 4 // 连接确认(s2c)
	Send       FrameType = 5 // 发送消息(c2s)
	SendAck    FrameType = 6 // 发送消息确认(s2c)
	Forward    FrameType = 7 // 转发消息(s2c)
	ForwardAck FrameType = 8 // 收到转发消息确认(c2s)
	Disconnect FrameType = 9 // 断开连接(服务端主动发起)
)

// 荷载协议类型(json, protobuf, thrift...)
type PayloadProtocolType uint8

const (
	JsonPayload     PayloadProtocolType = 1
	MsgpackPayload  PayloadProtocolType = 2
	ProtobufPayload PayloadProtocolType = 3
	ThriftPayload   PayloadProtocolType = 4
)

const (
	HeaderFrameSize = 2 + 1 + 1 + 4 + 4 // 帧头数据长度
	MaxPayloadSize  = 1 << 20           // 荷载最大长度 1MB
	MagicNumber     = 0xCAFE            // 模数验证
	Version         = 1                 // 版本
	ReqIdSize       = 16                // reqId大小
)

type (
	Frame struct {
		Header  FrameHeader // 帧头
		Payload Payload     // 帧荷载
	}

	FrameHeader struct {
		Magic      uint16    // 模数 2b
		Version    uint8     // 版本 1b
		Ftype      FrameType // 帧类型 1b
		PayloadLen uint32    // 荷载长度 4b
		CheckSum   uint32    // 校验和 4b
	}

	Payload struct {
		PayloadProtocol PayloadProtocolType // 荷载协议 1b
		ReqId           [ReqIdSize]byte     // 请求id(标准的uuid) 16b
		Body            []byte              // 数据 nb
	}
)

var ft2desc = map[FrameType]string{
	Ping:       "Ping",
	Pong:       "Pong",
	Conn:       "Conn",
	ConnAck:    "ConnAck",
	Send:       "Send",
	SendAck:    "SendAck",
	Forward:    "Forward",
	ForwardAck: "ForwardAck",
	Disconnect: "Disconnect",
}

func FrameType2desc(ft FrameType) string {
	return ft2desc[ft]
}
