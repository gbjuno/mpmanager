package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"strconv"
)

type PictureList struct {
	Count    int       `json:"count"`
	Pictures []Picture `json:"pictures"`
}

func (p Picture) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/picture").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(p.findPicture))
	ws.Route(ws.GET("/{picture_id}").To(p.findPicture))
	ws.Route(ws.POST("/{picture_id}").To(p.updatePicture))
	ws.Route(ws.PUT("").To(p.createPicture))
	ws.Route(ws.DELETE("/{picture_id}").To(p.deletePicture))
	container.Add(ws)
}

func (p Picture) findPicture(request *restful.Request, response *restful.Response) {
	glog.Infof("GET %s", request.Request.URL)
	picture_id := request.PathParameter("picture_id")

	if picture_id == "" {
		pictureList := PictureList{}
		pictureList.Pictures = make([]Picture, 0)
		db.Find(&pictureList.Pictures)
		pictureList.Count = len(pictureList.Pictures)
		response.WriteHeaderAndEntity(http.StatusOK, pictureList)
		return
	}

	id, err := strconv.Atoi(picture_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get picture, picture_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	picture := Picture{}
	db.First(&picture, id)
	if picture.ID == 0 {
		errmsg := fmt.Sprintf("cannot find picture with id %s", picture_id)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, picture)
		return
	}
}

func (p Picture) createPicture(request *restful.Request, response *restful.Response) {
	glog.Infof("PUT %s", request.Request.URL)
	picture := Picture{}
	err := request.ReadEntity(&picture)
	if err == nil {
		db.Create(&picture)
		response.WriteHeaderAndEntity(http.StatusCreated, picture)
		return
	} else {
		errmsg := fmt.Sprintf("cannot create picture, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (p Picture) updatePicture(request *restful.Request, response *restful.Response) {
	glog.Infof("POST %s", request.Request.URL)
	picture_id := request.PathParameter("picture_id")
	picture := Picture{}
	err := request.ReadEntity(&picture)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update picture, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(picture_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update picture, path picture_id is %s, err %s", picture_id, err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != picture.ID {
		errmsg := fmt.Sprintf("cannot update picture, path picture_id %d is not equal to id %d in body content", id, picture.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realPicture := Picture{}
	db.First(&realPicture, picture.ID)
	if realPicture.ID == 0 {
		errmsg := fmt.Sprintf("cannot update picture, picture_id %d is not exist", picture.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	db.Model(&realPicture).Update(picture)
	response.WriteHeaderAndEntity(http.StatusCreated, &realPicture)
	return
}

func (p Picture) deletePicture(request *restful.Request, response *restful.Response) {
	glog.Infof("DELETE %s", request.Request.URL)
	picture_id := request.PathParameter("picture_id")
	id, err := strconv.Atoi(picture_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete picture, picture_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	picture := Picture{}
	db.First(&picture, id)
	if picture.ID == 0 {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&picture)

	realPicture := Picture{}
	db.First(&realPicture, id)

	if realPicture.ID != 0 {
		errmsg := fmt.Sprintf("cannot delete picture,some of other object is referencing")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
