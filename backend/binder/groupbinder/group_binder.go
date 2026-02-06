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
	gm *chat.GroupManager
}

func NewGroupBinder(gi *groupapi.GroupApi, cm *chat.ConvManager, gm *chat.GroupManager) *GroupBinder {
	return &GroupBinder{
		gi: gi,
		cm: cm,
		gm: gm,
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

func (b *GroupBinder) FetchGroupData(groupNo string) (*groupapi.GroupDataResp, error) {
	uid := user.GetUid()
	mylog.GetLogger().Debug().Msgf("start fetch group data, groupNo=%s, uid=%s", groupNo, uid)

	// 先查缓存
	if cached := b.gm.FetchGroupData(groupNo); cached != nil {
		// 深拷贝返回，避免外部修改缓存
		resp := &groupapi.GroupDataResp{
			GroupNo:           cached.GroupNo,
			GroupName:         cached.GroupName,
			Role:              cached.Role,
			GroupLimitedNum:   cached.GroupLimitedNum,
			GroupMebCount:     cached.GroupMebCount,
			GroupAnnouncement: cached.GroupAnnouncement,
			GroupBak:          cached.GroupBak,
			NicknameInGroup:   cached.NicknameInGroup,
			MembersInfo:       make([]groupapi.MebInfoItem, len(cached.MembersInfo)),
		}
		for i, m := range cached.MembersInfo {
			resp.MembersInfo[i] = groupapi.MebInfoItem{
				Id:       m.Id,
				Uid:      m.Uid,
				Nickname: m.Nickname,
				Avatar:   m.Avatar,
				Role:     m.Role,
			}
		}
		mylog.GetLogger().Debug().Msgf("hit group data cache, groupNo=%s, uid=%s", groupNo, uid)
		return resp, nil
	}

	// API 查询
	respData, err := b.gi.FetchGroupData(groupNo, uid)
	if err != nil {
		return nil, err
	}

	// 构造缓存对象（深拷贝）
	cacheData := &chat.GroupData{
		GroupNo:           respData.GroupNo,
		GroupName:         respData.GroupName,
		Role:              respData.Role,
		GroupLimitedNum:   respData.GroupLimitedNum,
		GroupMebCount:     respData.GroupMebCount,
		GroupAnnouncement: respData.GroupAnnouncement,
		GroupBak:          respData.GroupBak,
		NicknameInGroup:   respData.NicknameInGroup,
		MembersInfo:       make([]chat.GroupMebItem, len(respData.MembersInfo)),
	}
	for i, m := range respData.MembersInfo {
		cacheData.MembersInfo[i] = chat.GroupMebItem{
			Id:       m.Id,
			Uid:      m.Uid,
			Nickname: m.Nickname,
			Avatar:   m.Avatar,
			Role:     m.Role,
		}
	}

	b.gm.UpsertGroupData(cacheData) // 存储新分配的对象
	return respData, nil
}

func (b *GroupBinder) SettingGroupName(groupNo, groupName string) error {
	uid := user.GetUid()
	mylog.GetLogger().Debug().Msgf("start fetch group data, groupNo=%s, uid=%s", groupNo, uid)

	return b.gi.SettingGroupName(groupNo, groupName)
}

// 设置群备注, 仅自身可见的备注
func (b *GroupBinder) SettingGroupBak(groupNo, bak string) error {
	uid := user.GetUid()
	mylog.GetLogger().Debug().Msgf("start setting group bak, groupNo=%s, uid=%s, bak=%s", groupNo, uid, bak)

	err := b.gi.SettingGroupBak(uid, groupNo, bak)
	if err != nil {
		return err
	}

	b.gm.ModifyGroupBak(groupNo, bak)

	return nil
}

// 设置群内昵称
func (b *GroupBinder) SettingNicknameInGroup(groupNo, nickname string) error {
	uid := user.GetUid()
	mylog.GetLogger().Debug().Msgf("start setting group nickname, groupNo=%s, uid=%s, nickname=%s", groupNo, uid, nickname)

	err := b.gi.SettingGroupNickname(uid, groupNo, nickname)
	if err != nil {
		return err
	}

	b.gm.ModifyGroupNickname(uid, groupNo, nickname)

	return nil
}

// 添加群成员
func (b *GroupBinder) GroupAddMembers(groupNo string, members []string) error {
	uid := user.GetUid()
	mylog.GetLogger().Debug().Msgf("start add group members, uid=%s, groupNo=%s, uids=%v", groupNo, uid, members)

	err := b.gi.GroupAddMembers(uid, groupNo, members)
	if err != nil {
		return err
	}

	//b.gm.ModifyGroupNickname(uid, groupNo, nickname)

	return nil
}

// 移除群成员
func (b *GroupBinder) GroupRemMembers(groupNo string, members []string) error {
	uid := user.GetUid()
	mylog.GetLogger().Debug().Msgf("start remove group members, uid=%s, groupNo=%s, uids=%v", groupNo, uid, members)

	err := b.gi.GroupRemMembers(uid, groupNo, members)
	if err != nil {
		return err
	}

	//b.gm.ModifyGroupNickname(uid, groupNo, nickname)

	return nil
}
