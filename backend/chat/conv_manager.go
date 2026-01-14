package chat

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"sdim_pc/backend/appctx"
	"sdim_pc/backend/preinld"
	"sdim_pc/backend/user"
	"sdim_pc/backend/utils"
	"sync"
	"time"
)

// 会话管理器
type ConvManager struct {
	items     []*ConvItem
	id2idx    map[string]int          // 会话id => idx
	msgId2msg map[string]*preinld.Msg // 消息id => 消息
	rw        sync.RWMutex
}

func NewConvManager() *ConvManager {
	cm := &ConvManager{
		items:     make([]*ConvItem, 0, 100),
		id2idx:    make(map[string]int, 100),
		msgId2msg: make(map[string]*preinld.Msg, 200),
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
		HasMore      bool             `json:"hasMore"`
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

func (cm *ConvManager) TailMsgId(convId string) (int64, bool) {
	cm.rw.Lock()
	defer cm.rw.Unlock()
	idx, ok := cm.id2idx[convId]

	if !ok {
		return 0, false
	}

	convItem := cm.items[idx]

	if len(convItem.RecentlyMsgs) == 0 {
		return 0, true
	}

	return convItem.RecentlyMsgs[len(convItem.RecentlyMsgs)-1].MsgId, true
}

func (cm *ConvManager) AppendTailMsgs(convId string, msgs []*preinld.Msg) ([]*ConvItem, int, bool) {
	hasMore := len(msgs) > 0

	cm.rw.Lock()
	defer cm.rw.Unlock()

	idx, ok := cm.id2idx[convId]

	if !ok {
		return nil, -1, false
	}

	convItem := cm.items[idx]
	if hasMore {
		convItem.HasMore = true
		convItem.RecentlyMsgs = append(convItem.RecentlyMsgs, msgs...)
	} else {
		convItem.HasMore = false
	}

	return cm.items, idx, true
}

func (cm *ConvManager) InsertMsgWhileSend(msd preinld.MsgSendData) ([]*ConvItem, int, string, bool) {
	ui := user.GetUserInfo()

	clientId := msd.ClientId
	if clientId == "" {
		clientId = utils.RandStr(32)
	}

	cm.rw.Lock()
	defer cm.rw.Unlock()

	idx, ok := cm.id2idx[msd.ConvId]
	if !ok {
		return nil, -1, "", false
	}

	convItem := cm.items[idx]

	lastMsg := &preinld.Msg{
		MsgId:    0,
		ConvId:   msd.ConvId,
		Sender:   ui.Uid,
		Receiver: msd.Receiver,
		ChatType: msd.ChatType,
		MsgType:  msd.MsgContent.Type,
		Content:  msd.MsgContent,
		SenderInfo: preinld.SenderInfo{
			Nickname: ui.Nickname,
			Avatar:   ui.Avatar,
		},
		Cts:      time.Now().UnixMilli(),
		State:    uint8(Sending),
		IsSelf:   true,
		ClientId: clientId,
	}

	convItem.LastMsg = lastMsg
	convItem.RecentlyMsgs = append([]*preinld.Msg{lastMsg}, convItem.RecentlyMsgs...)

	cm.msgId2msg[clientId] = lastMsg

	return cm.items, idx, clientId, true

}

// 本地插入收到的消息
func (cm *ConvManager) InsertMsgAfterReceived(ffb *preinld.ForwardFrameBody) ([]*ConvItem, int, bool) {
	ui := user.GetUserInfo()

	clientId := ffb.ClientUniqueId

	cm.rw.Lock()
	defer cm.rw.Unlock()

	idx, ok := cm.id2idx[ffb.ConvId]
	if !ok {
		return nil, -1, false
	}

	convItem := cm.items[idx]

	lastMsg := &preinld.Msg{
		MsgId:      ffb.MsgId,
		ConvId:     ffb.ConvId,
		Sender:     ffb.Sender,
		Receiver:   ffb.Receiver,
		ChatType:   ffb.ChatType,
		MsgType:    ffb.MsgContent.Type,
		Content:    ffb.MsgContent,
		SenderInfo: ffb.SenderInfo,
		MegSeq:     ffb.MsgSeq,
		Cts:        ffb.SendTs,
		State:      uint8(SendOk),
		IsSelf:     ui.Uid == ffb.Sender,
		ClientId:   clientId,
	}

	convItem.LastMsg = lastMsg
	convItem.RecentlyMsgs = append([]*preinld.Msg{lastMsg}, convItem.RecentlyMsgs...)

	return cm.items, idx, true

}

func (cm *ConvManager) UpdateMsgWhenSentSuccess(ackBody preinld.SendAckFrameBody) ([]*ConvItem, int, bool) {
	clientMsgId := ackBody.ClientUniqueId
	convId := ackBody.ConvId

	cm.rw.Lock()
	defer cm.rw.Unlock()

	idx, ok := cm.id2idx[convId]
	if !ok {
		return nil, -1, false
	}

	convItem := cm.items[idx]

	convItem.Uts = ackBody.SendTs
	convItem.MsgSeq = ackBody.MsgSeq

	lastMsg, ok := cm.msgId2msg[clientMsgId]
	if !ok {
		return nil, -1, false
	}

	lastMsg.MsgId = ackBody.MsgId
	lastMsg.State = uint8(SendOk)
	lastMsg.MegSeq = ackBody.MsgSeq
	lastMsg.Cts = ackBody.SendTs

	return cm.items, idx, true
}

func (cm *ConvManager) UpdateMsgWhenSentFailed(saf preinld.SendFrameAck) ([]*ConvItem, int, bool) {
	ackBody := saf.Data
	clientMsgId := ackBody.ClientUniqueId
	convId := ackBody.ConvId

	cm.rw.Lock()
	defer cm.rw.Unlock()

	idx, ok := cm.id2idx[convId]
	if !ok {
		return nil, -1, false
	}

	lastMsg, ok := cm.msgId2msg[clientMsgId]
	if !ok {
		return nil, -1, false
	}

	lastMsg.MsgId = ackBody.MsgId
	lastMsg.State = uint8(SendFailed)
	lastMsg.LastFailReason = fmt.Sprintf("[errCode=%d]-[errDesc=%s]", saf.ErrCode, saf.ErrDesc)

	return cm.items, idx, true
}

func (cm *ConvManager) emitConvListUpdateEvent() {
	cm.rw.RLock()
	convItems := cm.items
	cm.rw.RUnlock()

	cm.EmitConvListUpdateEvent(convItems, -1)

}

type ConvListUpdateData struct {
	Items []*ConvItem `json:"items,omitempty"`
	Idx   int         `json:"idx"`
}

func (cm *ConvManager) EmitConvListUpdateEvent(items []*ConvItem, idx int) {
	runtime.EventsEmit(
		appctx.GetAppCtx(),
		"event_backend:conv_list_update",
		ConvListUpdateData{
			Items: items,
			Idx:   idx,
		},
	)
}
