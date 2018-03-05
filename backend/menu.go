package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"gopkg.in/chanxuehong/wechat.v2/mp/menu"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
)

type Menu struct{}

func (m Menu) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/menu").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(m.getMenu))
	ws.Route(ws.POST("").To(m.updateMenu))
	container.Add(ws)
	glog.Infof("register menu service")
}

func (m Menu) getMenu(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] getMenu", request.Request.RemoteAddr)
	wechatMenu, _, err := menu.Get(wechatClient)
	if err != nil {
		errmsg := fmt.Sprintf("%s cannot get menu, err %s", prefix, err)
		glog.Errorf(errmsg)
		r := Response{Status: "error", Error: "无法连接腾讯服务器,请稍后重试"}
		response.WriteHeaderAndEntity(http.StatusInternalServerError, &r)
		return
	}
	glog.Infof("%s get wechat menu successfully", prefix)
	response.WriteHeaderAndEntity(http.StatusOK, &wechatMenu)
	return
}

func (m Menu) updateMenu(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] updateMenu", request.Request.RemoteAddr)
	if request.Request.Method == "GET" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	wechatMenu := menu.Menu{}
	err := request.ReadEntity(&wechatMenu)
	if err != nil {
		glog.Errorf("%s cannot parse wechat menu entity", prefix)
		r := Response{Status: "error", Error: "无法解析"}
		response.WriteHeaderAndEntity(http.StatusInternalServerError, &r)
		return
	}

	if err := menu.Delete(wechatClient); err != nil {
		glog.Errorf("%s cannot delete menu, err %s", prefix, err)
		r := Response{Status: "error", Error: "无法更新菜单"}
		response.WriteHeaderAndEntity(http.StatusInternalServerError, &r)
		return
	}

	disasterCheckButton := getDisasterCheckButton()

	length := len(wechatMenu.Buttons)
	if length == 0 {
		wechatMenu.Buttons = []menu.Button{*disasterCheckButton}
	} else {
		setDisasterCheckButton := false
		for key, b := range wechatMenu.Buttons {
			if b.Name == "隐患排查" {
				setDisasterCheckButton = true
				wechatMenu.Buttons[key] = *disasterCheckButton
			}
		}
		if !setDisasterCheckButton {
			wechatMenu.Buttons[length-1] = *disasterCheckButton
		}
	}

	if err := menu.Create(wechatClient, &wechatMenu); err != nil {
		glog.Errorf("%s cannot connect with weixin server, err %s", prefix, err)
		r := Response{Status: "error", Error: "无法更新菜单,请联系管理员"}
		response.WriteHeaderAndEntity(http.StatusInternalServerError, &r)
		return
	}

	response.WriteHeader(http.StatusOK)
	glog.Info("%s update menu succeed")
	return
}

func getDisasterCheckButton() *menu.Button {
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
	scanqrcodeButton.SetAsViewButton("扫描拍照", scanqrcodeURI)
	glog.Infof("set Button scanqrcode for uri %s, wechat redirecturi %s", scanqrcodeRedirectURI, scanqrcodeURI)

	companystatButton := &menu.Button{}
	companystatRedirectURI := fmt.Sprintf(redirectURIPrefix, "companystat")
	companystatState := "companystat"
	companystatURI := mpoauth2.AuthCodeURL(wxAppId, companystatRedirectURI, oauth2Scope, companystatState)
	companystatButton.SetAsViewButton("检查进度", companystatURI)
	glog.Infof("set Button companystat for uri %s, wechat redirecturi %s", companystatRedirectURI, companystatURI)

	topButton := &menu.Button{}
	topButton.Name = "隐患排查"
	topButton.SubButtons = []menu.Button{*bindingButton, *scanqrcodeButton, *companystatButton}

	return topButton
}
