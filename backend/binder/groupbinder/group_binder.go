package groupbinder

import (
	"sdim_pc/backend/api/groupapi"
	"sdim_pc/backend/chat"
	"sdim_pc/backend/mylog"
	"sdim_pc/backend/user"
	"strconv"
	"strings"
)

type GroupBinder struct {
	gi *groupapi.GroupApi
	cm *chat.ConvManager
}

func NewGroupBinder(gi *groupapi.GroupApi, cm *chat.ConvManager) *GroupBinder {
	return &GroupBinder{
		gi: gi,
		cm: cm,
	}
}

type StartGroupChatData struct {
	GroupName  string `json:"groupName"`
	Avatar     string `json:"avatar"`
	LimitedNum string `json:"limitedNum"`
	MembersStr string `json:"membersStr"`
}

func (b *GroupBinder) StarGroupChat(data StartGroupChatData) error {
	mylog.GetLogger().Debug().Msgf("start group chat, data=%+v", data)

	limitNum, err := strconv.Atoi(data.LimitedNum)
	if err != nil {
		return err
	}

	var req groupapi.StartGroupChatReq
	req.Avatar = data.Avatar
	req.GroupName = data.GroupName
	req.OwnerUid = user.GetUid()
	req.LimitedNum = limitNum
	req.Members = strings.Split(data.MembersStr, ",")

	_, err = b.gi.StarGroupChat(req)

	return err
}
