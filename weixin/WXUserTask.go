package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WXUserTaskResult struct {
	app.Result
	User *User `json:"user,omitempty"`
}

type WXUserTask struct {
	app.Task
	Openid  string `json:"openid"`
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
	Result  WXUserTaskResult
}

func (task *WXUserTask) GetResult() interface{} {
	return &task.Result
}

func (task *WXUserTask) GetInhertType() string {
	return "weixin"
}

func (task *WXUserTask) GetClientName() string {
	return "User.Get"
}
