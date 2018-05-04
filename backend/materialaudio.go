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

type MaterialAudioList struct {
	Count          int             `json:"count"`
	MaterialAudios []MaterialAudio `json:"audios"`
}

func (mp MaterialAudio) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/materialaudio").Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(mp.findAudio))
	ws.Route(ws.GET("/{audio_id}").To(mp.findAudio))
	ws.Route(ws.POST("").To(mp.uploadAudio))
	ws.Route(ws.DELETE("/{audio_id}").To(mp.deleteAudio))
	container.Add(ws)
}

func (mp MaterialAudio) findAudio(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findAudio]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	audio_id := request.PathParameter("audio_id")

	if audio_id == "" {
		materialAudioList := MaterialAudioList{}
		materialAudioList.MaterialAudios = make([]MaterialAudio, 0)
		db.Debug().Find(&materialAudioList.MaterialAudios)
		materialAudioList.Count = len(materialAudioList.MaterialAudios)
		response.WriteHeaderAndEntity(http.StatusOK, materialAudioList)
		glog.Infof("%s return material list", prefix)
		return
	}

	id, err := strconv.Atoi(audio_id)
	//fail to parse country id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get material audio, audio_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	materialAudio := MaterialAudio{}
	db.Debug().First(&materialAudio, id)
	if materialAudio.ID == 0 {
		errmsg := fmt.Sprintf("cannot find material audio with id %s", audio_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	glog.Infof("%s return material audio with id %s", prefix, audio_id)
	response.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(response.ResponseWriter)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(&materialAudio)
	return
}

func (mp MaterialAudio) uploadAudio(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [uploadAudio]", request.Request.RemoteAddr)
	glog.Infof("%s POST %s", prefix, request.Request.URL)

	request.Request.ParseMultipartForm(32 << 20)
	title := request.Request.FormValue("title")
	introduction := request.Request.FormValue("introduction")
	url := request.Request.FormValue("url")

	glog.Infof("%s title %s, introduction %s, url %s", prefix, title, introduction, url)
	if url != "" {
		materialAudio := MaterialAudio{}
		materialAudio.Title = title
		materialAudio.Introduction = introduction
		materialAudio.Url = url
		db.Debug().Create(&materialAudio)
		response.WriteHeaderAndEntity(http.StatusOK, &materialAudio)
		glog.Infof("%s upload anonymous material audio success")
		return
	}

	file, handler, err := request.Request.FormFile("uploadAudio")
	if err != nil {
		errmsg := fmt.Sprintf("cannot read audio, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: "无法上传音频,请稍后再试"})
		return
	}

	old := wechatClient.HttpClient
	client := *http.DefaultClient
	client.Timeout = time.Second * 300
	wechatClient.HttpClient = &client
	mediaID, err := material.UploadVoiceFromReader(wechatClient, handler.Filename, file)
	if err != nil {
		errmsg := ""
		errmsg = fmt.Sprintf("无法上传音频，音频格式不是mp4或大于10M. err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
	wechatClient.HttpClient = old

	materialAudio := MaterialAudio{}
	materialAudio.Title = title
	materialAudio.Introduction = introduction
	materialAudio.MediaId = mediaID

	db.Debug().Create(&materialAudio)
	response.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(response.ResponseWriter)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(&materialAudio)
	glog.Infof("%s upload material audio success")
	return
}

func (mp MaterialAudio) deleteAudio(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteAudio]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	audio_id := request.PathParameter("audio_id")
	id, err := strconv.Atoi(audio_id)
	//fail to parse country id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete material audio, audio_id is not integer, err %s", err)
		returnmsg := fmt.Sprintf("无法删除音频，提供的id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	materialAudio := MaterialAudio{}
	db.Debug().First(&materialAudio, id)
	if materialAudio.ID == 0 {
		//country with id doesn't exist
		glog.Infof("%s materialAudio with id %s doesn't exist, return ok", prefix, audio_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	if materialAudio.MediaId != "" {
		if err := material.Delete(wechatClient, materialAudio.MediaId); err != nil {
			glog.Errorf("%s materialAudio with id %s cannot delete", prefix, audio_id)
			response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "error", Error: "无法删除音频,请稍后再试"})
			return
		}
	}

	db.Debug().Delete(&materialAudio)
	glog.Infof("%s delete materialAudio with id %s successfully", prefix, audio_id)
	response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	return
}
