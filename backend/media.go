package main

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"gopkg.in/chanxuehong/wechat.v2/mp/base"
)

type MediaUrl struct {
	Url string `json:"url"`
}

func (mp MediaPicture) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/mediapicture").Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.POST("").To(mp.uploadPicture))
	container.Add(ws)
}

func (mp MediaPicture) uploadPicture(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [uploadPicture]", request.Request.RemoteAddr)
	glog.Infof("%s POST %s", prefix, request.Request.URL)

	request.Request.ParseMultipartForm(32 << 20)

	file, handler, err := request.Request.FormFile("uploadImage")
	if err != nil {
		errmsg := fmt.Sprintf("cannot read image, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: "无法上传图片,请稍后再试"})
		return
	}

	mediaUrl := MediaUrl{}
	mediaUrl.Url, err = base.UploadImageFromReader(wechatClient, handler.Filename, file)
	if err != nil {
		errmsg := fmt.Sprintf("无法上传图片, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, &mediaUrl)
	glog.Infof("%s upload media picture success, url %s", prefix, mediaUrl.Url)
	return
}
