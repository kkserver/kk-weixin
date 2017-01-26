package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WXAuthConfirmTaskResult struct {
	app.Result
	User *User `json:"user,omitempty"`
}

type WXAuthConfirmTask struct {
	app.Task
	Code   string `json:"code"`
	Result WXAuthTaskResult
}

func (task *WXAuthConfirmTask) GetResult() interface{} {
	return &task.Result
}

func (task *WXAuthConfirmTask) GetInhertType() string {
	return "weixin"
}

func (task *WXAuthConfirmTask) GetClientName() string {
	return "AuthConfirm"
}
