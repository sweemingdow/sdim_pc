package syncbinder

import (
	"sdim_pc/backend/api/convapi"
	"sdim_pc/backend/conv"
	"sdim_pc/backend/mylog"
)

type SyncBinder struct {
	ca *convapi.ConvApi
}

func NewSyncBinder(ca *convapi.ConvApi) *SyncBinder {
	return &SyncBinder{
		ca: ca,
	}
}

func (b *SyncBinder) SyncConvList(uid string) ([]*conv.ConvItem, error) {
	mylog.GetLogger().Debug().Msgf("uid=%s start sync conv list after connected success", uid)

	items, err := b.ca.SyncHotConvList(uid)
	if err != nil {
		return make([]*conv.ConvItem, 0), err
	}

	mylog.GetLogger().Debug().Msgf("uid=%s sync conv list completed", uid)

	return items, nil
}

func (b *SyncBinder) SyncConvMessages(convId string) ([]*conv.ConvItem, error) {
	return nil, nil
}
