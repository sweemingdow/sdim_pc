package chat

import (
	"fmt"
	"sdim_pc/backend/preinld"
	"strings"
)

var (
	cmdType2handler = map[preinld.SubCmdType]func(content map[string]any, myUid string){
		preinld.SubCmdGroupInvited: _cmdGroupInvitedHandler,
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
