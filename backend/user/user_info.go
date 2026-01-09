package user

import (
	"sync"
	"time"
)

var (
	ui *UserInfo
	rw sync.RWMutex
)

type UserInfo struct {
	Uid        string
	ApiToken   string
	SignKey    string
	ClientType uint8
	CTime      time.Time
	UTime      time.Time
}

func Replace(u UserInfo) {
	rw.Lock()
	ui = &u
	rw.Unlock()
}

func Reset() {
	rw.Lock()
	defer rw.Unlock()

	ui = &UserInfo{}
}

func ModifySignKey(sk string) {
	rw.Lock()
	ui.SignKey = sk
	ui.UTime = time.Now()
	rw.Unlock()
}

func GetUserInfo() UserInfo {
	rw.RLock()
	u := *ui
	rw.RUnlock()

	return u
}

func GetUid() string {
	return GetUserInfo().Uid
}
