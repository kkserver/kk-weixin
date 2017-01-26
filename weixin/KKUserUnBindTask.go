package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WXUserUnBindTaskResult struct {
	app.Result
	User *User `json:"user,omitempty"`
}

type WXUserUnBindTask struct {
	app.Task
	Openid string `json:"openid"`
	Uid    int64  `json:"uid"`
	Result WXUserUnBindTaskResult
}

func (task *WXUserUnBindTask) GetResult() interface{} {
	return &task.Result
}

func (task *WXUserUnBindTask) GetInhertType() string {
	return "weixin"
}

func (task *WXUserUnBindTask) GetClientName() string {
	return "User.UnBind"
}
