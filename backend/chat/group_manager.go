package chat

import (
	"sdim_pc/backend/preinld"
	"sdim_pc/backend/utils/usli"
	"sync"
)

type GroupData struct {
	GroupNo           string
	GroupName         string
	Role              preinld.GroupRole
	GroupLimitedNum   int
	GroupMebCount     int
	GroupAnnouncement string         // 群公告
	MembersInfo       []GroupMebItem // 群成员信息
	GroupBak          string         // 群备注(仅自己可见)
	NicknameInGroup   string         // 在群中的昵称
}

type GroupMebItem struct {
	Id       int64
	Uid      string
	Nickname string
	Avatar   string
	Role     preinld.GroupRole
}

type GroupManager struct {
	rw            sync.RWMutex
	groupNo2items map[string]*GroupData
}

func NewGroupManager() *GroupManager {
	return &GroupManager{
		groupNo2items: make(map[string]*GroupData),
	}
}

func (m *GroupManager) UpsertGroupData(data *GroupData) {
	m.rw.Lock()
	m.groupNo2items[data.GroupNo] = data
	m.rw.Unlock()
}

func (m *GroupManager) FetchGroupData(groupNo string) *GroupData {
	m.rw.RLock()
	item, ok := m.groupNo2items[groupNo]
	m.rw.RUnlock()

	if ok {
		return item
	}

	return nil
}

func (m *GroupManager) ModifyGroupName(groupNo, groupName string) *GroupData {
	m.rw.Lock()
	defer m.rw.Unlock()

	if item, ok := m.groupNo2items[groupNo]; ok {
		newItem := *item // 浅拷贝 struct
		newItem.GroupName = groupName
		m.groupNo2items[groupNo] = &newItem

		return &newItem
	}

	return nil
}

func (m *GroupManager) ModifyGroupBak(groupNo, groupBak string) *GroupData {
	m.rw.Lock()
	defer m.rw.Unlock()

	if item, ok := m.groupNo2items[groupNo]; ok {
		newItem := *item
		newItem.GroupBak = groupBak
		m.groupNo2items[groupNo] = &newItem

		return &newItem
	}

	return nil
}

func (m *GroupManager) OnMebBeKicked(groupNo string, removeUids []string) bool {
	m.rw.Lock()
	defer m.rw.Unlock()

	if item, ok := m.groupNo2items[groupNo]; ok {
		item.MembersInfo = usli.Diff(item.MembersInfo, removeUids, func(item GroupMebItem) string {
			return item.Uid
		})
		return true
	}

	return false
}

func (m *GroupManager) ResetWhileDisconnected() {
	m.rw.Lock()
	m.groupNo2items = make(map[string]*GroupData)
	m.rw.Unlock()
}

func (m *GroupManager) ModifyGroupNickname(uid, groupNo, nickname string) {
	m.rw.Lock()
	defer m.rw.Unlock()

	if item, ok := m.groupNo2items[groupNo]; ok {
		for idx, mebItem := range item.MembersInfo {
			if uid == mebItem.Uid {
				//mebItem.Nickname = nickname
				item.MembersInfo[idx] = GroupMebItem{
					Id:       mebItem.Id,
					Uid:      mebItem.Uid,
					Nickname: nickname,
					Avatar:   mebItem.Avatar,
					Role:     mebItem.Role,
				}
			}
		}
	}
}
