package msgbinder

import (
	"errors"
	"sdim_pc/backend/api/msgapi"
	"sdim_pc/backend/chat"
	"sdim_pc/backend/mylog"
)

type MsgBinder struct {
	ma *msgapi.MsgApi
	cm *chat.ConvManager
}

func NewMsgBinder(ma *msgapi.MsgApi, cm *chat.ConvManager) *MsgBinder {
	return &MsgBinder{
		ma: ma,
		cm: cm,
	}
}

func (b *MsgBinder) FetchNextMsgs(convId string) error {
	lastMsgId, ok := b.cm.TailMsgId(convId)
	if !ok {
		return errors.New("can not found conv with id:" + convId)
	}

	mylog.GetLogger().Debug().Msgf("convId=%s fetch next msgs start, lastMsgId=%d", convId, lastMsgId)

	msgs, err := b.ma.FetchNextMsgs(convId, lastMsgId)

	if err != nil {
		return err
	}

	// 更新会话中的最近消息
	items, idx, exists := b.cm.AppendTailMsgs(convId, msgs)

	if exists {
		b.cm.EmitConvListUpdateEvent(items, idx)
	}

	return nil
}
