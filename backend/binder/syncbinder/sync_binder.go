package syncbinder

import (
	"sdim_pc/backend/api/convapi"
	"sdim_pc/backend/chat"
	"sdim_pc/backend/mylog"
)

type SyncBinder struct {
	ca *convapi.ConvApi
	cm *chat.ConvManager
}

func NewSyncBinder(ca *convapi.ConvApi, cm *chat.ConvManager) *SyncBinder {
	return &SyncBinder{
		ca: ca,
		cm: cm,
	}
}

func (b *SyncBinder) SyncConvList(uid string) ([]*chat.ConvItem, error) {
	mylog.GetLogger().Debug().Msgf("uid=%s start sync conv list after connected success", uid)

	items, err := b.ca.SyncHotConvList(uid)
	if err != nil {
		return make([]*chat.ConvItem, 0), err
	}

	mylog.GetLogger().Debug().Msgf("uid=%s sync conv list completed, size=%d", uid, len(items))

	b.cm.ReplaceConvList(items)

	return items, nil
}

func (b *SyncBinder) SyncConvMessages(convId string) ([]*chat.ConvItem, error) {
	return nil, nil
}
