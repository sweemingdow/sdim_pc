package conv

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"sdim_pc/backend/appctx"
	"sdim_pc/backend/preinld"
	"sync"
)

// 会话管理器
type ConvManager struct {
	items  []*ConvItem
	id2idx map[string]int
	rw     sync.RWMutex
}

func NewConvManager() *ConvManager {
	cm := &ConvManager{
		items:  make([]*ConvItem, 0, 100),
		id2idx: make(map[string]int, 100),
	}

	return cm
}

type (
	ConvItem struct {
		ConvId       string           `json:"convId,omitempty"`
		ConvType     preinld.ConvType `json:"convType,omitempty"`
		Icon         string           `json:"icon,omitempty"`
		Title        string           `json:"title,omitempty"`
		RelationId   string           `json:"relationId,omitempty"`
		Remark       string           `json:"remark,omitempty"`
		PinTop       bool             `json:"pinTop,omitempty"`
		NoDisturb    bool             `json:"noDisturb,omitempty"`
		MsgSeq       int64            `json:"msgSeq,omitempty"`
		LastMsg      *preinld.Msg     `json:"lastMsg,omitempty"`
		BrowseMsgSeq int64            `json:"browseMsgSeq,omitempty"`
		UnreadCount  int64            `json:"unreadCount,omitempty"`
		Cts          int64            `json:"cts,omitempty"`
		Uts          int64            `json:"uts,omitempty"`
		RecentlyMsgs []*preinld.Msg   `json:"recentlyMsgs"`
	}
)

func (cm *ConvManager) List() []*ConvItem {
	cm.rw.RLock()
	defer cm.rw.RUnlock()

	return cm.items
}

func (cm *ConvManager) ReplaceConvList(convItems []*ConvItem) {
	cm.rw.Lock()
	cm.items = convItems
	for i, item := range convItems {
		cm.id2idx[item.ConvId] = i
	}
	cm.rw.Unlock()

	cm.emitConvListUpdateEvent()
}

func (cm *ConvManager) Exists(convId string) bool {
	cm.rw.RLock()
	_, ok := cm.id2idx[convId]
	cm.rw.RUnlock()

	return ok
}

func (cm *ConvManager) UpsertOne(item *ConvItem) {

}

func (cm *ConvManager) emitConvListUpdateEvent() {
	cm.rw.RLock()
	convItems := cm.items
	cm.rw.RUnlock()

	runtime.EventsEmit(appctx.GetAppCtx(), "event_backend:conv_list_update", convItems)
}
