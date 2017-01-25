package weixin

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
)

type Token struct {
	Id      int64  `json:"id"`
	Appid   string `json:"appid"`
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
	Ctime   int64  `json:"ctime"`
}

type Ticket struct {
	Id      int64  `json:"id"`
	Appid   string `json:"appid"`
	Ticket  string `json:"token"`
	Expires int64  `json:"expires"`
	Ctime   int64  `json:"ctime"`
}

type IWeixinApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetTokenTable() *kk.DBTable
	GetTicketTable() *kk.DBTable
	GetAppid() string
	GetSecret() string
}

type WeixinApp struct {
	app.App
	DB *app.DBConfig

	Remote *remote.Service

	Appid  string
	Secret string

	TokenTable  kk.DBTable
	TicketTable kk.DBTable

	Token  *WXTokenTask
	Ticket *WXTicketTask
	Config *WXConfigTask

	WX *WXService
}

func (C *WeixinApp) GetDB() (*sql.DB, error) {
	return C.DB.Get(C)
}

func (C *WeixinApp) GetPrefix() string {
	return C.DB.Prefix
}

func (C *WeixinApp) GetAppid() string {
	return C.Appid
}

func (C *WeixinApp) GetSecret() string {
	return C.Secret
}

func (C *WeixinApp) GetTokenTable() *kk.DBTable {
	return &C.TokenTable
}

func (C *WeixinApp) GetTicketTable() *kk.DBTable {
	return &C.TicketTable
}
