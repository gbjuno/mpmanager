package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"gopkg.in/chanxuehong/wechat.v2/mp/material"
)

type MaterialVideoList struct {
	Count          int             `json:"count"`
	MaterialVideos []MaterialVideo `json:"videos"`
}

func (mp MaterialVideo) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/materialvideo").Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(mp.findVideo))
	ws.Route(ws.GET("/{video_id}").To(mp.findVideo))
	ws.Route(ws.POST("").To(mp.uploadVideo))
	ws.Route(ws.DELETE("/{video_id}").To(mp.deleteVideo))
	container.Add(ws)
}

func (mp MaterialVideo) findVideo(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findVideo]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	video_id := request.PathParameter("video_id")

	if video_id == "" {
		materialVideoList := MaterialVideoList{}
		materialVideoList.MaterialVideos = make([]MaterialVideo, 0)
		db.Debug().Find(&materialVideoList.MaterialVideos)
		materialVideoList.Count = len(materialVideoList.MaterialVideos)
		response.WriteHeaderAndEntity(http.StatusOK, materialVideoList)
		glog.Infof("%s return material list", prefix)
		return
	}

	id, err := strconv.Atoi(video_id)
	//fail to parse country id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get material video, video_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	materialVideo := MaterialVideo{}
	db.Debug().First(&materialVideo, id)
	if materialVideo.ID == 0 {
		errmsg := fmt.Sprintf("cannot find material video with id %s", video_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	glog.Infof("%s return material video with id %s", prefix, video_id)
	response.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(response.ResponseWriter)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(&materialVideo)
	return
}

func (mp MaterialVideo) uploadVideo(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [uploadVideo]", request.Request.RemoteAddr)
	glog.Infof("%s POST %s", prefix, request.Request.URL)

	request.Request.ParseMultipartForm(32 << 20)
	title := request.Request.FormValue("title")
	introduction := request.Request.FormValue("introduction")
	url := request.Request.FormValue("url")

	glog.Infof("%s title %s, introduction %s, url %s", prefix, title, introduction, url)
	if url != "" {
		materialVideo := MaterialVideo{}
		materialVideo.Title = title
		materialVideo.Introduction = introduction
		materialVideo.Url = url
		db.Debug().Create(&materialVideo)
		response.WriteHeaderAndEntity(http.StatusOK, &materialVideo)
		glog.Infof("%s upload anonymous material video success")
		return
	}

	file, handler, err := request.Request.FormFile("uploadVideo")
	if err != nil {
		errmsg := fmt.Sprintf("cannot read video, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: "无法上传视频,请稍后再试"})
		return
	}

	old := wechatClient.HttpClient
	client := *http.DefaultClient
	client.Timeout = time.Second * 300
	wechatClient.HttpClient = &client
	mediaID, err := material.UploadVideoFromReader(wechatClient, handler.Filename, file, title, introduction)
	if err != nil {
		errmsg := ""
		errmsg = fmt.Sprintf("无法上传视频，视频格式不是mp4或大于10M. err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
	wechatClient.HttpClient = old

	info, err := material.GetVideo(wechatClient, mediaID)
	if err != nil {
		errmsg := fmt.Sprintf("无法上传视频，err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	materialVideo := MaterialVideo{}
	materialVideo.Title = title
	materialVideo.Introduction = introduction
	materialVideo.MediaId = mediaID
	materialVideo.Url = info.DownloadURL

	db.Debug().Create(&materialVideo)
	response.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(response.ResponseWriter)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(&materialVideo)
	glog.Infof("%s upload material video success")
	return
}

func (mp MaterialVideo) deleteVideo(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteVideo]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	video_id := request.PathParameter("video_id")
	id, err := strconv.Atoi(video_id)
	//fail to parse country id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete material video, video_id is not integer, err %s", err)
		returnmsg := fmt.Sprintf("无法删除视频，提供的id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	materialVideo := MaterialVideo{}
	db.Debug().First(&materialVideo, id)
	if materialVideo.ID == 0 {
		//country with id doesn't exist
		glog.Infof("%s materialVideo with id %s doesn't exist, return ok", prefix, video_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	if materialVideo.MediaId != "" {
		if err := material.Delete(wechatClient, materialVideo.MediaId); err != nil {
			glog.Errorf("%s materialVideo with id %s cannot delete", prefix, video_id)
			response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "error", Error: "无法删除视频,请稍后再试"})
			return
		}
	}

	db.Debug().Delete(&materialVideo)
	glog.Infof("%s delete materialVideo with id %s successfully", prefix, video_id)
	response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	return
}
