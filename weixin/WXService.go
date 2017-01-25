package weixin

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"log"
	"net/http"
	"sort"
	"time"
)

type WXService struct {
	app.Service
	ca *x509.CertPool

	Init   *app.InitTask
	Ticket *WXTicketTask
	Token  *WXTokenTask
	Config *WXConfigTask
}

func (S *WXService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *WXService) HandleInitTask(a app.IApp, task *app.InitTask) error {

	S.ca = x509.NewCertPool()
	S.ca.AppendCertsFromPEM(pemCerts)

	return nil
}

func (S *WXService) HandleWXTokenTask(a IWeixinApp, task *WXTokenTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	rows, err := kk.DBQuery(db, a.GetTokenTable(), a.GetPrefix(), " WHERE appid=?", a.GetAppid())

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	v := Token{}

	if rows.Next() {
		scanner := kk.NewDBScaner(&v)
		err = scanner.Scan(rows)
		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

		now := time.Now().Unix()

		if v.Ctime+v.Expires-10 > now {
			task.Result.Token = &v
			return nil
		}
	} else {
		v.Appid = a.GetAppid()
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: S.ca},
		},
	}

	resp, err := client.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", a.GetAppid(), a.GetSecret()))

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
	} else if resp.StatusCode == 200 {
		var body = make([]byte, resp.ContentLength)
		_, _ = resp.Body.Read(body)
		defer resp.Body.Close()

		log.Println(string(body))

		var data interface{} = nil

		err = json.Decode(body, &data)

		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

		errno := dynamic.IntValue(dynamic.Get(data, "errcode"), 0)

		if errno != 0 {
			task.Result.Errno = int(errno)
			task.Result.Errmsg = dynamic.StringValue(dynamic.Get(data, "errcode"), "")
			return nil
		}

		v.Token = dynamic.StringValue(dynamic.Get(data, "access_token"), v.Token)
		v.Expires = dynamic.IntValue(dynamic.Get(data, "expires_in"), v.Expires)
		v.Ctime = time.Now().Unix()

		if v.Id == 0 {
			_, err = kk.DBInsert(db, a.GetTokenTable(), a.GetPrefix(), &v)
			if err != nil {
				task.Result.Errno = ERROR_WEIXIN
				task.Result.Errmsg = err.Error()
				return nil
			}
		} else {
			_, err = kk.DBUpdate(db, a.GetTokenTable(), a.GetPrefix(), &v)
			if err != nil {
				task.Result.Errno = ERROR_WEIXIN
				task.Result.Errmsg = err.Error()
				return nil
			}
		}

		task.Result.Token = &v

	} else {
		var body = make([]byte, resp.ContentLength)
		_, _ = resp.Body.Read(body)
		defer resp.Body.Close()
		log.Println(string(body))
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = fmt.Sprintf("[%d] %s", resp.StatusCode, string(body))
	}

	return nil
}

func (S *WXService) HandleWXTicketTask(a IWeixinApp, task *WXTicketTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	rows, err := kk.DBQuery(db, a.GetTicketTable(), a.GetPrefix(), " WHERE appid=?", a.GetAppid())

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	v := Ticket{}

	if rows.Next() {
		scanner := kk.NewDBScaner(&v)
		err = scanner.Scan(rows)
		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

		now := time.Now().Unix()

		if v.Ctime+v.Expires-10 > now {
			task.Result.Ticket = &v
			return nil
		}
	} else {
		v.Appid = a.GetAppid()
	}

	var access_token = ""

	{
		t := WXTokenTask{}
		app.Handle(a, &t)
		if t.Result.Token == nil {
			task.Result.Errno = t.Result.Errno
			task.Result.Errmsg = t.Result.Errmsg
			return nil
		}
		access_token = t.Result.Token.Token
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: S.ca},
		},
	}

	resp, err := client.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", access_token))

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
	} else if resp.StatusCode == 200 {
		var body = make([]byte, resp.ContentLength)
		_, _ = resp.Body.Read(body)
		defer resp.Body.Close()

		log.Println(string(body))

		var data interface{} = nil

		err = json.Decode(body, data)

		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

		errno := dynamic.IntValue(dynamic.Get(data, "errcode"), 0)

		if errno != 0 {
			task.Result.Errno = int(errno)
			task.Result.Errmsg = dynamic.StringValue(dynamic.Get(data, "errcode"), "")
			return nil
		}

		v.Ticket = dynamic.StringValue(dynamic.Get(data, "ticket"), v.Ticket)
		v.Expires = dynamic.IntValue(dynamic.Get(data, "expires_in"), v.Expires)
		v.Ctime = time.Now().Unix()

		if v.Id == 0 {
			_, err = kk.DBInsert(db, a.GetTicketTable(), a.GetPrefix(), &v)
			if err != nil {
				task.Result.Errno = ERROR_WEIXIN
				task.Result.Errmsg = err.Error()
				return nil
			}
		} else {
			_, err = kk.DBUpdate(db, a.GetTicketTable(), a.GetPrefix(), &v)
			if err != nil {
				task.Result.Errno = ERROR_WEIXIN
				task.Result.Errmsg = err.Error()
				return nil
			}
		}

		task.Result.Ticket = &v

	} else {
		var body = make([]byte, resp.ContentLength)
		_, _ = resp.Body.Read(body)
		defer resp.Body.Close()
		log.Println(string(body))
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = fmt.Sprintf("[%d] %s", resp.StatusCode, string(body))
	}

	return nil
}

func NewNonceStr() string {
	m := md5.New()
	m.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return hex.EncodeToString(m.Sum(nil))
}

func (S *WXService) HandleWXConfigTask(a IWeixinApp, task *WXConfigTask) error {

	var ticket string = ""

	{
		t := WXTicketTask{}
		app.Handle(a, &t)
		if t.Result.Ticket == nil {
			task.Result.Errno = t.Result.Errno
			task.Result.Errmsg = t.Result.Errmsg
			return nil
		}
		ticket = t.Result.Ticket.Ticket
	}

	noncestr := NewNonceStr()
	timestamp := time.Now().Unix()

	data := map[string]interface{}{}
	data["noncestr"] = noncestr
	data["timestamp"] = timestamp
	data["jsapi_ticket"] = ticket
	data["url"] = task.Url

	keys := []string{}

	for key, _ := range data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	b := bytes.NewBuffer(nil)

	for i, key := range keys {
		if i != 0 {
			b.WriteString("&")
		}
		b.WriteString(key)
		b.WriteString("=")
		b.WriteString(dynamic.StringValue(data[key], ""))
	}

	m := sha1.New()
	m.Write(b.Bytes())

	task.Result.Signature = hex.EncodeToString(m.Sum(nil))
	task.Result.NonceStr = noncestr
	task.Result.Timestamp = timestamp
	task.Result.Appid = a.GetAppid()

	return nil
}
