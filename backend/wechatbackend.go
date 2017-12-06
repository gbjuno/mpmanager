package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/jssdk"
	"gopkg.in/chanxuehong/wechat.v2/mp/media"
	"gopkg.in/chanxuehong/wechat.v2/mp/menu"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/callback/request"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/callback/response"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/sid"
	myTemplate "github.com/gbjuno/mpmanager/backend/templates"
	//"github.com/gbjuno/mpmanager/backend/utils"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"gopkg.in/chanxuehong/wechat.v2/oauth2"
)

const (
	wxAppId           = "wx6eb571f36f6b1c10"
	wxAppSecret       = "555210c557802c8a0c6a930cf2e4c159"
	wxOriId           = "gh_75a8e6a73da5"
	wxToken           = "wechatcms"
	wxEncodedAESKey   = "aeskeyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
	noticePageSuccess = "success"
	noticePagefail    = "warn"
)

var (
	wxDESkey                                 = []byte(wxEncodedAESKey)[:8]
	sessionStorage                           = session.New(20*60, 60*60)
	oauth2Endpoint    oauth2.Endpoint        = mpoauth2.NewEndpoint(wxAppId, wxAppSecret)
	accessTokenServer core.AccessTokenServer = core.NewDefaultAccessTokenServer(wxAppId, wxAppSecret, nil)
	wechatClient      *core.Client           = core.NewClient(accessTokenServer, nil)
	ticketServer                             = jssdk.NewDefaultTicketServer(wechatClient)

	// 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
	msgHandler core.Handler
	msgServer  *core.Server
)

type NoticePage struct {
	Title    string
	Type     string
	Msgtitle string
	Msgbody  string
}

func WechatBackendInit() {
	prefix := fmt.Sprintf("[%s]", "WechatBackendInit")
	glog.Infof("%s initialization", prefix)

	if wechatMenu, _, err := menu.Get(wechatClient); err != nil {
		glog.Errorf("%s cannot get menu, err %s", prefix, err)
		if wechatMenu == nil || len(wechatMenu.Buttons) != 3 {
			createMenu()
		}
	} else {
		glog.Errorf("%s get Menu from wechat %v", prefix, wechatMenu)
		if err := menu.Delete(wechatClient); err != nil {
			glog.Errorf("%s cannot delete menu, err %s", prefix, err)
		}
		createMenu()
	}

	mux := core.NewServeMux()
	mux.DefaultMsgHandleFunc(defaultMsgHandler)
	mux.DefaultEventHandleFunc(defaultEventHandler)
	mux.MsgHandleFunc(request.MsgTypeText, textMsgHandler)
	mux.EventHandleFunc(menu.EventTypeClick, menuClickEventHandler)

	msgHandler = mux
	msgServer = core.NewServer(wxOriId, wxAppId, wxToken, wxEncodedAESKey, msgHandler, nil)
}

func createMenu() {
	redirectURIPrefix := fmt.Sprintf("https://%s/backend/%%s", domain)
	oauth2Scope := "snsapi_base"

	bindingButton := &menu.Button{}
	bindingRedirectURI := fmt.Sprintf(redirectURIPrefix, "binding")
	bindingState := "binding"
	bindingURI := mpoauth2.AuthCodeURL(wxAppId, bindingRedirectURI, oauth2Scope, bindingState)
	bindingButton.SetAsViewButton("绑定账户", bindingURI)
	glog.Infof("set Button binding for uri %s, wechat redirecturi %s", bindingRedirectURI, bindingURI)

	scanqrcodeButton := &menu.Button{}
	scanqrcodeRedirectURI := fmt.Sprintf(redirectURIPrefix, "scanqrcode")
	scanqrcodeState := "scanqrcode"
	scanqrcodeURI := mpoauth2.AuthCodeURL(wxAppId, scanqrcodeRedirectURI, oauth2Scope, scanqrcodeState)
	scanqrcodeButton.SetAsViewButton("拍    照", scanqrcodeURI)
	glog.Infof("set Button scanqrcode for uri %s, wechat redirecturi %s", scanqrcodeRedirectURI, scanqrcodeURI)

	checkButton := &menu.Button{}
	checkButton.SetAsClickButton("检查进度", "CL_CHECKPROGRESS")
	wechatMenu := &menu.Menu{Buttons: []menu.Button{*bindingButton, *scanqrcodeButton, *checkButton}}

	if err := menu.Create(wechatClient, wechatMenu); err != nil {
		glog.Fatalf("cannot connect with weixin server, err %s", err)
	}
	glog.Info("create menu succeed")
}

func textMsgHandler(ctx *core.Context) {
	prefix := fmt.Sprintf("[%s]", "textMsgHandler")
	glog.Infof("%s 收到文本消息:\n%s", prefix, ctx.MsgPlaintext)
	msg := request.GetText(ctx.MixedMsg)
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
	ctx.RawResponse(resp) // 明文回复
	//ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultMsgHandler(ctx *core.Context) {
	prefix := fmt.Sprintf("[%s]", "defaultMsgHandler")
	glog.Infof("%s 收到消息:\n%s", prefix, ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func menuClickEventHandler(ctx *core.Context) {
	prefix := fmt.Sprintf("[%s]", "menuClickHandler")
	glog.Infof("%s 收到菜单 click 事件:\n%s", prefix, ctx.MsgPlaintext)
	event := menu.GetClickEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")
	ctx.RawResponse(resp) // 明文回复
	//ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultEventHandler(ctx *core.Context) {
	prefix := fmt.Sprintf("[%s]", "defaultEventHandler")
	glog.Infof("%s 收到事件:\n%s", prefix, ctx.MsgPlaintext)
	ctx.NoneResponse()
}

// wxCallbackHandler 是处理回调请求的 http handler.
//  1. 不同的 web 框架有不同的实现
//  2. 一般一个 handler 处理一个公众号的回调请求(当然也可以处理多个, 这里我只处理一个)
func wxCallbackHandler(w http.ResponseWriter, r *http.Request) {
	msgServer.ServeHTTP(w, r, nil)
}

// if cookie is valid, return session.(string)
// if cookie is invalid, return error
func parseSession(cookieValue string) (string, error) {
	prefix := fmt.Sprintf("[%s]", "parseSession")
	glog.Infof("%s parse cookie, cookie value %s", prefix, cookieValue)

	session, err := sessionStorage.Get(cookieValue)
	if err != nil {
		glog.Errorf("%s session is outdate or invalid, err %s", prefix, err)
		return "", err
	}
	// user session is valid
	return session.(string), nil
}

func newSession(openId string) (string, error) {
	prefix := fmt.Sprintf("[%s]", "newSession")
	glog.Infof("%s openId %s", prefix, openId)
	sid := sid.New()
	if err := sessionStorage.Add(sid, openId); err != nil {
		//fail to set session
		errmsg := fmt.Sprintf("cannot set sid to sessionStorage, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		return "", errors.New(errmsg)
	}
	return sid, nil
}

func getUserOpenId(code string) (string, error) {
	prefix := fmt.Sprintf("[%s]", "getUserOpenId")
	glog.Infof("%s code %s", prefix, code)

	oauth2Client := oauth2.Client{
		Endpoint: oauth2Endpoint,
	}

	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get userinfo token from wechat by using code %s, err %s", code, err)
		glog.Errorf("%s %s", prefix, errmsg)
		return "", errors.New(errmsg)
	}

	glog.Infof("%s get user token %v ", prefix, token)
	glog.Infof("%s get user openid %v ", prefix, token.OpenId)

	return token.OpenId, nil
}

func isUserExist(openId string) bool {
	// user session is valid
	user := User{}
	db.Where("wx_openid = ?", openId).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

// if a user has already bind to a company, return a notify page
// if a user hasn't binded before, send a login page to bind
func bindingHandler(w http.ResponseWriter, r *http.Request) {
	prefix := fmt.Sprintf("[%s] [%s]", r.RemoteAddr, "BindingHandler")
	glog.Infof("%s %s %s", prefix, r.Method, r.RequestURI)

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		glog.Errorf("%s cannot parse url querystring, err %s", prefix, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	code := queryValues.Get("code")
	state := queryValues.Get("state")

	//request isn't redirected by wechat, return notice page
	if code == "" || state == "" {
		glog.Infof("%s request isn't redirect by weixin, return notice page")
		w.WriteHeader(http.StatusOK)
		return
	}

	//if openid related to a user
	var openId string
	var sid string
	var validSession bool = false

	cookie, err := r.Cookie("sid")
	//no cookie sid
	if err == nil {
		glog.Infof("%s get cookie sid %s", prefix, cookie.Value)
		openId, err = parseSession(cookie.Value)
		if err == nil {
			validSession = true
			glog.Infof("%s session is valid, user openid %s", prefix, openId)
		}
	} else {
		glog.Infof("%s no cookie sid", prefix, cookie.Value)
	}

	//cookie is not exist or invalid cookie
	if !validSession {
		glog.Infof("%s cookie not exist or invalid cookie, generating newOne", prefix)
		//get openid
		openId, err = getUserOpenId(code)
		if err != nil {
			glog.Errorf("%s cannot get user openid from wechat, err %s", prefix, err)
			return
		}
		glog.Infof("%s get user openid %s", prefix, openId)

		//set session
		sid, err = newSession(openId)
		if err != nil {
			glog.Errorf("%s cannot get new session, err %s", prefix, err)
			return
		}

		cookie := http.Cookie{
			Name:     "sid",
			Value:    sid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		glog.Infof("%s new session create for user with openid %s", prefix, openId)
	}

	user := User{}
	db.Where("wx_openid = ?", openId).First(&user)
	if user.ID != 0 {
		company := Company{}
		db.First(&company, user.CompanyId)
		glog.Infof("%s openid is related to a user", prefix)
		// openid is related to a user
		msgbody := fmt.Sprintf("用户%s已经绑定企业%s，无须再次绑定即可拍照", user.Name, company.Name)
		n := NoticePage{Title: "绑定企业", Type: noticePageSuccess, Msgtitle: "绑定企业成功", Msgbody: msgbody}
		noticepageTmpl := template.Must(template.New("noticepage").Parse(myTemplate.NOTICEPAGE))
		noticepageTmpl.Execute(w, n)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		glog.Infof("%s openid is not related to a user", prefix)
		// openid isn't related to a user
		bind_tmpl := template.Must(template.New("bind").Parse(myTemplate.BIND))
		bind_tmpl.Execute(w, nil)
		return
	}

	return
}

func confirmHandler(w http.ResponseWriter, r *http.Request) {
	prefix := fmt.Sprintf("[%s]", "confirmHandler")
	glog.Infof("%s %s %s", prefix, r.Method, r.RequestURI)
	glog.Infof("%s user binding started", prefix)

	oauth2RedirectURI := fmt.Sprintf("https://%s/backend/binding", domain)
	oauth2Scope := "snsapi_base"
	state := "binding"
	AuthCodeURL := mpoauth2.AuthCodeURL(wxAppId, oauth2RedirectURI, oauth2Scope, state)

	cookie, err := r.Cookie("sid")
	if err != nil {
		//no cookie sid => request is not acccessed from wechat
		http.Redirect(w, r, AuthCodeURL, http.StatusFound)
		glog.Errorf("%s no cookie sid in request, return with redirect", prefix)
		return
	}

	openId, err := parseSession(cookie.Value)
	if err != nil {
		//cookie is invalid, redirect to start page
		glog.Errorf("%s cookie is invalid, redirect to /backend/binding", prefix)
		http.Redirect(w, r, AuthCodeURL, http.StatusFound)
		return
	}

	var response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
	//openid is related to a user
	user := User{}
	company := Company{}

	db.Where("wx_openid = ?", openId).First(&user)
	if user.ID != 0 {
		db.First(&company, user.CompanyId)
		response.Status = true
		response.Message = fmt.Sprintf("用户%s已成功绑定企业%s，可以进行拍照", user.Name, company.Name)
		glog.Infof("%s openid is related to a user, return notice page", prefix)
		w.WriteHeader(http.StatusOK)
		return
	}

	//openid is not related to a user
	r.ParseForm()
	phone := r.Form.Get("phone")
	password := r.Form.Get("password")
	glog.Infof("%s, postform data phone %s, password %s", prefix, phone, password)

	user = User{}
	db.Where("phone = ?", phone).First(&user)
	if user.ID == 0 {
		glog.Errorf("%s user with phone %s doesn't exist", prefix, phone)
		return
	}
	glog.Infof("%s find user with id %d", prefix, user.ID)
	/*
		encryptedPassword, err := utils.DesEncrypt([]byte(password), wxDESkey)
		if err != nil {
			glog.Errorf("%s unable to encrypt password %s", password)
			glog.Infof("%s password not match", prefix)
			return
		}*/

	// password match
	if password == user.Password {
		user.WxOpenId = openId
		db.Save(&user)
		response.Status = true
		response.Message = fmt.Sprintf("用户%s首次成功绑定企业%s，可以进行拍照", user.Name, company.Name)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		//password not match
		glog.Infof("%s password not match, input password %s, actual password %s", prefix, password, user.Password)
		response.Status = false
		response.Message = "手机号或密码错误，请重试"
		w.WriteHeader(http.StatusOK)
		return
	}
	return
}

//scan qrcode show page
func scanqrcodeHandler(w http.ResponseWriter, r *http.Request) {
	prefix := fmt.Sprintf("[%s]", "scanQrcodeHandler")
	glog.Infof("%s %s %s", prefix, r.Method, r.RequestURI)

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		glog.Errorf("%s cannot parse url querystring, err %s", prefix, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	code := queryValues.Get("code")
	state := queryValues.Get("state")

	//request isn't redirected by wechat, return notice page
	if code == "" || state == "" {
		glog.Infof("%s request isn't redirect by weixin, return notice page")
		w.WriteHeader(http.StatusOK)
		return
	}

	var openId string
	var sid string
	var validSession bool = false

	cookie, err := r.Cookie("sid")
	//no cookie sid
	if err == nil {
		glog.Infof("%s get cookie sid %s", prefix, cookie.Value)
		openId, err = parseSession(cookie.Value)
		if err == nil {
			validSession = true
			glog.Infof("%s session is valid, user openid %s", prefix, openId)
		}
	} else {
		glog.Infof("%s no cookie sid", prefix, cookie.Value)
	}

	//cookie is not exist or invalid cookie
	if !validSession {
		glog.Infof("%s cookie not exist or invalid cookie, generating newOne", prefix)
		//get openid
		openId, err = getUserOpenId(code)
		if err != nil {
			glog.Errorf("%s cannot get user openid from wechat, err %s", prefix, err)
			return
		}
		glog.Infof("%s get user openid %s", prefix, openId)

		//set session
		sid, err = newSession(openId)
		if err != nil {
			glog.Errorf("%s cannot get new session, err %s", prefix, err)
			return
		}

		cookie := http.Cookie{
			Name:     "sid",
			Value:    sid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		glog.Infof("%s new session create for user with openid %s", prefix, openId)
	}

	user := User{}
	db.Where("wx_openid = ?", openId).First(&user)
	if user.ID == 0 {
		// openid isn't related to a user
		glog.Infof("%s openid is not related to a user", prefix)
		msgbody := "在菜单中点击绑定企业，绑定企业后再进行拍照，谢谢"
		n := NoticePage{Title: "扫描监控地点二维码", Type: noticePagefail, Msgtitle: "用户未绑定企业", Msgbody: msgbody}
		noticepageTmpl := template.Must(template.New("noticepage").Parse(myTemplate.NOTICEPAGE))
		noticepageTmpl.Execute(w, n)
		return
	}

	company := Company{}
	db.First(&company, user.CompanyId)

	// openid is related to a user
	nonceStr := string(rand.NewHex())
	timeNow := time.Now().Unix()
	ticket, err := ticketServer.Ticket()
	//cannot get ticket
	if err != nil {
		errmsg := fmt.Sprintf("cannot get ticket %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		return
	}

	var jssdkObj struct {
		Timestamp string
		Noncestr  string
		Wxappid   string
		Signature string
		Title     string
		User      string
		Company   string
		Phone     string
	}
	jssdkObj.Timestamp = fmt.Sprintf("%d", timeNow)
	jssdkObj.Noncestr = nonceStr
	jssdkObj.Wxappid = wxAppId
	jssdkObj.Signature = jssdk.WXConfigSign(ticket, nonceStr, jssdkObj.Timestamp, fmt.Sprintf("https://%s%s", domain, r.URL))
	jssdkObj.User = user.Name
	jssdkObj.Company = company.Name
	jssdkObj.Phone = user.Phone

	glog.Infof("%s get ticket %s, signature %s, noncestr %s, uri %s", prefix, ticket, jssdkObj.Signature, nonceStr, r.URL)
	glog.Infof("%s user %s, phone %s, company %s", prefix, jssdkObj.User, jssdkObj.Phone, jssdkObj.Company)
	scanqrcodeTmpl := template.Must(template.New("scanqrcode").Parse(myTemplate.SCANQRCODE))
	scanqrcodeTmpl.Execute(w, jssdkObj)
	return
}

func photoHandler(w http.ResponseWriter, r *http.Request) {
	prefix := fmt.Sprintf("[%s]", "photoHandler")
	glog.Infof("%s %s %s", prefix, r.Method, r.RequestURI)
	glog.Infof("%s user binding started")

	redirectURI := fmt.Sprintf("https://%s/backend/scanqrcode", domain)
	scope := "snsapi_base"
	state := "scanqrcode"
	authCodeURL := mpoauth2.AuthCodeURL(wxAppId, redirectURI, scope, state)

	cookie, err := r.Cookie("sid")
	if err != nil {
		//no cookie sid => request is not acccessed from wechat
		http.Redirect(w, r, authCodeURL, http.StatusFound)
		glog.Errorf("%s no cookie sid in request, return with redirect", prefix)
		return
	}

	openId, err := parseSession(cookie.Value)
	if err != nil {
		//cookie is invalid, redirect to start page
		http.Redirect(w, r, authCodeURL, http.StatusFound)
		glog.Errorf("%s cookie is invalid, redirect to /backend/scanqrcode", prefix)
		return
	}

	user := User{}
	db.Where("wx_openid = ?", openId).First(&user)
	//openid is not related to a user
	if user.ID == 0 {
		http.Redirect(w, r, authCodeURL, http.StatusFound)
		glog.Errorf("%s openid %s is not related to user, redirect to /backend/scanqrcode", prefix, openId)
		return
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		glog.Errorf("%s cannot parse url querystring, err %s", prefix, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	place := queryValues.Get("place")
	monitor_place := MonitorPlace{}
	db.Where("id = ?", place).First(&monitor_place)
	if monitor_place.ID == 0 {
		//cannot find place
		glog.Errorf("%s cannot find monitor_place with id %s", place)
		w.WriteHeader(http.StatusOK)
		return
	}
	if monitor_place.CompanyId != user.CompanyId {
		glog.Errorf("%s place %d belongs to company %d, but user %d belongs to company %d",
			prefix, monitor_place.ID, monitor_place.CompanyId, user.ID, user.CompanyId)
		w.WriteHeader(http.StatusOK)
		return
	}

	nonceStr := string(rand.NewHex())
	timeNow := time.Now().Unix()
	ticket, err := ticketServer.Ticket()
	//cannot get ticket
	if err != nil {
		errmsg := fmt.Sprintf("cannot get ticket %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		return
	}

	var jssdkObj struct {
		Timestamp string
		Noncestr  string
		Wxappid   string
		Signature string
		Userid    int
		Placeid   int
	}
	jssdkObj.Timestamp = fmt.Sprintf("%d", timeNow)
	jssdkObj.Noncestr = nonceStr
	jssdkObj.Wxappid = wxAppId
	jssdkObj.Signature = jssdk.WXConfigSign(ticket, nonceStr, jssdkObj.Timestamp, fmt.Sprintf("https://%s%s", domain, r.URL))

	glog.Infof("%s get ticket %s, signature %s, noncestr %s, uri %s", prefix, ticket, jssdkObj.Signature, nonceStr, r.URL)
	photoTmpl := template.Must(template.New("photo").Parse(myTemplate.PHOTO))
	photoTmpl.Execute(w, jssdkObj)
	return
}

// URL: /backend/download
// Used for download image from backend server
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	prefix := fmt.Sprintf("[%s]", "photoHandler")
	glog.Infof("%s %s %s", prefix, r.Method, r.RequestURI)

	redirectURI := fmt.Sprintf("https://%s/backend/scanqrcode", domain)
	scope := "snsapi_base"
	state := "scanqrcode"
	authCodeURL := mpoauth2.AuthCodeURL(wxAppId, redirectURI, scope, state)

	cookie, err := r.Cookie("sid")
	if err != nil {
		//no cookie sid => request is not acccessed from wechat
		http.Redirect(w, r, authCodeURL, http.StatusFound)
		glog.Errorf("%s no cookie sid in request, return with redirect", prefix)
		return
	}

	openId, err := parseSession(cookie.Value)
	if err != nil {
		//cookie is invalid, redirect to start page
		glog.Errorf("%s cookie is invalid, redirect to /backend/binding", prefix)
		return
	}

	//openid is related to a user
	if isUserExist(openId) {
		glog.Infof("%s openid is related to a user, return notice page", prefix)
		return
	}

	r.ParseForm()
	userId := r.Form.Get("userId")
	placeId := r.Form.Get("placeId")
	serverId := r.Form.Get("serverId")
	glog.Infof("%s userId %s, placeId %s, serverId %s", userId, placeId, serverId)

	var response struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}

	user := User{}
	db.Where("id = ?", userId).First(&user)
	//openid is not related to a user
	if user.ID == 0 {
		errmsg := fmt.Sprintf("invalid user id %s, user not found", userId)
		glog.Errorf("%s %s", prefix, errmsg)
		response.Status = false
		response.Message = "无效用户,请使用有效用户上传"
		returnContent, err := json.Marshal(response)
		if err != nil {
			glog.Errorf("%s json marshal error %s, response %v", prefix, err, response)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(returnContent))
		w.WriteHeader(http.StatusOK)
		return
	}

	monitor_place := MonitorPlace{}
	db.Where("id = ?", placeId).First(&monitor_place)
	if monitor_place.ID == 0 {
		errmsg := fmt.Sprintf("invalid monitor_place id %s, monitor_place not found", placeId)
		glog.Errorf("%s %s", prefix, errmsg)
		response.Status = false
		response.Message = "无效监控地点,请拍摄有效监控地点"
		returnContent, err := json.Marshal(response)
		if err != nil {
			glog.Errorf("%s json marshal error %s, response %v", prefix, err, response)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(returnContent))
		w.WriteHeader(http.StatusOK)
		return
	}

	if user.CompanyId != monitor_place.CompanyId {
		errmsg := fmt.Sprintf("monitor_place belong to company %d, user belong to company %d, unmatch", user)
		glog.Errorf("%s %s", prefix, errmsg)
		response.Status = false
		response.Message = "用户对该地点无权限,请使用该地点的管理人员账户上传"
		returnContent, err := json.Marshal(response)
		if err != nil {
			glog.Errorf("%s json marshal error %s, response %v", prefix, err, response)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(returnContent))
		w.WriteHeader(http.StatusOK)
		return
	}

	timeNow := time.Now()
	timeToday := fmt.Sprintf("%s%02s%02s", timeNow.Year(), timeNow.Month(), timeNow.Day())

	picture := Picture{CreateAt: timeNow, UpdateAt: timeNow, MonitorPlaceId: monitor_place.ID, UserId: user.ID}
	picture.ThumbPath = fmt.Sprintf("/picture/%s/%d/%d/thumb_%s.png", timeToday, monitor_place.CompanyId, monitor_place.ID, picture.ID)
	picture.ThumbURI = fmt.Sprintf("/picture/%s/%d/%d/thumb_%s.png", timeToday, monitor_place.CompanyId, monitor_place.ID, picture.ID)
	picture.FullPath = fmt.Sprintf("/static/picture/%s/%d/%d/full_%s.png", timeToday, monitor_place.CompanyId, monitor_place.ID, picture.ID)
	picture.FullURI = fmt.Sprintf("/static/picture/%s/%d/%d/full_%s.png", timeToday, monitor_place.CompanyId, monitor_place.ID, picture.ID)

	//picture save place: /picture/20171206/<monitor_place.CompanyId>/<monitor_place.ID>/full_<picture_id>.png
	written, err := media.Download(wechatClient, serverId, fmt.Sprintf("%s/picture/%s/%d/%d/full_%s.png", imgRepo, monitor_place.CompanyId, timeToday, monitor_place.ID))
	if err != nil {
		errmsg := fmt.Sprintf("cannot download media id %s for place %s", serverId, placeId)
		response.Status = false
		response.Message = "上传失败,请重新上传"
		returnContent, err := json.Marshal(response)
		if err != nil {
			glog.Errorf("%s json marshal error %s, response %v", prefix, err, response)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		glog.Errorf("%s %s", prefix, errmsg)
		io.WriteString(w, string(returnContent))
		w.WriteHeader(http.StatusOK)
		return
	}

	response.Status = true
	response.Message = "上传成功"
	returnContent, err := json.Marshal(response)
	glog.Infof("%s download serverId %s success, bytes %d", prefix, serverId, written)
	if err != nil {
		glog.Errorf("%s json marshal error %s, response %v", prefix, err, response)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(returnContent))
	w.WriteHeader(http.StatusOK)
	return
}