package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WXConfigTaskResult struct {
	app.Result
	Appid     string `json:"appid,omitempty"`
	NonceStr  string `json:"nonceStr,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type WXConfigTask struct {
	app.Task
	Url    string `json:"url"`
	Result WXConfigTaskResult
}

func (task *WXConfigTask) GetResult() interface{} {
	return &task.Result
}

func (task *WXConfigTask) GetInhertType() string {
	return "weixin"
}

func (task *WXConfigTask) GetClientName() string {
	return "Config"
}
