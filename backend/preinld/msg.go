package preinld

type MsgSendData struct {
	Sender     string      `json:"sender,omitempty"`   // 发送者uid
	ConvId     string      `json:"convId,omitempty"`   // 发送的会话id
	Receiver   string      `json:"receiver,omitempty"` // 接收者, 单聊是对方的uid, 群聊是群id
	ChatType   ChatType    `json:"chatType,omitempty"`
	Ttl        int32       `json:"ttl,omitempty"` // 消息过期时间(sec), -1:阅后即焚,0:不过期
	MsgContent *MsgContent `json:"msgContent,omitempty"`
	ClientId   string      `json:"clientId,omitempty"`
}

type MsgType uint16

const (
	// chat
	// text
	TextType MsgType = 1

	// image
	ImageType MsgType = 2

	// custom
	CustomType MsgType = 100

	// cmd
	CmdType MsgType = 1000
)

type SubCmdType uint16

const (
	// 邀请入群
	SubCmdGroupInvited SubCmdType = 1001
)

const SysSendUser = "sys:send:sys_auto"

type SenderType uint8

const (
	UserSender SenderType = 1

	SysCmdSender SenderType = 10
)

type (
	MsgContent struct {
		Type    MsgType        `json:"type,omitempty"`    // 消息类型
		Content map[string]any `json:"content,omitempty"` // 消息内容
		Custom  map[string]any `json:"custom,omitempty"`  // 自定义内容
		Extra   map[string]any `json:"extra,omitempty"`   // extra内容
	}

	SenderInfo struct {
		SenderType SenderType `json:"senderType"`
		Nickname   string     `json:"nickname,omitempty"`
		Avatar     string     `json:"avatar,omitempty"`
	}

	Msg struct {
		MsgId          int64       `json:"msgId"`
		ConvId         string      `json:"convId"`
		Sender         string      `json:"sender"`
		Receiver       string      `json:"receiver"`
		ChatType       ChatType    `json:"chatType"`
		MsgType        MsgType     `json:"msgType"`
		Content        *MsgContent `json:"content"`
		SenderInfo     SenderInfo  `json:"senderInfo"`
		MegSeq         int64       `json:"megSeq"`
		Cts            int64       `json:"cts"`
		State          uint8       `json:"state"` // 1=发送中,2=发送成功,3=发送失败
		LastFailReason string      `json:"lastFailReason"`
		RetryTimes     uint8       `json:"retryTimes"`
		IsSelf         bool        `json:"isSelf"`
		ClientId       string      `json:"clientId"`
	}
)
