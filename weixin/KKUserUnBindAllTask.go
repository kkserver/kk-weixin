package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WXUserUnBindAllTaskResult struct {
	app.Result
}

type WXUserUnBindAllTask struct {
	app.Task
	Uid    int64 `json:"uid"`
	Result WXUserUnBindTaskResult
}

func (task *WXUserUnBindAllTask) GetResult() interface{} {
	return &task.Result
}

func (task *WXUserUnBindAllTask) GetInhertType() string {
	return "weixin"
}

func (task *WXUserUnBindAllTask) GetClientName() string {
	return "User.UnBindAll"
}
