package convbinder

import (
	"sdim_pc/backend/api/convapi"
	"sdim_pc/backend/chat"
	"sdim_pc/backend/mylog"
	"sdim_pc/backend/user"
)

type ConvBinder struct {
	ca *convapi.ConvApi
	cm *chat.ConvManager
}

func NewConvBinder(ca *convapi.ConvApi, cm *chat.ConvManager) *ConvBinder {
	return &ConvBinder{
		ca: ca,
		cm: cm,
	}
}

func (b *ConvBinder) ClearUnreadCount(convId string) error {
	if !b.cm.ShouldClearUnread(convId) {
		return nil
	}

	mylog.GetLogger().Trace().Msgf("conv clear unread, convId=%s", convId)

	uid := user.GetUid()
	err := b.ca.ClearUnread(convId, uid)
	if err != nil {
		return err
	}

	items, idx, ok := b.cm.UpdateAfterClearUnread(convId)

	if ok {
		b.cm.EmitConvListUpdateEvent(items, idx)
	}

	return nil
}
