package main

import (
	"bytes"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"strconv"
)

type PictureList struct {
	Count    int       `json:"count"`
	Pictures []Picture `json:"pictures"`
}

func (p Picture) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/picture").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(p.findPicture))
	ws.Route(ws.GET("/{picture_id}").To(p.findPicture))
	ws.Route(ws.POST("").To(p.createPicture))
	ws.Route(ws.PUT("/{picture_id}").To(p.updatePicture))
	ws.Route(ws.DELETE("/{picture_id}").To(p.deletePicture))
	container.Add(ws)
}

func (p Picture) findPicture(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [FIND_Picture]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	picture_id := request.PathParameter("picture_id")

	//return picture list
	if picture_id == "" {
		pictureList := PictureList{}
		pictureList.Pictures = make([]Picture, 0)
		db.Debug().Find(&pictureList.Pictures)
		pictureList.Count = len(pictureList.Pictures)
		glog.Infof("%s Return picture list", prefix)
		response.WriteHeaderAndEntity(http.StatusOK, pictureList)
		return
	}

	//get picture id
	id, err := strconv.Atoi(picture_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get picture, picture_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	picture := Picture{}
	db.Debug().First(&picture, id)

	//cannot find picture
	if picture.ID == 0 {
		errmsg := fmt.Sprintf("cannot find picture with id %s", picture_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	//find picture
	glog.Infof("%s picture with id %d found and return", prefix, picture.ID)
	response.WriteHeaderAndEntity(http.StatusOK, picture)
	return
}

func (p Picture) createPicture(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createPicture]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	picture := Picture{}
	err := request.ReadEntity(&picture)
	if err == nil {
		db.Debug().Create(&picture)
		if picture.ID != 0 {
			//create picture successfully
			glog.Infof("%s create picture with id %d successfully", prefix, picture.ID)
			response.WriteHeaderAndEntity(http.StatusCreated, picture)
			return
		} else {
			//fail to create picture
			errmsg := fmt.Sprintf("%s cannot create picture on database", prefix)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		}
	} else {
		//parse picture entity failed
		errmsg := fmt.Sprintf("cannot create picture, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (p Picture) updatePicture(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updatePicture]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	picture_id := request.PathParameter("picture_id")
	picture := Picture{}
	err := request.ReadEntity(&picture)
	//fail to parse the picture entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update picture, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(picture_id)
	//fail to parse picture id
	if err != nil {
		errmsg := fmt.Sprintf("cannot update picture, path picture_id is %s, err %s", picture_id, err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != picture.ID {
		errmsg := fmt.Sprintf("cannot update picture, path picture_id %d is not equal to id %d in body content", id, picture.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realPicture := Picture{}
	db.Debug().First(&realPicture, picture.ID)
	//cannot find picture
	if realPicture.ID == 0 {
		errmsg := fmt.Sprintf("cannot update picture, picture_id %d is not exist", picture.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	//find picture and update
	db.Debug().Model(&realPicture).Update(picture)
	glog.Infof("%s update picture with id %d on database", prefix, realPicture.ID)
	response.WriteHeaderAndEntity(http.StatusCreated, realPicture)
	return
}

func (p Picture) deletePicture(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deletePicture]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	picture_id := request.PathParameter("picture_id")
	id, err := strconv.Atoi(picture_id)
	//fail to parse picture id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete picture, picture_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	picture := Picture{}
	db.Debug().First(&picture, id)
	if picture.ID == 0 {
		//picture with id doesn't exist, return ok
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		glog.Infof("%s picture with id %s doesn't exist, return ok", prefix, id)
		return
	}

	db.Debug().Delete(&picture)

	realPicture := Picture{}
	db.Debug().First(&realPicture, id)

	if realPicture.ID != 0 {
		//fail to delete picture
		errmsg := fmt.Sprintf("cannot delete picture,some of other object is referencing")
		glog.Infof("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	//delete picture successfully
	glog.Infof("%s delete picture with id %d successfully", prefix, realPicture.ID)
	response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	return
}
