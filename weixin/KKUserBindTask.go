package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WXUserBindTaskResult struct {
	app.Result
	User *User `json:"user,omitempty"`
}

type WXUserBindTask struct {
	app.Task
	Openid string `json:"openid"`
	Uid    int64  `json:"uid"`
	Result WXUserBindTaskResult
}

func (task *WXUserBindTask) GetResult() interface{} {
	return &task.Result
}

func (task *WXUserBindTask) GetInhertType() string {
	return "weixin"
}

func (task *WXUserBindTask) GetClientName() string {
	return "User.Bind"
}
