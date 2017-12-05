package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"strconv"
)

type MonitorTypeList struct {
	Count        int           `json:"count"`
	MonitorTypes []MonitorType `json:"monitor_types"`
}

type MonitorPlaceWithMonitorType struct {
	MonitorTypeId int            `json:"monitor_type_id"`
	Count         int            `json:"count"`
	MonitorPlaces []MonitorPlace `json:"monitor_places"`
}

func (m MonitorType) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/monitor_type").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(m.findMonitorType))
	ws.Route(ws.GET("/{monitor_type_id}").To(m.findMonitorType))
	ws.Route(ws.GET("/{monitor_type_id}/{scope}").To(m.findMonitorType))
	ws.Route(ws.POST("").To(m.createMonitorType))
	ws.Route(ws.PUT("/{monitor_type_id}").To(m.updateMonitorType))
	ws.Route(ws.DELETE("/{monitor_type_id}").To(m.deleteMonitorType))
	container.Add(ws)
}

func (m MonitorType) findMonitorType(request *restful.Request, response *restful.Response) {
	glog.Infof("GET %s", request.Request.URL)
	monitor_type_id := request.PathParameter("monitor_type_id")
	scope := request.PathParameter("scope")

	if monitor_type_id == "" {
		monitor_typeList := MonitorTypeList{}
		monitor_typeList.MonitorTypes = make([]MonitorType, 0)
		db.Find(&monitor_typeList.MonitorTypes)
		monitor_typeList.Count = len(monitor_typeList.MonitorTypes)
		response.WriteHeaderAndEntity(http.StatusOK, monitor_typeList)
		return
	}

	id, err := strconv.Atoi(monitor_type_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get monitor_type, monitor_type_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	monitor_type := MonitorType{}
	db.First(&monitor_type, id)
	if monitor_type.ID == 0 {
		errmsg := fmt.Sprintf("cannot find monitor_type with id %s", monitor_type_id)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	if scope == "" {
		response.WriteHeaderAndEntity(http.StatusOK, monitor_type)
		return
	}

	if scope == "monitorplace" {
		monitor_placeList := MonitorPlaceWithMonitorType{}
		monitor_placeList.MonitorTypeId = monitor_type.ID
		monitor_placeList.MonitorPlaces = make([]MonitorPlace, 0)
		db.Model(&monitor_type).Related(&monitor_placeList)
		monitor_placeList.Count = len(monitor_placeList.MonitorPlaces)
		response.WriteHeaderAndEntity(http.StatusOK, monitor_placeList)
		return
	}

	errmsg := fmt.Sprintf("cannot find object with scope %s", scope)
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}

func (m MonitorType) createMonitorType(request *restful.Request, response *restful.Response) {
	glog.Infof("POST %s", request.Request.URL)
	monitor_type := MonitorType{}
	err := request.ReadEntity(&monitor_type)
	if err == nil {
		db.Create(&monitor_type)
	} else {
		errmsg := fmt.Sprintf("cannot create monitor_type, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (m MonitorType) updateMonitorType(request *restful.Request, response *restful.Response) {
	glog.Infof("PUT %s", request.Request.URL)
	monitor_type_id := request.PathParameter("monitor_type_id")
	monitor_type := MonitorType{}
	err := request.ReadEntity(&monitor_type)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update monitor_type, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(monitor_type_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update monitor_type, path monitor_type_id is %s, err %s", monitor_type_id, err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != monitor_type.ID {
		errmsg := fmt.Sprintf("cannot update monitor_type, path monitor_type_id %d is not equal to id %d in body content", id, monitor_type.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realMonitorType := MonitorType{}
	db.First(&realMonitorType, monitor_type.ID)
	if realMonitorType.ID == 0 {
		errmsg := fmt.Sprintf("cannot update monitor_type, monitor_type_id %d is not exist", monitor_type.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	db.Model(&realMonitorType).Update(monitor_type)
	response.WriteHeaderAndEntity(http.StatusCreated, &realMonitorType)
}

func (m MonitorType) deleteMonitorType(request *restful.Request, response *restful.Response) {
	glog.Infof("DELETE %s", request.Request.URL)
	monitor_type_id := request.PathParameter("monitor_type_id")
	id, err := strconv.Atoi(monitor_type_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete monitor_type, monitor_type_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	monitor_type := MonitorType{}
	db.First(&monitor_type, id)
	if monitor_type.ID == 0 {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&monitor_type)

	realMonitorType := MonitorType{}
	db.First(&realMonitorType, id)

	if realMonitorType.ID != 0 {
		errmsg := fmt.Sprintf("cannot delete monitor_type,some of other object is referencing")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	}
}
