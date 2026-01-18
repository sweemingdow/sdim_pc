package chat

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"sdim_pc/backend/appctx"
	"sdim_pc/backend/mylog"
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
	msgId2msg map[string]*preinld.Msg // 维护消息发送: 消息id => 消息
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

func (cm *ConvManager) ResetWhileDisconnected() {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	cm.items = make([]*ConvItem, 0)
	cm.id2idx = make(map[string]int)
	cm.msgId2msg = make(map[string]*preinld.Msg)
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

func (cm *ConvManager) UpdateWhenConvUpdate(cuf *preinld.ConvUpdateFrame) ([]*ConvItem, int, bool) {
	convId := cuf.ConvId
	lg := mylog.GetLogger().With().Str("conv_id", cuf.ConvId).Logger()
	if cuf.Type == preinld.ConvLastMsgUpdated {
		//ui := user.GetUserInfo()

		lg.Debug().Msgf("start handle conv update event, data=%+v", cuf)

		cm.rw.Lock()
		defer cm.rw.Unlock()

		idx, ok := cm.id2idx[convId]
		if !ok {
			return nil, -1, false
		}

		convItem := cm.items[idx]

		/*
			type LastMsg struct {
				MsgId      int64       `json:"msgId"`
				SenderInfo SenderInfo  `json:"senderInfo"`
				Content    *MsgContent `json:"content"`
			}
		*/
		lastMsgMap, ok := cuf.Data["lastMsg"].(map[string]any)
		if !ok {
			return nil, -1, false
		}

		lastActiveTsVal, ok := cuf.Data["lastActiveTs"].(float64)
		if !ok {
			return nil, -1, false
		}

		unreadCountVal, ok := cuf.Data["unreadCount"].(float64)
		if !ok {
			return nil, -1, false
		}

		var msgIdInMap = int64(lastMsgMap["msgId"].(float64))

		convLastMsg := convItem.LastMsg
		if convLastMsg != nil {
			// 更新
			if msgIdInMap == convLastMsg.MsgId {
				convItem.UnreadCount = int64(unreadCountVal)
				convItem.Uts = int64(lastActiveTsVal)
			}
		}

		return cm.items, idx, true
	} else if cuf.Type == preinld.ConvAdded { // 会话新增
		lg.Debug().Msgf("start handle conv add event, data=%+v", cuf)

		myUid := user.GetUid()

		cm.rw.Lock()
		defer cm.rw.Unlock()

		_, ok := cm.id2idx[convId]
		if ok {
			// 更换下图标? 可能要引入版本机制才行
			return nil, -1, false
		}

		// 会话不存在, 前插
		dataMap := cuf.Data
		icon, _ := dataMap["icon"].(string)
		title, _ := dataMap["title"].(string)
		ts, _ := dataMap["ts"].(float64)
		relationId, _ := dataMap["relationId"].(string)
		convType, _ := dataMap["convType"].(float64)
		chatType, _ := dataMap["chatType"].(float64)
		sender, _ := dataMap["sender"].(string)
		receiver, _ := dataMap["receiver"].(string)
		followMsgMap, ok := dataMap["followMsg"].(map[string]any)

		newConvItem := &ConvItem{
			ConvId:     convId,
			ConvType:   preinld.ConvType(convType),
			Title:      title,
			Icon:       icon,
			RelationId: relationId,
			Cts:        int64(ts),
			Uts:        int64(ts),
		}

		if ok {
			msgId, _ := followMsgMap["msgId"].(float64)
			var senderInfo preinld.SenderInfo
			if senderInfoMap, ok := followMsgMap["senderInfo"].(map[string]any); ok {
				nickname, _ := senderInfoMap["nickname"].(string)
				avatar, _ := senderInfoMap["avatar"].(string)

				senderInfo = preinld.SenderInfo{
					Nickname: nickname,
					Avatar:   avatar,
				}
			}
			var (
				msgType    preinld.MsgType
				msgContent *preinld.MsgContent
			)

			if contentMap, ok := followMsgMap["content"].(map[string]any); ok {
				_msgType, _ := contentMap["type"].(float64)
				msgType = preinld.MsgType(_msgType)
				if ok {
					var content, custom, extra map[string]any
					content, _ = contentMap["content"].(map[string]any)
					custom, _ = contentMap["custom"].(map[string]any)
					extra, _ = contentMap["extra"].(map[string]any)
					msgContent = &preinld.MsgContent{
						Type:    msgType,
						Content: content,
						Custom:  custom,
						Extra:   extra,
					}
				}
			}

			lastMsg := &preinld.Msg{
				MsgId:      int64(msgId),
				ConvId:     convId,
				Sender:     sender,
				Receiver:   receiver,
				ChatType:   preinld.ChatType(chatType),
				MsgType:    msgType,
				Content:    msgContent,
				SenderInfo: senderInfo,
				MegSeq:     0,
				Cts:        int64(ts),
				State:      uint8(SendOk),
				IsSelf:     myUid == sender,
			}

			newConvItem.LastMsg = lastMsg
			newConvItem.RecentlyMsgs = append([]*preinld.Msg{lastMsg}, newConvItem.RecentlyMsgs...)
		}

		cm.id2idx = make(map[string]int)
		cm.items = append([]*ConvItem{newConvItem}, cm.items...)
		for i, item := range cm.items {
			cm.id2idx[item.ConvId] = i
		}

		return cm.items, -1, true
	}

	return nil, -1, false
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

func (cm *ConvManager) ShouldClearUnread(convId string) bool {
	cm.rw.RLock()
	defer cm.rw.RUnlock()

	idx, ok := cm.id2idx[convId]
	if !ok {
		return false
	}

	convItem := cm.items[idx]

	return convItem.UnreadCount > 0
}

func (cm *ConvManager) UpdateAfterClearUnread(convId string) ([]*ConvItem, int, bool) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	idx, ok := cm.id2idx[convId]
	if !ok {
		return nil, -1, false
	}

	convItem := cm.items[idx]

	convItem.UnreadCount = 0

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
