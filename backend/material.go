package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"gopkg.in/chanxuehong/wechat.v2/mp/material"
)

type MaterialPictureList struct {
	Count            int               `json:"count"`
	MaterialPictures []MaterialPicture `json:"pictures"`
}

func (mp MaterialPicture) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/materialpicture").Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(mp.findPicture))
	ws.Route(ws.GET("/{picture_id}").To(mp.findPicture))
	ws.Route(ws.POST("").To(mp.uploadPicture))
	ws.Route(ws.DELETE("/{picture_id}").To(mp.deletePicture))
	container.Add(ws)
}

func (mp MaterialPicture) findPicture(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findPicture]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	picture_id := request.PathParameter("picture_id")

	if picture_id == "" {
		materialPictureList := MaterialPictureList{}
		materialPictureList.MaterialPictures = make([]MaterialPicture, 0)
		db.Debug().Find(&materialPictureList.MaterialPictures)
		materialPictureList.Count = len(materialPictureList.MaterialPictures)
		response.WriteHeaderAndEntity(http.StatusOK, materialPictureList)
		glog.Infof("%s return material list", prefix)
		return
	}

	id, err := strconv.Atoi(picture_id)
	//fail to parse country id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get material picture, picture_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	materialPicture := MaterialPicture{}
	db.Debug().First(&materialPicture, id)
	if materialPicture.ID == 0 {
		errmsg := fmt.Sprintf("cannot find material picture with id %s", picture_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	glog.Infof("%s return material picture with id %s", prefix, picture_id)
	response.WriteHeaderAndEntity(http.StatusOK, &materialPicture)
	return
}

func (mp MaterialPicture) uploadPicture(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [uploadPicture]", request.Request.RemoteAddr)
	glog.Infof("%s POST %s", prefix, request.Request.URL)
	media_id, url, err := material.UploadImageFromReader(wechatClient, "title", request.Request.Body)
	if err != nil {
		errmsg := fmt.Sprintf("无法上传图片, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	materialPicture := MaterialPicture{}
	materialPicture.MediaId = media_id
	materialPicture.Url = url

	db.Debug().Create(&materialPicture)
	response.WriteHeaderAndEntity(http.StatusOK, &materialPicture)
	glog.Info("%s upload material picture success")
	return
}

func (mp MaterialPicture) deletePicture(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deletePicture]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	picture_id := request.PathParameter("picture_id")
	id, err := strconv.Atoi(picture_id)
	//fail to parse country id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete material picture, picture_id is not integer, err %s", err)
		returnmsg := fmt.Sprintf("无法删除图片，提供的id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	materialPicture := MaterialPicture{}
	db.Debug().First(&materialPicture, id)
	if materialPicture.ID == 0 {
		//country with id doesn't exist
		glog.Infof("%s materialPicture with id %s doesn't exist, return ok", prefix, picture_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	if err := material.Delete(wechatClient, materialPicture.MediaId); err != nil {
		glog.Errorf("%s materialPicture with id %s cannot delete", prefix, picture_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "error", Error: "无法删除图片,请稍后再试"})
		return
	}

	db.Debug().Delete(&materialPicture)
	glog.Infof("%s delete materialPicture with id successfully", prefix, picture_id)
	response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	return
}
