package groupapi

import (
	"fmt"
	"sdim_pc/backend/chat"
	"sdim_pc/backend/config"
	"sdim_pc/backend/preinld"
	"sdim_pc/backend/user"
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
func (a *GroupApi) StarGroupChat(req StartGroupChatReq) ([]*chat.ConvItem, error) {
	status, buf, err := a.reqSender.JsonPost(
		fmt.Sprintf("%s/group/start_chat", a.host),
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

type GroupDataResp struct {
	GroupNo           string            `json:"groupNo,omitempty"`
	GroupName         string            `json:"groupName,omitempty"`
	Role              preinld.GroupRole `json:"role,omitempty"`
	GroupLimitedNum   int               `json:"groupLimitedNum,omitempty"`
	GroupMebCount     int               `json:"groupMebCount,omitempty"`
	GroupAnnouncement string            `json:"groupAnnouncement,omitempty"` // 群公告
	MembersInfo       []MebInfoItem     `json:"membersInfo"`                 // 群成员信息
	GroupBak          string            `json:"groupBak,omitempty"`          // 群备注(仅自己可见)
	NicknameInGroup   string            `json:"nicknameInGroup,omitempty"`   // 在群中的昵称
}

type MebInfoItem struct {
	Id       int64             `json:"id,omitempty"`
	Uid      string            `json:"uid,omitempty"`
	Nickname string            `json:"nickname,omitempty"`
	Avatar   string            `json:"avatar,omitempty"`
	Role     preinld.GroupRole `json:"role,omitempty"`
}

// 群资料
func (a *GroupApi) FetchGroupData(groupNo, uid string) (*GroupDataResp, error) {
	status, buf, err := a.reqSender.JsonGet(
		fmt.Sprintf("%s/group/fetch_group_data", a.host),
		map[string]string{
			"group_no": groupNo,
			"uid":      uid,
		},
	)

	if err != nil {
		return nil, err
	}

	resp, err := unet.ParseOrError[*GroupDataResp](status, buf)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (a *GroupApi) SettingGroupName(groupNo, groupName string) error {
	status, buf, err := a.reqSender.JsonPost(
		fmt.Sprintf("%s/group/setting_group_name", a.host),
		nil,
		map[string]string{
			"group_no":   groupNo,
			"group_name": groupName,
			"uid":        user.GetUid(),
		},
	)

	if err != nil {
		return err
	}

	_, err = unet.ParseOrError[any](status, buf)
	if err != nil {
		return err
	}

	return nil
}

func (a *GroupApi) SettingGroupBak(uid, groupNo, bak string) error {
	status, buf, err := a.reqSender.JsonPost(
		fmt.Sprintf("%s/group/setting_group_bak", a.host),
		nil,
		map[string]string{
			"group_no":  groupNo,
			"group_bak": bak,
			"uid":       uid,
		},
	)

	if err != nil {
		return err
	}

	_, err = unet.ParseOrError[any](status, buf)
	if err != nil {
		return err
	}

	return nil
}

func (a *GroupApi) SettingGroupNickname(uid, groupNo string, nickname string) error {
	status, buf, err := a.reqSender.JsonPost(
		fmt.Sprintf("%s/group/setting_group_nickname", a.host),
		nil,
		map[string]string{
			"group_no": groupNo,
			"nickname": nickname,
			"uid":      uid,
		},
	)

	if err != nil {
		return err
	}

	_, err = unet.ParseOrError[any](status, buf)
	if err != nil {
		return err
	}

	return nil
}

func (a *GroupApi) GroupAddMembers(uid, groupNo string, members []string) error {
	status, buf, err := a.reqSender.JsonPost(
		fmt.Sprintf("%s/group/add_members", a.host),
		map[string]any{
			"uid":     uid,
			"groupNo": groupNo,
			"members": members,
		},
		nil,
	)

	if err != nil {
		return err
	}

	_, err = unet.ParseOrError[any](status, buf)
	if err != nil {
		return err
	}
	return nil
}

func (a *GroupApi) GroupRemMembers(uid, groupNo string, members []string) error {
	status, buf, err := a.reqSender.JsonPost(
		fmt.Sprintf("%s/group/remove_members", a.host),
		map[string]any{
			"uid":     uid,
			"groupNo": groupNo,
			"members": members,
		},
		nil,
	)

	if err != nil {
		return err
	}

	_, err = unet.ParseOrError[any](status, buf)
	if err != nil {
		return err
	}
	return nil
}
