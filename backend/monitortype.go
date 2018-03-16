package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
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
	ws.Path(RESTAPIVERSION + "/monitor_type").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(m.findMonitorType))
	ws.Route(ws.GET("/{monitor_type_id}").To(m.findMonitorType))
	ws.Route(ws.GET("/{monitor_type_id}/{scope}").To(m.findMonitorType))
	ws.Route(ws.POST("").To(m.createMonitorType))
	ws.Route(ws.PUT("/{monitor_type_id}").To(m.updateMonitorType))
	ws.Route(ws.DELETE("/{monitor_type_id}").To(m.deleteMonitorType))
	container.Add(ws)
}

func (m MonitorType) findMonitorType(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findMonitorType]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	monitor_type_id := request.PathParameter("monitor_type_id")
	scope := request.PathParameter("scope")

	//get monitory_type list
	if monitor_type_id == "" {
		monitor_typeList := MonitorTypeList{}
		monitor_typeList.MonitorTypes = make([]MonitorType, 0)
		db.Debug().Find(&monitor_typeList.MonitorTypes)
		monitor_typeList.Count = len(monitor_typeList.MonitorTypes)
		glog.Infof("%s return monitor_type list", prefix)
		response.WriteHeaderAndEntity(http.StatusOK, monitor_typeList)
		return
	}

	id, err := strconv.Atoi(monitor_type_id)
	//fail to parse monitor_type id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get monitor_type, monitor_type_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	monitor_type := MonitorType{}
	db.Debug().First(&monitor_type, id)
	//cannot find monitor_type
	if monitor_type.ID == 0 {
		errmsg := fmt.Sprintf("cannot find monitor_type with id %s", monitor_type_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	//find monitor_type
	if scope == "" {
		glog.Infof("%s return monitor_place with id %d", prefix, monitor_type.ID)
		response.WriteHeaderAndEntity(http.StatusOK, monitor_type)
		return
	}

	//find monitor_place related to monitor_type with id
	if scope == "monitorplace" {
		monitor_placeList := MonitorPlaceWithMonitorType{}
		monitor_placeList.MonitorTypeId = monitor_type.ID
		monitor_placeList.MonitorPlaces = make([]MonitorPlace, 0)
		db.Debug().Model(&monitor_type).Related(&monitor_placeList)
		monitor_placeList.Count = len(monitor_placeList.MonitorPlaces)
		glog.Infof("%s return monitor_place related to monitor_type with id %d", prefix, monitor_type.ID)
		response.WriteHeaderAndEntity(http.StatusOK, monitor_placeList)
		return
	}

	errmsg := fmt.Sprintf("cannot find object with scope %s", scope)
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}

func (m MonitorType) createMonitorType(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createMonitorType]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	monitor_type := MonitorType{}
	err := request.ReadEntity(&monitor_type)
	if err == nil {
		sameNameMonitorType := MonitorType{}
		db.Debug().Where("name = ?", monitor_type.Name).First(&sameNameMonitorType)
		if sameNameMonitorType.ID != 0 {
			errmsg := fmt.Sprintf("monitor_type %s already exists", sameNameMonitorType.Name)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		}

		db.Debug().Create(&monitor_type)
		if monitor_type.ID == 0 {
			//fail to create monitor_type on database
			errmsg := fmt.Sprintf("cannot create monitor_type on database")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		} else {
			//create monitor_type on database
			glog.Infof("%s create monitor_type with id %d succesfully", prefix, monitor_type.ID)
			response.WriteHeaderAndEntity(http.StatusOK, monitor_type)
			return
		}
	} else {
		//fail to parse monitor_type entity
		errmsg := fmt.Sprintf("cannot create monitor_type, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (m MonitorType) updateMonitorType(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateMonitorType]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	monitor_type_id := request.PathParameter("monitor_type_id")
	monitor_type := MonitorType{}
	err := request.ReadEntity(&monitor_type)

	//fail to parse monitor_type entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update monitor_type, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(monitor_type_id)

	//fail to parse monitor_type id
	if err != nil {
		errmsg := fmt.Sprintf("cannot update monitor_type, path monitor_type_id is %s, err %s", monitor_type_id, err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != monitor_type.ID {
		errmsg := fmt.Sprintf("cannot update monitor_type, path monitor_type_id %d is not equal to id %d in body content", id, monitor_type.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realMonitorType := MonitorType{}
	db.Debug().First(&realMonitorType, monitor_type.ID)

	//cannot find monitor_type
	if realMonitorType.ID == 0 {
		errmsg := fmt.Sprintf("cannot update monitor_type, monitor_type_id %d is not exist", monitor_type.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
	//find monitor_type
	db.Debug().Model(&realMonitorType).Update(monitor_type)
	glog.Infof("%s update monitor_type with id %d successfully and return", prefix, realMonitorType.ID)
	response.WriteHeaderAndEntity(http.StatusCreated, realMonitorType)
	return
}

func (m MonitorType) deleteMonitorType(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteMonitorType]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	monitor_type_id := request.PathParameter("monitor_type_id")
	id, err := strconv.Atoi(monitor_type_id)
	//fail to parse monitor
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete monitor_type, monitor_type_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	monitor_type := MonitorType{}
	db.Debug().First(&monitor_type, id)
	if monitor_type.ID == 0 {
		//monitor_place with id doesn't exist, return ok
		glog.Infof("%s company with id %s doesn't exist, return ok", prefix, monitor_type.ID)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Debug().Delete(&monitor_type)

	realMonitorType := MonitorType{}
	db.Debug().First(&realMonitorType, id)

	if realMonitorType.ID != 0 {
		//fail to delete monitor_place
		errmsg := fmt.Sprintf("cannot delete monitor_type,some of other object is referencing")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	} else {
		//delete monitor_place successfully
		glog.Infof("%s delete monitor_place with id %s successfully", prefix, monitor_type_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
