package msgapi

import (
	"fmt"
	"net/http"
	"sdim_pc/backend/chat"
	"sdim_pc/backend/config"
	"sdim_pc/backend/mylog"
	"sdim_pc/backend/preinld"
	"sdim_pc/backend/user"
	"sdim_pc/backend/utils/unet"
	"sdim_pc/backend/wrapper"
	"strconv"
)

type MsgApi struct {
	host      string
	reqSender *unet.HttpSender
}

func NewMsgApi(cfg config.Config, reqSender *unet.HttpSender) *MsgApi {
	return &MsgApi{
		host:      fmt.Sprintf("%s", "http://192.168.1.6:6040"),
		reqSender: reqSender,
	}
}

func (ma *MsgApi) FetchNextMsgs(convId string, lastMsgId int64) ([]*preinld.Msg, error) {
	status, buf, err := ma.reqSender.JsonGet(fmt.Sprintf("%s/msg/conv_history_msgs", ma.host), map[string]string{
		"conv_id":     convId,
		"last_msg_id": strconv.FormatInt(lastMsgId, 10),
	})

	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("http status not ok, status=%d", status)
	}

	var resp wrapper.HttpRespWrapper[[]*preinld.Msg]
	err = wrapper.ParseResp(buf, &resp)
	if err != nil {
		return nil, err
	}

	if resp.IsOK() {
		if e := mylog.GetLogger().Trace(); e.Enabled() {
			e.Msgf("fetch next msgs success, size:%d", len(resp.Data))
		}

		myUid := user.GetUid()

		for _, msg := range resp.Data {
			msg.IsSelf = myUid == msg.Sender
			msg.State = uint8(chat.SendOk)
		}

		return resp.Data, nil
	}

	return nil, fmt.Errorf("respone code not ok, code=%s, subCode=%s, msg=%s", resp.Code, resp.SubCode, resp.Msg)
}
