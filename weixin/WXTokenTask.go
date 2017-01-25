package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WXTokenTaskResult struct {
	app.Result
	Token *Token `json:"token,omitempty"`
}

type WXTokenTask struct {
	app.Task
	Result WXTokenTaskResult
}

func (task *WXTokenTask) GetResult() interface{} {
	return &task.Result
}

func (task *WXTokenTask) GetInhertType() string {
	return "weixin"
}

func (task *WXTokenTask) GetClientName() string {
	return "Token"
}
