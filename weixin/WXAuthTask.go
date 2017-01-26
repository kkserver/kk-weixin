package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

const WXAuthScopeBase = "snsapi_base"
const WXAuthScopeUserInfo = "snsapi_userinfo"

type WXAuthTaskResult struct {
	app.Result
	Url string `json:"url,omitempty"`
}

type WXAuthTask struct {
	app.Task
	Url    string `json:"url"`
	Scope  string `json:"scope"`
	State  string `json:"state"`
	Result WXAuthTaskResult
}

func (task *WXAuthTask) GetResult() interface{} {
	return &task.Result
}

func (task *WXAuthTask) GetInhertType() string {
	return "weixin"
}

func (task *WXAuthTask) GetClientName() string {
	return "Auth"
}
