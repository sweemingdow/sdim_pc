package groupapi

import (
	"fmt"
	"sdim_pc/backend/chat"
	"sdim_pc/backend/config"
	"sdim_pc/backend/utils/unet"
)

type GroupApi struct {
	host      string
	reqSender *unet.HttpSender
}

func NewGroupApi(cfg config.Config, reqSender *unet.HttpSender) *GroupApi {
	return &GroupApi{
		host:      fmt.Sprintf("%s", "http://192.168.1.155:6020"),
		reqSender: reqSender,
	}
}

type StartGroupChatReq struct {
	GroupName  string   `json:"groupName"`
	Avatar     string   `json:"avatar"`
	OwnerUid   string   `json:"ownerUid"`
	LimitedNum int      `json:"limitedNum"`
	Members    []string `json:"members"`
}

// 发起群聊
func (ca *GroupApi) StarGroupChat(req StartGroupChatReq) ([]*chat.ConvItem, error) {
	status, buf, err := ca.reqSender.JsonPost(
		fmt.Sprintf("%s/group/start_chat", ca.host),
		req,
		nil,
	)

	if err != nil {
		return nil, err
	}

	_, err = unet.ParseOrError[string](status, buf)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
