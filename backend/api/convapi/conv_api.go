package convapi

import (
	"fmt"
	"net/http"
	"sdim_pc/backend/config"
	"sdim_pc/backend/conv"
	"sdim_pc/backend/utils/unet"
	"sdim_pc/backend/wrapper"
)

type ConvApi struct {
	host      string
	reqSender *unet.HttpSender
}

func NewConvApi(cfg config.Config, reqSender *unet.HttpSender) *ConvApi {
	return &ConvApi{
		host:      fmt.Sprintf("%s", "http://192.168.1.155:6020"),
		reqSender: reqSender,
	}
}

// 活跃的会话列表
func (ca *ConvApi) RecentlyConvList(uid string) ([]*conv.ConvItem, error) {
	status, buf, err := ca.reqSender.JsonGet(fmt.Sprintf("%s/conv/recently_list", ca.host), map[string]string{
		"uid": uid,
	})

	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("http status not ok, status=%d", status)
	}

	var resp wrapper.HttpRespWrapper[[]*conv.ConvItem]
	err = wrapper.ParseResp(buf, &resp)
	if err != nil {
		return nil, err
	}

	if resp.IsOK() {
		return resp.Data, nil
	}

	return nil, fmt.Errorf("respone code not ok, code=%s, subCode=%s, msg=%s", resp.Code, resp.SubCode, resp.Msg)
}

func (ca *ConvApi) SyncHotConvList(uid string) ([]*conv.ConvItem, error) {
	status, buf, err := ca.reqSender.JsonGet(fmt.Sprintf("%s/conv/sync/hot_list", ca.host), map[string]string{
		"uid": uid,
	})

	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("http status not ok, status=%d", status)
	}

	var resp wrapper.HttpRespWrapper[[]*conv.ConvItem]
	err = wrapper.ParseResp(buf, &resp)
	if err != nil {
		return nil, err
	}

	if resp.IsOK() {
		return resp.Data, nil
	}

	return nil, fmt.Errorf("respone code not ok, code=%s, subCode=%s, msg=%s", resp.Code, resp.SubCode, resp.Msg)
}
