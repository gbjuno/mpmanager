package main

import (
	"bytes"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/gbjuno/mpmanager/backend/utils"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"os"
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
	ws.Path(RESTAPIVERSION + "/monitor_place").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(m.findMonitorPlace))
	ws.Route(ws.GET("/{monitor_place_id}").To(m.findMonitorPlace))
	ws.Route(ws.GET("/{monitor_place_id}/{scope}?after={after}&limit={limit}").To(m.findMonitorPlace))
	ws.Route(ws.POST("").To(m.createMonitorPlace))
	ws.Route(ws.PUT("/{monitor_place_id}").To(m.updateMonitorPlace))
	ws.Route(ws.DELETE("/{monitor_place_id}").To(m.deleteMonitorPlace))
	container.Add(ws)
}

func (m MonitorPlace) findMonitorPlace(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findMonitorPlace]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	monitor_place_id := request.PathParameter("monitor_place_id")
	scope := request.PathParameter("scope")
	after := request.QueryParameter("after")
	limit := request.QueryParameter("limit")

	//get monitor_place list
	if monitor_place_id == "" {
		monitor_placeList := MonitorPlaceList{}
		monitor_placeList.MonitorPlaces = make([]MonitorPlace, 0)
		db.Find(&monitor_placeList.MonitorPlaces)
		monitor_placeList.Count = len(monitor_placeList.MonitorPlaces)
		response.WriteHeaderAndEntity(http.StatusOK, monitor_placeList)
		glog.Infof("%s return monitor_place list", prefix)
		return
	}

	id, err := strconv.Atoi(monitor_place_id)
	//fail to parse monitor_place id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get monitor_place, monitor_place_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	monitor_place := MonitorPlace{}
	db.First(&monitor_place, id)
	//cannot find monitor_place
	if monitor_place.ID == 0 {
		errmsg := fmt.Sprintf("cannot find monitor_place with id %s", monitor_place_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	//find monitor_place, set QrcodePath
	if scope == "" {
		glog.Infof("%s return monitor_place with id %d", prefix, monitor_place_id)
		response.WriteHeaderAndEntity(http.StatusOK, monitor_place)
		return
	}

	//find pictures related to monitor_place
	if scope == "picture" {
		pictureList := PictureWithMonitorPlace{}
		pictureList.MonitorPlaceId = monitor_place.ID
		pictureList.Pictures = make([]Picture, 0)
		var after_str string
		var limit_trans int
		if after != "" {
			loc, _ := time.LoadLocation("Asia/Shanghai")
			const shortFormat = "20160102"
			after_trans, err := time.ParseInLocation(shortFormat, after, loc)
			if err != nil {
				errmsg := fmt.Sprintf("cannot find object with after %s, err", after, err)
				glog.Errorf("%s %s", prefix, errmsg)
				response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
				return
			}
			after_str = fmt.Sprintf("%d-%d-%d", after_trans.Year(), after_trans.Month(), after_trans.Day())
			if limit == "" {
				db.Where("create_at >= ?", after_str).Model(&monitor_place).Related(&pictureList.Pictures)
				pictureList.Count = len(pictureList.Pictures)
				response.WriteHeaderAndEntity(http.StatusOK, pictureList)
				glog.Infof("%s return pictureList after %s with monitor_place id %d", prefix, after_str, monitor_place_id)
				return
			} else {
				limit_trans, err = strconv.Atoi(limit)
				if err != nil {
					errmsg := fmt.Sprintf("cannot find object with limit %s, err", limit, err)
					glog.Errorf("%s %s", prefix, errmsg)
					response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
					return
				}
				db.Where("create_at >= ?", after_str).Model(&monitor_place).Related(&pictureList.Pictures).Limit(limit_trans)
				pictureList.Count = len(pictureList.Pictures)
				response.WriteHeaderAndEntity(http.StatusOK, pictureList)
				glog.Infof("%s return pictureList after %s with limit %s with monitor_place id %d and limit ", prefix, after_str, limit, monitor_place_id)
				return
			}
		} else {
			if limit == "" {
				db.Model(&monitor_place).Related(&pictureList.Pictures)
				pictureList.Count = len(pictureList.Pictures)
				response.WriteHeaderAndEntity(http.StatusOK, pictureList)
				glog.Infof("%s return pictureList after %s with limit %s with monitor_place id %d and limit ", prefix, after_str, limit, monitor_place_id)
				return
			} else {
				limit_trans, err := strconv.Atoi(limit)
				if err != nil {
					errmsg := fmt.Sprintf("cannot find object with limit %s, err", limit, err)
					glog.Errorf("%s %s", prefix, errmsg)
					response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
					return
				}
				db.Model(&monitor_place).Related(&pictureList.Pictures).Limit(limit_trans)
				pictureList.Count = len(pictureList.Pictures)
				response.WriteHeaderAndEntity(http.StatusOK, pictureList)
				glog.Infof("%s return pictureList with limit %s with monitor_place id %d and limit ", prefix, limit, monitor_place_id)
				return
			}
		}
		db.Model(&monitor_place).Related(&pictureList.Pictures)
		pictureList.Count = len(pictureList.Pictures)
		response.WriteHeaderAndEntity(http.StatusOK, pictureList)
		glog.Infof("%s return pictureList with limit %s with monitor_place id %d and limit ", prefix, limit, monitor_place_id)
		return
	}

	errmsg := fmt.Sprintf("cannot find object with scope %s", scope)
	glog.Errorf("%s %s", prefix, errmsg)
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}

func (m MonitorPlace) createMonitorPlace(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createMonitorPlace]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	monitor_place := MonitorPlace{}
	err := request.ReadEntity(&monitor_place)
	if err == nil {
		db.Create(&monitor_place)
		if monitor_place.ID == 0 {
			//fail to create monitor_place on database
			errmsg := fmt.Sprintf("cannot create monitor_place, err %s", err)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		} else {
			//create monitor_place on databse
			monitor_place.QrcodePath = fmt.Sprintf("/qrcode/%d/%d.png", monitor_place.CompanyId, monitor_place.ID)
			monitor_place.QrcodeURI = fmt.Sprintf("/static/qrcode/%d/%d.png", monitor_place.CompanyId, monitor_place.ID)
			db.Save(&monitor_place)
			company := Company{}
			db.First(&company, monitor_place.CompanyId)
			companyName := company.Name
			qrcodePath := fmt.Sprintf("https://%s/backend/photo?place=%d", domain, monitor_place.ID)
			//create monitor_place qrcode image
			if err := utils.GenerateQrcodeImage(qrcodePath, companyName+monitor_place.Name, imgRepo+monitor_place.QrcodePath); err != nil {
				errmsg := fmt.Sprintf("cannot create qrcode for monitor_place %d, err %s", monitor_place.ID, err)
				glog.Errorf("%s %s", prefix, errmsg)
			}
			glog.Infof("%s create monitor_place, id %d", prefix, monitor_place.ID)
			response.WriteHeaderAndEntity(http.StatusOK, monitor_place)
			return
		}
	} else {
		//failed to parse monitor_place entity
		errmsg := fmt.Sprintf("cannot create monitor_place, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
	return
}

func (m MonitorPlace) updateMonitorPlace(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateMonitorPlace]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	monitor_place_id := request.PathParameter("monitor_place_id")
	monitor_place := MonitorPlace{}
	err := request.ReadEntity(&monitor_place)

	//fail to parse monitor_place entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update monitor_place, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	//fail to parse monitor_place id
	id, err := strconv.Atoi(monitor_place_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update monitor_place, path monitor_place_id is %s, err %s", monitor_place_id, err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != monitor_place.ID {
		errmsg := fmt.Sprintf("cannot update monitor_place, path monitor_place_id %d is not equal to id %d in body content", id, monitor_place.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realMonitorPlace := MonitorPlace{}
	db.First(&realMonitorPlace, monitor_place.ID)

	//cannot find monitor_place
	if realMonitorPlace.ID == 0 {
		errmsg := fmt.Sprintf("cannot update monitor_place, monitor_place_id %d is not exist", monitor_place.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	//find monitor_place
	db.Model(&realMonitorPlace).Update(monitor_place)
	glog.Infof("%s update monitor place with id %d succeed", prefix, realMonitorPlace.ID)
	response.WriteHeaderAndEntity(http.StatusOK, realMonitorPlace)
	return
}

func (m MonitorPlace) deleteMonitorPlace(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteMonitorPlace]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s DELETE %s, content %s", prefix, request.Request.URL, content)
	monitor_place_id := request.PathParameter("monitor_place_id")
	id, err := strconv.Atoi(monitor_place_id)
	//fail to parse monitor_place id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete monitor_place, monitor_place_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	monitor_place := MonitorPlace{}
	db.First(&monitor_place, id)
	if monitor_place.ID == 0 {
		//monitor_place with id doesn't exist, return ok
		glog.Infof("%s monitor_place with id %d doesn't exist, delete successfully", prefix, monitor_place_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&monitor_place)

	realMonitorPlace := MonitorPlace{}
	db.First(&realMonitorPlace, id)

	if realMonitorPlace.ID != 0 {
		//fail to delete monitor_place
		errmsg := fmt.Sprintf("cannot delete monitor_place,some of other object is referencing")
		glog.Infof("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	} else {
		//delete monitor_place successfully
		os.Remove(imgRepo + monitor_place.QrcodePath)
		glog.Infof("%s delete monitor_place with id %d, qrcode path %s successfully", prefix, monitor_place_id, monitor_place.QrcodePath)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
