package main

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
)

type MediaUrl struct {
	url string
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

	var incompleteURL = "https://api.weixin.qq.com/cgi-bin/media/uploadimg&access_token="
	var fields = []core.MultipartFormField{
		{
			IsFile:   true,
			Name:     "media",
			FileName: "picture",
			Value:    request.Request.Body,
		},
	}
	var result struct {
		core.Error
		MediaUrl
	}

	var err error
	if err = wechatClient.PostMultipartForm(incompleteURL, fields, &result); err != nil {
		errmsg := fmt.Sprintf("无法上传图片, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	}

	if result.ErrCode != core.ErrCodeOK {
		errmsg := fmt.Sprintf("无法上传图片, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, &result.MediaUrl)
	glog.Info("%s upload media picture success, url %s", result.MediaUrl.url)
	return
}
