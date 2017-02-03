package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WXUserQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
	RowCount  int `json:"rowCount"`
}

type WXUserQueryTaskResult struct {
	app.Result
	Counter *WXUserQueryCounter `json:"counter,omitempty"`
	Users   []User              `json:"users,omitempty"`
}

type WXUserQueryTask struct {
	app.Task
	Uid       int64  `json:"uid"`
	Openid    string `json:"openid"`
	Keyword   string `json:"q"`
	OrderBy   string `json:"orderBy"` // desc, asc
	PageIndex int    `json:"p"`
	PageSize  int    `json:"size"`
	Counter   bool   `json:"counter"`
	Result    WXUserQueryTaskResult
}

func (T *WXUserQueryTask) GetResult() interface{} {
	return &T.Result
}

func (T *WXUserQueryTask) GetInhertType() string {
	return "weixin"
}

func (T *WXUserQueryTask) GetClientName() string {
	return "User.Query"
}
