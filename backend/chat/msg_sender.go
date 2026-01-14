package chat

type MsgSendState uint8

const (
	Sending    MsgSendState = 1
	SendOk     MsgSendState = 2
	SendFailed MsgSendState = 3
)

/*
type msgSendState struct {
	msd        *preinld.MsgSendData
	state      MsgSendState // 1=发送中,2=发送成功,3=发送失败
	retryTimes uint8
}

type MsgSender struct {
	clientId2mss map[string]*msgSendState
	cli          *client.Client
	cm           *ConvManager
	mu           sync.Mutex
}

func NewMsgSender(cli *client.Client, cm *ConvManager) *MsgSender {
	return &MsgSender{
		clientId2mss: make(map[string]*msgSendState, 200),
		cli:          cli,
		cm:           cm,
	}
}

func (ms *MsgSender) SendMsg(msd *preinld.MsgSendData) error {
	clientId := msd.ClientId
	var retried = false
	if clientId == "" {
		clientId = utils.RandStr(32)
	} else {
		retried = true
	}

	ms.mu.Lock()

	if retried {
		// 重试发送
		mss, ok := ms.clientId2mss[clientId]
		if !ok {
			return fmt.Errorf("can not found msg by clientId:%s, when send retry", clientId)
		}

		if mss.state == SendOk {
			return nil
		}

		mss.state = Sending
		mss.retryTimes += 1
	} else {
		// 新消息发送
		ms.clientId2mss[clientId] = &msgSendState{
			msd:        msd,
			state:      Sending,
			retryTimes: 0,
		}
	}

	ms.mu.Unlock()

	return ms.cli.SendMsgFrame(msd)
}

func (ms *MsgSender) ModifyAfterSendSuccess(ackBody preinld.SendAckFrameBody) (*preinld.MsgSendData, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	mss, ok := ms.clientId2mss[ackBody.ClientUniqueId]
	if !ok {
		return nil, false
	}

	if mss.state == Sending {
		mss.state = SendOk
		return mss.msd, ok
	}

	return nil, false
}*/
