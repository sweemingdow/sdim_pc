package preinld

import (
	"time"
)

const (
	HeaderFrameSize = 12 // 2+1+1+4+4
	ReqIdSize       = 16
	MagicNumber     = 0xCAFE
	Version         = 1
	MaxPayloadSize  = 1 << 20 // 1MB

	// 心跳间隔：根据服务器要求调整，这里设为 15 秒
	HeartbeatInterval = 120 * time.Second

	// 重连最大重试次数
	ReconnectMaxRetryTimes = 10

	// 重连最大重试间隔
	ReconnectMaxRetryInterval = 10 * time.Second

	// 重连重试延时间隔
	ReconnectDelayInterval = 200 * time.Millisecond
)

// 客户端类型
const (
	Pc uint8 = 3
)

type ChatType uint8

const (
	P2pChat      ChatType = 1
	GroupChat    ChatType = 2
	CustomerChat ChatType = 3
	Danmaku      ChatType = 4
)

type (
	ConnFrame struct {
		Uid     string `json:"uid,omitempty"`
		CType   uint8  `json:"ctype,omitempty"`
		TsMills int64  `json:"tsMills,omitempty"`
	}

	ConnAckFrame struct {
		ErrCode  ErrCode `json:"errCode,omitempty"`
		ErrDesc  string  `json:"errDesc,omitempty"`
		TimeDiff int64   `json:"timeDiff"` // 客户端和服务器之间的时差
		SignKey  string  `json:"signKey,omitempty"`
	}
)

type (
	ErrCode uint32

	SendFrame struct {
		Sender         string      `json:"sender,omitempty"`   // 发送者uid
		Receiver       string      `json:"receiver,omitempty"` // 接收者, 单聊是对方的uid, 群聊是群id
		ChatType       ChatType    `json:"chatType,omitempty"`
		SendMills      int64       `json:"sendMills,omitempty"`
		Sign           string      `json:"sign,omitempty"`           // 消息签名, 防纂改
		Ttl            int32       `json:"ttl,omitempty"`            // 消息过期时间(sec), -1:阅后即焚,0:不过期
		ClientUniqueId string      `json:"clientUniqueId,omitempty"` // 客户端唯一id
		MsgContent     *MsgContent `json:"msgContent,omitempty"`
	}

	SendFrameAck struct {
		ErrCode ErrCode          `json:"errCode,omitempty"`
		ErrDesc string           `json:"errDesc,omitempty"`
		Data    SendAckFrameBody `json:"data,omitempty"`
	}

	SendAckFrameBody struct {
		MsgId          int64  `json:"msgId,omitempty"`          // 消息id
		ClientUniqueId string `json:"clientUniqueId,omitempty"` // 客户端唯一id
		ConvId         string `json:"convId,omitempty"`         // 会话id
	}
)

const (
	OK         ErrCode = 0
	BizErr     ErrCode = 1000
	ServerErr  ErrCode = 2000
	RpcRespErr ErrCode = 3000
)

type ErrItem struct {
	ErrCode ErrCode
	Desc    string
}

func IsOk(code ErrCode) bool {
	return code == OK
}
