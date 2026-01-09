package ucopy

import (
	"fmt"
	"sdim_pc/backend/preinld"
	"testing"
)

type (
	ConvItemV1 struct {
		Id           string           `json:"id,omitempty"`
		ConvType     preinld.ConvType `json:"convType,omitempty"`
		MsgSeq       uint32           `json:"msgSeq,omitempty"`
		OwnerUid     string           `json:"ownerUid,omitempty"`
		RelationId   string           `json:"relationId,omitempty"`
		Title        string           `json:"title,omitempty"`
		Icon         string           `json:"icon,omitempty"`
		LastMsg      *preinld.Msg     `json:"lastMsg,omitempty"`
		UnreadCount  int              `json:"unreadCount,omitempty"`
		BrowseMsgReq int              `json:"browseMsgReq,omitempty"`
		PinTop       bool             `json:"pinTop,omitempty"`
		NoDisturb    bool             `json:"noDisturb,omitempty"`
	}

	ConvItemV2 struct {
		Id           string           `json:"id,omitempty"`
		ConvType     preinld.ConvType `json:"convType,omitempty"`
		MsgSeq       uint32           `json:"msgSeq,omitempty"`
		OwnerUid     string           `json:"ownerUid,omitempty"`
		RelationId   string           `json:"relationId,omitempty"`
		Title        string           `json:"title,omitempty"`
		Icon         string           `json:"icon,omitempty"`
		LastMsg      *preinld.Msg     `json:"lastMsg,omitempty"`
		UnreadCount  int              `json:"unreadCount,omitempty"`
		BrowseMsgReq int              `json:"browseMsgReq,omitempty"`
		PinTop       bool             `json:"pinTop,omitempty"`
		NoDisturb    bool             `json:"noDisturb,omitempty"`
		MsgItems     []*preinld.Msg   `json:"msgItems,omitempty"`
	}
)

func TestCopy(t *testing.T) {
	var v1 = []*ConvItemV1{
		{
			Id:         "1",
			ConvType:   preinld.P2pConv,
			MsgSeq:     1,
			OwnerUid:   "woshidage",
			RelationId: "dageshiwo",
			Title:      "niha",
			Icon:       "123.jpg",
			LastMsg: &preinld.Msg{
				MsgId:          123,
				ClientUniqueId: "321",
				SenderInfo: preinld.SenderInfo{
					Nickname: "dagehao",
					Avatar:   "dagehao.jpg",
				},
				Content: &preinld.MsgContent{
					Type: preinld.TextType,
					Content: map[string]any{
						"text": "fd",
					},
				},
				Cts:   1,
				Uts:   1,
				State: 1,
			},
			UnreadCount:  21,
			BrowseMsgReq: 1,
			PinTop:       true,
			NoDisturb:    false,
		},
		{
			Id:         "2",
			ConvType:   preinld.P2pConv,
			MsgSeq:     2,
			OwnerUid:   "woshidage2",
			RelationId: "dageshiwo2",
			Title:      "niha2",
			Icon:       "123.jpg2",
			LastMsg: &preinld.Msg{
				MsgId:          123,
				ClientUniqueId: "3212",
				SenderInfo: preinld.SenderInfo{
					Nickname: "dagehao2",
					Avatar:   "dagehao.jpg2",
				},
				Content: &preinld.MsgContent{
					Type: preinld.TextType,
					Content: map[string]any{
						"text": "fd2",
					},
				},
				Cts:   2,
				Uts:   2,
				State: 2,
			},
			UnreadCount:  212,
			BrowseMsgReq: 12,
			PinTop:       true,
			NoDisturb:    true,
		},
	}

	var v2 []*ConvItemV2
	Cp(&v1, &v2)

	fmt.Printf("v2:%+v\n", v2)
}
