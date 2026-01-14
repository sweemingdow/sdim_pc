package userbinder

import (
	"sdim_pc/backend/api/userapi"
	"sdim_pc/backend/mylog"
)

type UserBinder struct {
	ua *userapi.UserApi
}

func NewUserBinder(ua *userapi.UserApi) *UserBinder {
	return &UserBinder{
		ua: ua,
	}
}

func (ub *UserBinder) UserProfile(uid string) (userapi.UserProfile, error) {
	mylog.GetLogger().Debug().Msgf("uid=%s fetch user profile start", uid)
	return ub.ua.UserProfile(uid)
}
