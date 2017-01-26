package weixin

import (
	"crypto/x509"
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

type User struct {
	Id            int64  `json:"id"`
	Appid         string `json:"appid"`
	Openid        string `json:"openid"`
	Uid           int64  `json:"uid"`
	Token         string `json:"token"`
	Expires       int64  `json:"expires"`
	Nick          string `json:"nick"`
	Logo          string `json:"logo"`
	Province      string `json:"province"`
	City          string `json:"city"`
	Country       string `json:"country"`
	Subscribe     int    `json:"subscribe"`
	SubscribeTime int64  `json:"subscribeTime"`
	Language      string `json:"language"`
	Sex           int    `json:"sex"`
	Mtime         int64  `json:"mtime"`
	Ctime         int64  `json:"ctime"`
}

type IWeixinApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetTokenTable() *kk.DBTable
	GetTicketTable() *kk.DBTable
	GetUserTable() *kk.DBTable
	GetAppid() string
	GetSecret() string
	GetCA() *x509.CertPool
}

type WeixinApp struct {
	app.App
	DB *app.DBConfig

	Remote *remote.Service

	Appid  string
	Secret string

	TokenTable  kk.DBTable
	TicketTable kk.DBTable
	UserTable   kk.DBTable

	Token  *WXTokenTask
	Ticket *WXTicketTask
	Config *WXConfigTask

	WX *WXService

	User *WXUserService

	ca *x509.CertPool
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

func (C *WeixinApp) GetUserTable() *kk.DBTable {
	return &C.TicketTable
}

func (C *WeixinApp) GetCA() *x509.CertPool {
	if C.ca == nil {
		C.ca = x509.NewCertPool()
		C.ca.AppendCertsFromPEM(pemCerts)
	}
	return C.ca
}
