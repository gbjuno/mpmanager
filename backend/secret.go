package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

	session, err := sessionStorage.Get(cookieValue)
	if err != nil {
		glog.Errorf("%s session is outdate or invalid, err %s", prefix, err)
		return "", err
	}
	// user session is valid
	return session.(string), nil
}

func PasswordAuthenticate(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	sessionid := request.HeaderParameter("sessionid")
	if sessionid == "" {
		response.WriteHeaderAndEntity(http.StatusUnauthorized, Response{Status: "error", Error: "401 please logi"})
		return
	}
	user_id, err := parsePasswordSession(sessionid)
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusUnauthorized, Response{Status: "error", Error: err.Error()})
		return
	}
	id, _ := strconv.Atoi(user_id)
	user := User{}
	db.Debug().Where("id = ?", id).First(&user)
	if user.ID == 0 {
		response.WriteHeaderAndEntity(http.StatusUnauthorized, Response{Status: "error", Error: "401 user not found"})
		return
	}

	chain.ProcessFilter(request, response)
}
