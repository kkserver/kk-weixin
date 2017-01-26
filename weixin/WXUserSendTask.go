package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WXUserSendTaskResult struct {
	app.Result
}

type WXUserSendTask struct {
	app.Task
	Openid     string      `json:"openid"`
	Uid        int64       `json:"uid"`
	TempleteId string      `json:"templeteId"`
	Url        string      `json:"url"`
	Data       interface{} `json:"data"`
	Result     WXUserSendTaskResult
}

func (task *WXUserSendTask) GetResult() interface{} {
	return &task.Result
}

func (task *WXUserSendTask) GetInhertType() string {
	return "weixin"
}

func (task *WXUserSendTask) GetClientName() string {
	return "User.Send"
}
