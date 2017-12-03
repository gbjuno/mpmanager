package main

import (
	"fmt"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/jssdk"
	"gopkg.in/chanxuehong/wechat.v2/mp/menu"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/callback/request"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/callback/response"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/sid"
	"github.com/emicklei/go-restful"
	myTemplate "github.com/gbjuno/mpmanager/backend/templates"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"gopkg.in/chanxuehong/wechat.v2/oauth2"
)

type WxResponse struct {
}

func (w WxResponse) Register(container *restful.Container) {
	ws := new(restful.WebService)
	container.Add(ws)
}

const (
	wxAppId     = "wx6eb571f36f6b1c10"
	wxAppSecret = "555210c557802c8a0c6a930cf2e4c159"

	wxOriId         = "gh_75a8e6a73da5"
	wxToken         = "wechatcms"
	wxEncodedAESKey = "aeskeyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
)

var (
	sessionStorage                           = session.New(20*60, 60*60)
	oauth2Endpoint    oauth2.Endpoint        = mpoauth2.NewEndpoint(wxAppId, wxAppSecret)
	accessTokenServer core.AccessTokenServer = core.NewDefaultAccessTokenServer(wxAppId, wxAppSecret, nil)
	wechatClient      *core.Client           = core.NewClient(accessTokenServer, nil)
	ticketServer                             = jssdk.NewDefaultTicketServer(wechatClient)

	// 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
	msgHandler core.Handler
	msgServer  *core.Server
)

func createMenu() {
	bindButton := &menu.Button{}
	bind_redirect_url := url.PathEscape("https://www.juntengshoes.cn/backend/session")
	bind_url := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=test#wechat_redirect", wxAppId, bind_redirect_url)
	bindButton.SetAsViewButton("绑定账户", bind_url)

	photoButton := &menu.Button{}
	photo_redirect_url := url.PathEscape("https://www.juntengshoes.cn/backend/photo")
	photo_url := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=test#wechat_redirect", wxAppId, photo_redirect_url)
	photoButton.SetAsViewButton("拍  照", photo_url)

	checkButton := &menu.Button{}
	checkButton.SetAsClickButton("check", "CL_01")
	wechatMenu := &menu.Menu{Buttons: []menu.Button{*bindButton, *photoButton, *checkButton}}

	if err := menu.Create(wechatClient, wechatMenu); err != nil {
		log.Fatalf("cannot connect with weixin server, err %s\n", err)
	}
	log.Printf("create menu succeed")
}

func init() {
	if wechatMenu, _, err := menu.Get(wechatClient); err != nil {
		log.Printf("cannot get menu, err %s", err)
		if wechatMenu == nil || len(wechatMenu.Buttons) != 3 {
			createMenu()
		}
	} else {
		log.Printf("Menu: %s", wechatMenu)
		if err := menu.Delete(wechatClient); err != nil {
			log.Fatalf("cannot delete menu, err %s", err)
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

func textMsgHandler(ctx *core.Context) {
	log.Printf("收到文本消息:\n%s\n", ctx.MsgPlaintext)

	msg := request.GetText(ctx.MixedMsg)
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, msg.Content)
	ctx.RawResponse(resp) // 明文回复
	//ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultMsgHandler(ctx *core.Context) {
	log.Printf("收到消息:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func menuClickEventHandler(ctx *core.Context) {
	log.Printf("收到菜单 click 事件:\n%s\n", ctx.MsgPlaintext)

	event := menu.GetClickEvent(ctx.MixedMsg)
	resp := response.NewText(event.FromUserName, event.ToUserName, event.CreateTime, "收到 click 类型的事件")
	ctx.RawResponse(resp) // 明文回复
	//ctx.AESResponse(resp, 0, "", nil) // aes密文回复
}

func defaultEventHandler(ctx *core.Context) {
	log.Printf("收到事件:\n%s\n", ctx.MsgPlaintext)
	ctx.NoneResponse()
}

func init() {
	http.HandleFunc("/wx_callback", wxCallbackHandler)
	http.HandleFunc("/session", sessionHandler)
	http.HandleFunc("/bind", bindingHandler)
	http.HandleFunc("/confirm", confirmHandler)
	//http.HandleFunc("/picture", pictureHandler)
	http.HandleFunc("/photo", photoHandler)
}

// wxCallbackHandler 是处理回调请求的 http handler.
//  1. 不同的 web 框架有不同的实现
//  2. 一般一个 handler 处理一个公众号的回调请求(当然也可以处理多个, 这里我只处理一个)
func wxCallbackHandler(w http.ResponseWriter, r *http.Request) {
	msgServer.ServeHTTP(w, r, nil)
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("session handler")
	sid := sid.New()
	state := string(rand.NewHex())

	if err := sessionStorage.Add(sid, state); err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	cookie := http.Cookie{
		Name:     "sid",
		Value:    sid,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	oauth2RedirectURI := "https://www.juntengshoes.cn/backend/bind"
	oauth2Scope := "snsapi_base"

	AuthCodeURL := mpoauth2.AuthCodeURL(wxAppId, oauth2RedirectURI, oauth2Scope, state)
	log.Println("AuthCodeURL:", AuthCodeURL)

	http.Redirect(w, r, AuthCodeURL, http.StatusFound)
}

func bindingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)

	cookie, err := r.Cookie("sid")
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	session, err := sessionStorage.Get(cookie.Value)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	savedState := session.(string) // 一般是要序列化的, 这里保存在内存所以可以这么做

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	code := queryValues.Get("code")
	if code == "" {
		log.Println("用户禁止授权")
		return
	}

	queryState := queryValues.Get("state")
	if queryState == "" {
		log.Println("state 参数为空")
		return
	}
	if savedState != queryState {
		str := fmt.Sprintf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, queryState)
		io.WriteString(w, str)
		log.Println(str)
		return
	}

	oauth2Client := oauth2.Client{
		Endpoint: oauth2Endpoint,
	}

	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	log.Printf("token: %+v\r\n", token)
	log.Printf("openid: %s", token.OpenId)

	bind_tmpl := template.Must(template.New("bind").Parse(myTemplate.BIND))
	bind_tmpl.Execute(w, nil)

	return
}

func confirmHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("confirmHandler")
	cookie, err := r.Cookie("sid")
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}

	_, err = sessionStorage.Get(cookie.Value)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Println(err)
		return
	}
	r.ParseForm()
	log.Printf("form value %s\n", r.Form)
	phone := r.Form.Get("phone")
	password := r.Form.Get("password")

	log.Printf("phone %s\n", phone)
	log.Printf("password %s\n", password)

	if phone == "123456" && password == "123456" {
		io.WriteString(w, "YES")
	} else {
		io.WriteString(w, "NO")
	}
}

/*
func pictureHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)

	str1 := "abcdefghijklmnopqrstuvwxyz"
	time1 := time.Now().Unix()

	ticket, err := ticketServer.Ticket()
	if err != nil {
		log.Fatalf("cannot get ticket %s", err)
	}

	var obj1 struct {
		Timestamp string
		Noncestr  string
		Wxappid   string
		Signature string
	}

	obj1.Timestamp = fmt.Sprintf("%d", time1)
	obj1.Noncestr = str1
	obj1.Wxappid = wxAppId
	obj1.Signature = jssdk.WXConfigSign(ticket, str1, obj1.Timestamp, fmt.Sprintf("http://www.juntengshoes.cn%s", r.URL))

	log.Printf("%s", ticket)
	log.Printf("%s", obj1.Signature)
	log.Printf("http://www.juntengshoes.cn%s", r.URL)

	picture_tmpl := template.Must(template.New("picture").Parse(myTemplate.Picture_html))
	log.Println("send html to weixin")
	picture_tmpl.Execute(w, obj1)
	picture_tmpl.Execute(os.Stdout, obj1)

	return
}
*/
func photoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)

	str1 := "abcdefghijklmnopqrstuvwxyz"
	time1 := time.Now().Unix()

	ticket, err := ticketServer.Ticket()
	if err != nil {
		log.Fatalf("cannot get ticket %s", err)
	}

	var obj1 struct {
		Timestamp string
		Noncestr  string
		Wxappid   string
		Signature string
	}

	obj1.Timestamp = fmt.Sprintf("%d", time1)
	obj1.Noncestr = str1
	obj1.Wxappid = wxAppId
	obj1.Signature = jssdk.WXConfigSign(ticket, str1, obj1.Timestamp, fmt.Sprintf("http://www.juntengshoes.cn%s", r.URL))

	log.Printf("%s", ticket)
	log.Printf("%s", obj1.Signature)
	log.Printf("https://www.juntengshoes.cn%s", r.URL)

	photo_tmpl := template.Must(template.New("photo").Parse(myTemplate.PHOTO))
	log.Println("send html to weixin")
	photo_tmpl.Execute(w, obj1)
	photo_tmpl.Execute(os.Stdout, obj1)

	return
}

func main() {
	log.Println(http.ListenAndServe(":8000", nil))
}