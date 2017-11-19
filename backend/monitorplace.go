package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"strconv"
	"time"
)

type MonitorPlaceList struct {
	Count         int            `json:"count"`
	MonitorPlaces []MonitorPlace `json:"monitor_places"`
}

type PictureWithMonitorPlace struct {
	MonitorPlaceId int       `json:"monitor_place_id"`
	Count          int       `json:"count"`
	Pictures       []Picture `json:"picture"`
}

func (m MonitorPlace) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/monitor_place").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(m.findMonitorPlace))
	ws.Route(ws.GET("/{monitor_place_id}").To(m.findMonitorPlace))
	ws.Route(ws.GET("/{monitor_place_id}/{scope}?after={after}&limit={limit}").To(m.findMonitorPlace))
	ws.Route(ws.POST("/{monitor_place_id}").To(m.updateMonitorPlace))
	ws.Route(ws.PUT("").To(m.createMonitorPlace))
	ws.Route(ws.DELETE("/{monitor_place_id}").To(m.deleteMonitorPlace))
	container.Add(ws)
}

func (m MonitorPlace) findMonitorPlace(request *restful.Request, response *restful.Response) {
	glog.Infof("GET %s", request.Request.URL)
	monitor_place_id := request.PathParameter("monitor_place_id")
	scope := request.PathParameter("scope")
	after := request.QueryParameter("after")
	limit := request.QueryParameter("limit")

	if monitor_place_id == "" {
		monitor_placeList := MonitorPlaceList{}
		monitor_placeList.MonitorPlaces = make([]MonitorPlace, 0)
		db.Find(&monitor_placeList.MonitorPlaces)
		monitor_placeList.Count = len(monitor_placeList.MonitorPlaces)
		response.WriteHeaderAndEntity(http.StatusOK, monitor_placeList)
		return
	}

	id, err := strconv.Atoi(monitor_place_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get monitor_place, monitor_place_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	monitor_place := MonitorPlace{}
	db.First(&monitor_place, id)
	if monitor_place.ID == 0 {
		errmsg := fmt.Sprintf("cannot find monitor_place with id %s", monitor_place_id)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	if scope == "" {
		response.WriteHeaderAndEntity(http.StatusOK, monitor_place)
		return
	}

	if scope == "picture" {
		pictureList := PictureWithMonitorPlace{}
		pictureList.MonitorPlaceId = monitor_place.ID
		pictureList.Pictures = make([]Picture, 0)
		if after != "" {
			loc, _ := time.LoadLocation("Asia/Shanghai")
			const shortFormat = "20160102"
			after_trans, err := time.ParseInLocation(shortFormat, after, loc)
			if err != nil {
				errmsg := fmt.Sprintf("cannot find object with after %s, err", after, err)
				response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
				return
			}
			after_str := fmt.Sprintf("%d-%d-%d", after_trans.Year(), after_trans.Month(), after_trans.Day())
			if limit == "" {
				db.Where("create_at >= ?", after_str).Model(&monitor_place).Related(&pictureList.Pictures)
				pictureList.Count = len(pictureList.Pictures)
				response.WriteHeaderAndEntity(http.StatusOK, pictureList)
				return
			} else {
				limit_trans, err := strconv.Atoi(limit)
				if err != nil {
					errmsg := fmt.Sprintf("cannot find object with limit %s, err", limit, err)
					response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
					return
				}
				db.Where("create_at >= ?", after_str).Model(&monitor_place).Related(&pictureList.Pictures).Limit(limit_trans)
				pictureList.Count = len(pictureList.Pictures)
				response.WriteHeaderAndEntity(http.StatusOK, pictureList)
				return
			}
		} else {
			if limit == "" {
				db.Model(&monitor_place).Related(&pictureList.Pictures)
				pictureList.Count = len(pictureList.Pictures)
				response.WriteHeaderAndEntity(http.StatusOK, pictureList)
				return
			} else {
				limit_trans, err := strconv.Atoi(limit)
				if err != nil {
					errmsg := fmt.Sprintf("cannot find object with limit %s, err", limit, err)
					response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
					return
				}
				db.Model(&monitor_place).Related(&pictureList.Pictures).Limit(limit_trans)
				pictureList.Count = len(pictureList.Pictures)
				response.WriteHeaderAndEntity(http.StatusOK, pictureList)
				return
			}
		}
		db.Model(&monitor_place).Related(&pictureList.Pictures)
		pictureList.Count = len(pictureList.Pictures)
		response.WriteHeaderAndEntity(http.StatusOK, pictureList)
		return
	}

	errmsg := fmt.Sprintf("cannot find object with scope %s", scope)
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}

func (m MonitorPlace) createMonitorPlace(request *restful.Request, response *restful.Response) {
	glog.Infof("PUT %s", request.Request.URL)
	monitor_place := MonitorPlace{}
	err := request.ReadEntity(&monitor_place)
	if err == nil {
		db.Create(&monitor_place)
	} else {
		errmsg := fmt.Sprintf("cannot create monitor_place, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (m MonitorPlace) updateMonitorPlace(request *restful.Request, response *restful.Response) {
	glog.Infof("POST %s", request.Request.URL)
	monitor_place_id := request.PathParameter("monitor_place_id")
	monitor_place := MonitorPlace{}
	err := request.ReadEntity(&monitor_place)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update monitor_place, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(monitor_place_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update monitor_place, path monitor_place_id is %s, err %s", monitor_place_id, err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != monitor_place.ID {
		errmsg := fmt.Sprintf("cannot update monitor_place, path monitor_place_id %d is not equal to id %d in body content", id, monitor_place.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realMonitorPlace := MonitorPlace{}
	db.First(&realMonitorPlace, monitor_place.ID)
	if realMonitorPlace.ID == 0 {
		errmsg := fmt.Sprintf("cannot update monitor_place, monitor_place_id %d is not exist", monitor_place.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	db.Model(&realMonitorPlace).Update(monitor_place)
	response.WriteHeaderAndEntity(http.StatusCreated, &realMonitorPlace)
}

func (m MonitorPlace) deleteMonitorPlace(request *restful.Request, response *restful.Response) {
	glog.Infof("DELETE %s", request.Request.URL)
	monitor_place_id := request.PathParameter("monitor_place_id")
	id, err := strconv.Atoi(monitor_place_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete monitor_place, monitor_place_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	monitor_place := MonitorPlace{}
	db.First(&monitor_place, id)
	if monitor_place.ID == 0 {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&monitor_place)

	realMonitorPlace := MonitorPlace{}
	db.First(&realMonitorPlace, id)

	if realMonitorPlace.ID != 0 {
		errmsg := fmt.Sprintf("cannot delete monitor_place,some of other object is referencing")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	}
}
