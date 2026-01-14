package userapi

import (
	"fmt"
	"net/http"
	"sdim_pc/backend/config"
	"sdim_pc/backend/user"
	"sdim_pc/backend/utils/unet"
	"sdim_pc/backend/wrapper"
)

type UserApi struct {
	host      string
	reqSender *unet.HttpSender
}

func NewUserApi(cfg config.Config, reqSender *unet.HttpSender) *UserApi {
	return &UserApi{
		host:      fmt.Sprintf("%s", "http://192.168.1.155:6050"),
		reqSender: reqSender,
	}
}

type UserProfile struct {
	Uid      string `json:"uid,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

func (ca *UserApi) UserProfile(uid string) (UserProfile, error) {
	status, buf, err := ca.reqSender.JsonGet(fmt.Sprintf("%s/user/profile", ca.host), map[string]string{
		"uid": uid,
	})

	var zero UserProfile
	if err != nil {
		return zero, err
	}
	if status != http.StatusOK {
		return zero, fmt.Errorf("http status not ok, status=%d", status)
	}

	var resp wrapper.HttpRespWrapper[UserProfile]
	err = wrapper.ParseResp(buf, &resp)
	if err != nil {
		return zero, err
	}

	if resp.IsOK() {
		user.ModifyUnitInfo(resp.Data.Nickname, resp.Data.Avatar)
		return resp.Data, nil
	}

	return zero, fmt.Errorf("respone code not ok, code=%s, subCode=%s, msg=%s", resp.Code, resp.SubCode, resp.Msg)
}
