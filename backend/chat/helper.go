package chat

import (
	"fmt"
	"sdim_pc/backend/preinld"
	"strings"
)

var (
	cmdType2handler = map[preinld.SubCmdType]func(content map[string]any, myUid string){
		preinld.SubCmdGroupInvited:       _cmdGroupInvitedHandler,
		preinld.SubCmdGroupSettingName:   _cmdGroupSettingHandler,
		preinld.SubCmdGroupRemoveMembers: _cmdGroupRemMembersHandler,
		preinld.SubCmdGroupAddMembers:    _cmdGroupAddMembersHandler,
	}
)

func RewriteContentIfNeed(msg *preinld.Msg, myUid string) {
	msgContent := msg.Content
	msgType := msgContent.Type
	if msgType <= preinld.CustomType {
		return
	}

	content := msgContent.Content
	if msgType == preinld.CmdType {
		_subCmd, _ := content["subCmd"].(float64)
		subCmd := preinld.SubCmdType(_subCmd)

		if subCmd == preinld.SubCmdGroupInvited {
			if f, ok := cmdType2handler[preinld.SubCmdGroupInvited]; ok {
				f(content, myUid)
			}
		} else if subCmd == preinld.SubCmdGroupSettingName {
			if f, ok := cmdType2handler[preinld.SubCmdGroupSettingName]; ok {
				f(content, myUid)
			}
		} else if subCmd == preinld.SubCmdGroupRemoveMembers {
			if f, ok := cmdType2handler[preinld.SubCmdGroupRemoveMembers]; ok {
				f(content, myUid)
			}
		} else if subCmd == preinld.SubCmdGroupAddMembers {
			if f, ok := cmdType2handler[preinld.SubCmdGroupAddMembers]; ok {
				f(content, myUid)
			}
		}
	}
}

func _cmdGroupInvitedHandler(content map[string]any, myUid string) {
	inviteContent, _ := content["inviteContent"].(map[string]any)
	//groupMebCount, _ := inviteContent["groupMebCount"].(float64)
	inviteFmtItems, _ := inviteContent["inviteFmtItems"].([]any)
	inviteHint, _ := inviteContent["inviteHint"].(string)

	for idx, val := range inviteFmtItems {
		m, _ := val.(map[string]any)
		nickname, _ := m["nickname"].(string)
		uid, _ := m["uid"].(string)
		if uid == myUid {
			nickname = "你"
		}

		inviteHint = strings.ReplaceAll(inviteHint, fmt.Sprintf("{%d}", idx), nickname)
	}

	inviteContent["inviteHint"] = inviteHint
}

func _cmdGroupSettingHandler(content map[string]any, myUid string) {
	settingContent, _ := content["settingContent"].(map[string]any)
	settingFmtItems, _ := settingContent["settingFmtItems"].([]any)
	settingHint, _ := settingContent["settingHint"].(string)

	for idx, val := range settingFmtItems {
		m, _ := val.(map[string]any)
		nickname, _ := m["nickname"].(string)
		uid, _ := m["uid"].(string)
		if uid == myUid {
			nickname = "你"
		}

		settingHint = strings.ReplaceAll(settingHint, fmt.Sprintf("{%d}", idx), nickname)
	}

	settingContent["settingHint"] = settingHint
}

func _cmdGroupRemMembersHandler(content map[string]any, myUid string) {
	realContent, _ := content["removeContent"].(map[string]any)
	fmtItems, _ := content["groupRemItems"].([]any)
	hint, _ := content["remHint"].(string)

	for idx, val := range fmtItems {
		m, _ := val.(map[string]any)
		nickname, _ := m["nickname"].(string)
		uid, _ := m["uid"].(string)
		if uid == myUid {
			nickname = "你"
		}

		hint = strings.ReplaceAll(hint, fmt.Sprintf("{%d}", idx), nickname)
	}

	realContent["remHint"] = hint
}

func _cmdGroupAddMembersHandler(content map[string]any, myUid string) {
	realContent, _ := content["addContent"].(map[string]any)
	fmtItems, _ := content["groupAddItems"].([]any)
	hint, _ := content["addHint"].(string)

	for idx, val := range fmtItems {
		m, _ := val.(map[string]any)
		nickname, _ := m["nickname"].(string)
		uid, _ := m["uid"].(string)
		if uid == myUid {
			nickname = "你"
		}

		hint = strings.ReplaceAll(hint, fmt.Sprintf("{%d}", idx), nickname)
	}

	realContent["addHint"] = hint
}
