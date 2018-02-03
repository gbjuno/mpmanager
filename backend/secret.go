package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/chanxuehong/session"
	"github.com/chanxuehong/sid"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

var PasswordSessionStorage = session.New(30*60, 30*60)

func newPasswordSession(user_id string) (string, error) {
	prefix := fmt.Sprintf("[%s]", "newPasswordSession")
	glog.Infof("%s user_id %s", prefix, user_id)
	sid := sid.New()
	if err := PasswordSessionStorage.Add(sid, user_id); err != nil {
		//fail to set session
		errmsg := fmt.Sprintf("cannot set sid to sessionStorage, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		return "", errors.New(errmsg)
	}
	return sid, nil
}

func parsePasswordSession(cookieValue string) (string, error) {
	prefix := fmt.Sprintf("[%s]", "parsePasswordSession")
	glog.Infof("%s parse cookie, cookie value %s", prefix, cookieValue)

	session, err := PasswordSessionStorage.Get(cookieValue)
	if err != nil {
		glog.Errorf("%s session is outdate or invalid, err %s", prefix, err)
		return "", err
	}
	// user session is valid
	return session.(string), nil
}

func PasswordAuthenticate(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	prefix := fmt.Sprintf("[%s]", "PasswordAuthenticate")
	cookieStr := request.HeaderParameter("Cookie")
	var sessionid string
	for _, cookie := range strings.Split(cookieStr, ";") {
		temp := strings.Split(cookie, "=")
		if temp[0] == "sessionid" {
			sessionid = temp[1]
		}
	}

	if sessionid == "" {
		glog.Infof("%s no sessionid", prefix)
		response.WriteHeaderAndEntity(http.StatusUnauthorized, Response{Status: "error", Error: "请登陆后进行操作"})
		return
	}

	_, err := parsePasswordSession(sessionid)
	if err != nil {
		glog.Infof("%s sessionid is not valid ,err %s", prefix, err)
		response.WriteHeaderAndEntity(http.StatusUnauthorized, Response{Status: "error", Error: "请重新登陆"})
		return
	}

	chain.ProcessFilter(request, response)
}
