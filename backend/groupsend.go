package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/mass/mass2all"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/mass/preview"
)

type GroupSendList struct {
	Count     int         `json:"count"`
	GroupSend []GroupSend `json:"groupsend"`
}

const PREVIEW_USER = "o1k8S0nPewTiG3vE6ZSl_1pQLDWA"

func (g GroupSend) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/groupsend").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(g.findGroupSend))
	ws.Route(ws.GET("/?pageNo={pageNo}&pageSize={pageSize}&order={order}").To(g.findGroupSend))
	ws.Route(ws.GET("/{groupSend_id}").To(g.findGroupSend))
	ws.Route(ws.POST("").To(g.createGroupSend))
	ws.Route(ws.POST("/?preview={preview}").To(g.createGroupSend))
	container.Add(ws)
}

func (g GroupSend) findGroupSend(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findGroupSend]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	groupSend_id := request.PathParameter("groupSend_id")
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")

	var searchGroupSend *gorm.DB = db.Debug()

	if order != "asc" && order != "desc" {
		errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
		glog.Errorf("%s %s", prefix, errmsg)
		order = "desc"
	}

	if order == "" {
		order = "desc"
	}

	glog.Infof("%s find groupSend with order %s", prefix, order)

	groupSends := make([]GroupSend, 0)
	count := 0
	searchGroupSend.Find(&groupSends).Count(&count)
	searchGroupSend = searchGroupSend.Order("id " + order)

	if groupSend_id == "" {
		isPageSizeOk := true
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			isPageSizeOk = false
			errmsg := fmt.Sprintf("cannot find groupSend with pageSize %s, err %s, ignore", pageSize, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		//pageNo depends on pageSize
		isPageNoOk := true
		pageNoInt, err := strconv.Atoi(pageNo)
		if err != nil {
			isPageNoOk = false
			errmsg := fmt.Sprintf("cannot find groupSend with pageNo %s, err %s, ignore", pageNo, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		if isPageSizeOk && isPageNoOk {
			limit := pageSizeInt
			offset := (pageNoInt - 1) * limit
			glog.Infof("%s set find groupSend db with pageSize %s, pageNo %s(limit %d, offset %d)", prefix, pageSize, pageNo, limit, offset)
			searchGroupSend = searchGroupSend.Offset(offset).Limit(limit)
		}

		groupSendList := GroupSendList{}
		groupSendList.GroupSend = make([]GroupSend, 0)
		searchGroupSend.Find(&groupSendList.GroupSend)

		response.WriteHeaderAndEntity(http.StatusOK, &groupSendList)
		glog.Infof("%s return groupSend list", prefix)
		return
	}

	id, err := strconv.Atoi(groupSend_id)
	//fail to parse groupSend id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get groupSend, groupSend_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	groupSend := GroupSend{}
	db.Debug().First(&groupSend, id)
	//cannot find groupSend
	if groupSend.ID == 0 {
		errmsg := fmt.Sprintf("cannot find groupSend with id %s", groupSend_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, &groupSend)
	glog.Infof("%s find groupSend with id %d", prefix, groupSend.ID)
	return
}

func (g GroupSend) createGroupSend(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createGroupSend]", request.Request.RemoteAddr)
	param_preview := request.QueryParameter("preview")
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	groupSend := GroupSend{}
	err := request.ReadEntity(&groupSend)
	if err != nil {
		errmsg := fmt.Sprintf("cannot create groupSend, err %s", err)
		returnmsg := fmt.Sprintf("无法群发消息,提供的信息错误,请联系管理员")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	news := News{}
	db.Debug().First(&news, groupSend.NewsID)
	groupSend.NewsName = news.Name
	groupSend.MediaId = news.MediaId

	if param_preview == "true" {
		sendNews := preview.NewNews(PREVIEW_USER, groupSend.MediaId)
		err := preview.Send(wechatClient, sendNews)
		if err != nil {
			errmsg := fmt.Sprintf("cannot create groupSend, err %s", err)
			returnmsg := fmt.Sprintf("无法预览消息,请重试")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		}
		glog.Infof("%s preview groupSend with id %d succesfully", prefix, groupSend.ID)
		response.WriteHeader(http.StatusOK)
		return
	}

	sendNews := mass2all.NewNews(groupSend.MediaId)
	result, err := mass2all.Send(wechatClient, sendNews)
	if err != nil {
		errmsg := fmt.Sprintf("cannot create groupSend, err %s", err)
		returnmsg := fmt.Sprintf("无法群发消息,请重试")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	groupSend.MsgId = result.MsgId
	groupSend.MsgDataId = result.MsgDataId
	db.Debug().Create(&groupSend)
	glog.Info("%s create groupSend with id %d succesfully", prefix, groupSend.ID)
	response.WriteHeader(http.StatusOK)
	return
}
