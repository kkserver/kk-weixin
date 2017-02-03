package weixin

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type WXUserService struct {
	app.Service

	Auth        *WXAuthTask
	AuthConfirm *WXAuthConfirmTask
	Get         *WXUserTask
	Bind        *WXUserBindTask
	UnBind      *WXUserUnBindTask
	UnBindAll   *WXUserUnBindAllTask
	Send        *WXUserSendTask
	Query       *WXUserQueryTask
}

func (S *WXUserService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *WXUserService) HandleWXAuthTask(a IWeixinApp, task *WXAuthTask) error {

	if task.Scope == "" {
		task.Scope = WXAuthScopeBase
	}

	task.Result.Url = fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect",
		a.GetAppid(), url.QueryEscape(task.Url), url.QueryEscape(task.Scope), url.QueryEscape(task.State))

	return nil
}

func (S *WXUserService) HandleWXAuthConfirmTask(a IWeixinApp, task *WXAuthConfirmTask) error {

	var url = fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		a.GetAppid(), a.GetSecret(), task.Code)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: a.GetCA()},
		},
	}

	resp, err := client.Get(url)

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

		openid := dynamic.StringValue(dynamic.Get(data, "openid"), "")
		access_token := dynamic.StringValue(dynamic.Get(data, "access_token"), "")
		expires_in := dynamic.IntValue(dynamic.Get(data, "expires_in"), 0)

		t := WXUserTask{}
		t.Openid = openid
		t.Token = access_token
		t.Expires = expires_in

		app.Handle(a, &t)

		if t.Result.User == nil {
			task.Result.Errno = t.Result.Errno
			task.Result.Errmsg = t.Result.Errmsg
			return nil
		}

		task.Result.User = t.Result.User

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

func getUser(a IWeixinApp, user *User, openid string, access_token string) error {

	if access_token == "" {

		token := WXTokenTask{}

		app.Handle(a, &token)

		if token.Result.Token == nil {
			return app.NewError(token.Result.Errno, token.Result.Errmsg)
		}

		var url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN",
			token.Result.Token.Token, openid)

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: a.GetCA()},
			},
		}

		resp, err := client.Get(url)

		if err != nil {
			return err
		} else if resp.StatusCode == 200 {
			var body = make([]byte, resp.ContentLength)
			_, _ = resp.Body.Read(body)
			defer resp.Body.Close()

			log.Println(string(body))

			var data interface{} = nil

			err = json.Decode(body, &data)

			if err != nil {
				return err
			}

			errno := dynamic.IntValue(dynamic.Get(data, "errcode"), 0)

			if errno != 0 {
				return app.NewError(int(errno), dynamic.StringValue(dynamic.Get(data, "errmsg"), ""))
			}

			user.Subscribe = int(dynamic.IntValue(dynamic.Get(data, "subscribe"), int64(user.Subscribe)))
			user.SubscribeTime = dynamic.IntValue(dynamic.Get(data, "subscribe_time"), user.SubscribeTime)
			user.Sex = int(dynamic.IntValue(dynamic.Get(data, "sex"), int64(user.Sex)))
			user.Nick = dynamic.StringValue(dynamic.Get(data, "nickname"), user.Nick)
			user.Logo = dynamic.StringValue(dynamic.Get(data, "headimgurl"), user.Logo)
			user.Language = dynamic.StringValue(dynamic.Get(data, "language"), user.Language)
			user.Province = dynamic.StringValue(dynamic.Get(data, "province"), user.Province)
			user.City = dynamic.StringValue(dynamic.Get(data, "city"), user.City)
			user.Country = dynamic.StringValue(dynamic.Get(data, "country"), user.Country)

			return nil

		} else {
			var body = make([]byte, resp.ContentLength)
			_, _ = resp.Body.Read(body)
			defer resp.Body.Close()
			log.Println(string(body))
			return app.NewError(ERROR_WEIXIN, fmt.Sprintf("[%d] %s", resp.StatusCode, string(body)))
		}

	} else {

		var url = fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", access_token, openid)

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: a.GetCA()},
			},
		}

		resp, err := client.Get(url)

		if err != nil {
			return err
		} else if resp.StatusCode == 200 {
			var body = make([]byte, resp.ContentLength)
			_, _ = resp.Body.Read(body)
			defer resp.Body.Close()

			log.Println(string(body))

			var data interface{} = nil

			err = json.Decode(body, &data)

			if err != nil {
				return err
			}

			errno := dynamic.IntValue(dynamic.Get(data, "errcode"), 0)

			if errno != 0 {
				return app.NewError(int(errno), dynamic.StringValue(dynamic.Get(data, "errmsg"), ""))
			}

			user.Sex = int(dynamic.IntValue(dynamic.Get(data, "sex"), int64(user.Sex)))
			user.Nick = dynamic.StringValue(dynamic.Get(data, "nickname"), user.Nick)
			user.Logo = dynamic.StringValue(dynamic.Get(data, "headimgurl"), user.Logo)
			user.Province = dynamic.StringValue(dynamic.Get(data, "province"), user.Province)
			user.City = dynamic.StringValue(dynamic.Get(data, "city"), user.City)
			user.Country = dynamic.StringValue(dynamic.Get(data, "country"), user.Country)

			return nil

		} else {
			var body = make([]byte, resp.ContentLength)
			_, _ = resp.Body.Read(body)
			defer resp.Body.Close()
			log.Println(string(body))
			return app.NewError(ERROR_WEIXIN, fmt.Sprintf("[%d] %s", resp.StatusCode, string(body)))
		}

	}

	return nil
}

func (S *WXUserService) HandleWXUserTask(a IWeixinApp, task *WXUserTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.Uid != 0 {
		sql := bytes.NewBuffer(nil)
		args := []interface{}{}
		sql.WriteString(" WHERE appid=? AND uid=?")
		args = append(args, a.GetAppid(), task.Uid)
		if task.Openid != "" {
			sql.WriteString(" AND openid=?")
			args = append(args, task.Openid)
		}

		rows, err := kk.DBQuery(db, a.GetUserTable(), a.GetPrefix(), sql.String(), args...)

		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

		defer rows.Close()

		if rows.Next() {

			v := User{}
			scanner := kk.NewDBScaner(&v)

			err = scanner.Scan(rows)

			if err != nil {
				task.Result.Errno = ERROR_WEIXIN
				task.Result.Errmsg = err.Error()
				return nil
			}

			task.Result.User = &v

			return nil
		} else {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = "Not Found user"
			return nil
		}
	}

	if task.Openid == "" {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = "Not Found openid"
		return nil
	}

	rows, err := kk.DBQuery(db, a.GetUserTable(), a.GetPrefix(), " WHERE appid=? AND openid=?", a.GetAppid(), task.Openid)

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	v := User{}

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)
		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

		if task.Token != "" {
			_ = getUser(a, &v, task.Openid, task.Token)
			v.Token = task.Token
			v.Expires = task.Expires
			v.Mtime = time.Now().Unix()
			_, _ = kk.DBUpdate(db, a.GetUserTable(), a.GetPrefix(), &v)
		}

		task.Result.User = &v

		return nil

	} else {

		err = getUser(a, &v, task.Openid, "")

		if err != nil {
			if task.Token != "" {
				err = getUser(a, &v, task.Openid, task.Token)
			}
		}

		if err != nil {
			e, ok := err.(*app.Error)
			if ok {
				task.Result.Errno = e.Errno
				task.Result.Errmsg = e.Errmsg
				return nil
			} else {
				task.Result.Errno = ERROR_WEIXIN
				task.Result.Errmsg = err.Error()
				return nil
			}
		}

		v.Appid = a.GetAppid()
		v.Openid = task.Openid
		v.Token = task.Token
		v.Expires = task.Expires
		v.Mtime = time.Now().Unix()
		v.Ctime = v.Mtime

		_, err = kk.DBInsert(db, a.GetUserTable(), a.GetPrefix(), &v)

		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.User = &v

	}

	return nil
}

func (S *WXUserService) HandleWXUserBindTask(a IWeixinApp, task *WXUserBindTask) error {

	user := WXUserTask{}

	user.Openid = task.Openid

	app.Handle(a, &user)

	if user.Result.User == nil {
		task.Result.Errno = user.Result.Errno
		task.Result.Errmsg = user.Result.Errmsg
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	user.Result.User.Uid = task.Uid

	_, err = kk.DBUpdateWithKeys(db, a.GetUserTable(), a.GetPrefix(), user.Result.User, map[string]bool{"uid": true})

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.User = user.Result.User

	return nil
}

func (S *WXUserService) HandleWXUserUnBindTask(a IWeixinApp, task *WXUserUnBindTask) error {

	user := WXUserTask{}

	user.Openid = task.Openid

	app.Handle(a, &user)

	if user.Result.User == nil {
		task.Result.Errno = user.Result.Errno
		task.Result.Errmsg = user.Result.Errmsg
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	if user.Result.User.Uid == task.Uid {

		user.Result.User.Uid = 0

		_, err = kk.DBUpdateWithKeys(db, a.GetUserTable(), a.GetPrefix(), user.Result.User, map[string]bool{"uid": true})

		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

	}

	task.Result.User = user.Result.User

	return nil
}

func (S *WXUserService) HandleWXUserUnBindAllTask(a IWeixinApp, task *WXUserUnBindAllTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	_, err = db.Exec(fmt.Sprintf("UPDATE %s%s SET uid=0 WHERE uid=?", a.GetPrefix(), a.GetUserTable().Name), task.Uid)

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	return nil
}

func (S *WXUserService) HandleWXUserSendTask(a IWeixinApp, task *WXUserSendTask) error {

	opendids := []string{}

	if task.Openid != "" {
		opendids = append(opendids, task.Openid)
	}

	if task.Uid != 0 {

		var db, err = a.GetDB()

		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

		rows, err := db.Query(fmt.Sprintf("SELECT openid FROM %s%s WHERE uid=? ORDER BY id ASC", a.GetPrefix(), a.GetUserTable().Name), task.Uid)

		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

		defer rows.Close()

		openid := ""

		for rows.Next() {

			err = rows.Scan(&openid)

			if err != nil {
				task.Result.Errno = ERROR_WEIXIN
				task.Result.Errmsg = err.Error()
				return nil
			}

			opendids = append(opendids, openid)

		}
	}

	token := WXTokenTask{}

	app.Handle(a, &token)

	if token.Result.Token == nil {
		task.Result.Errno = token.Result.Errno
		task.Result.Errmsg = token.Result.Errmsg
		return nil
	}

	var url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", token.Result.Token.Token)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: a.GetCA()},
		},
	}

	for _, openid := range opendids {

		data := map[interface{}]interface{}{}

		data["touser"] = openid
		data["template_id"] = task.TempleteId
		data["url"] = task.Url
		data["data"] = task.Data

		b, _ := json.Encode(data)

		resp, err := client.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(b))

		if err != nil {
			log.Println("WXUserSendTask", err)
		} else {
			var body = make([]byte, resp.ContentLength)
			_, _ = resp.Body.Read(body)
			defer resp.Body.Close()
			log.Println("WXUserSendTask", resp.StatusCode, string(body))
		}

	}

	return nil
}

func (S *WXUserService) HandleWXUserQueryTask(a IWeixinApp, task *WXUserQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	sql := bytes.NewBuffer(nil)
	args := []interface{}{}

	sql.WriteString(" WHERE appid=?")
	args = append(args, a.GetAppid())

	if task.Uid != 0 {
		sql.WriteString(" AND uid=?")
		args = append(args, task.Uid)
	}

	if task.Openid != "" {
		sql.WriteString(" AND openid=?")
		args = append(args, task.Openid)
	}

	if task.Keyword != "" {
		q := "%" + task.Keyword + "%"
		sql.WriteString(" AND (nick LIKE ? OR province LIKE ? OR city LIKE ?)")
		args = append(args, q, q, q)
	}

	if task.OrderBy == "asc" {
		sql.WriteString(" ORDER BY id ASC")
	} else {
		sql.WriteString(" ORDER BY id DESC")
	}

	var pageIndex = task.PageIndex
	var pageSize = task.PageSize

	if pageIndex < 1 {
		pageIndex = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	if task.Counter {
		var counter = WXUserQueryCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize
		total, err := kk.DBQueryCount(db, a.GetUserTable(), a.GetPrefix(), sql.String(), args...)
		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}
		if total%pageSize == 0 {
			counter.PageCount = total / pageSize
		} else {
			counter.PageCount = total/pageSize + 1
		}
		task.Result.Counter = &counter
	}

	sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	var users = []User{}
	var v = User{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetUserTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_WEIXIN
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_WEIXIN
			task.Result.Errmsg = err.Error()
			return nil
		}

		users = append(users, v)
	}

	task.Result.Users = users

	return nil
}
