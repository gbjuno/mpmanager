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
	ws.Route(ws.GET("/{monitor_place_id}/").To(m.findMonitorPlace))
	ws.Route(ws.GET("/{monitor_place_id}/{scope}/?after={after}&limit={limit}&order={order}").To(m.findMonitorPlace))
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
	order := request.QueryParameter("order")

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
		glog.Infof("%s return monitor_place with id %s", prefix, monitor_place_id)
		response.WriteHeaderAndEntity(http.StatusOK, monitor_place)
		return
	}

	//find pictures related to monitor_place
	if scope == "picture" {
		pictureList := PictureWithMonitorPlace{}
		pictureList.MonitorPlaceId = monitor_place.ID
		pictureList.Pictures = make([]Picture, 0)

		var limitInt int
		var afterOk bool
		var limitOk bool
		var err error

		if after != "" {
			loc, _ := time.LoadLocation("Asia/Shanghai")
			const shortFormat = "20060102"
			_, err = time.ParseInLocation(shortFormat, after, loc)
			if err != nil {
				errmsg := fmt.Sprintf("cannot find object with after %s, err %s, ignore", after, err)
				glog.Errorf("%s %s", prefix, errmsg)
			}
			afterOk = true
		}

		if limit != "" {
			limitInt, err = strconv.Atoi(limit)
			if err != nil {
				errmsg := fmt.Sprintf("cannot find object with limit %s, err %s, ignore", limit, err)
				glog.Errorf("%s %s", prefix, errmsg)
			}
			limitOk = true
		}

		if order != "asc" && order != "desc" {
			errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		if order == "" {
			order = "desc"
		}

		if afterOk && limitOk {
			condition := fmt.Sprintf("create_at >= str_to_date(%s, '%%%%Y%%%%m%%%%d')", after)
			glog.Infof("%s find pictureList with condition %s", prefix, condition)
			db.Where(condition).Where("monitor_place_id = ?", monitor_place_id).Limit(limitInt).Order("id " + order).Find(&pictureList.Pictures)
			pictureList.Count = len(pictureList.Pictures)
			response.WriteHeaderAndEntity(http.StatusOK, pictureList)
			glog.Infof("%s return pictureList after %s with limit %s with monitor_place id %s and limit ", prefix, after, limit, monitor_place_id)
			return
		} else if afterOk {
			condition := fmt.Sprintf("create_at >= str_to_date(%s, '%%%%Y%%%%m%%%%d')", after)
			glog.Infof("%s find pictureList with condition %s", prefix, condition)
			db.Where(condition).Where("monitor_place_id = ?", monitor_place_id).Order("id" + order).Find(&pictureList.Pictures)
			pictureList.Count = len(pictureList.Pictures)
			response.WriteHeaderAndEntity(http.StatusOK, pictureList)
			glog.Infof("%s return pictureList after %s with monitor_place id %s", prefix, after, monitor_place_id)
			return
		} else if limitOk {
			db.Limit(limitInt).Where("monitor_place_id = ?", monitor_place.ID).Order("id " + order).Find(&pictureList.Pictures)
			pictureList.Count = len(pictureList.Pictures)
			response.WriteHeaderAndEntity(http.StatusOK, pictureList)
			glog.Infof("%s return pictureList with limit %s with monitor_place id %", prefix, limit, monitor_place_id)
			return
		} else {
			db.Where("monitor_place_id = ?", monitor_place.ID).Order("id " + order).Find(&pictureList.Pictures)
			pictureList.Count = len(pictureList.Pictures)
			response.WriteHeaderAndEntity(http.StatusOK, pictureList)
			glog.Infof("%s return pictureList with monitor_place id %s", prefix, monitor_place_id)
			return
		}
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
		company := Company{}
		db.First(&company, monitor_place.CompanyId)
		if company.ID == 0 {
			errmsg := fmt.Sprintf("company id %d not exists", company.ID)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		}

		//whether monitor_place name is unique in the same town
		monitor_places := make([]MonitorPlace, 0)
		db.Where("company_id = ?", company.ID).Find(&monitor_places)
		for _, m := range monitor_places {
			if m.Name == monitor_place.Name {
				errmsg := fmt.Sprintf("monitor_place %s already exists int the same company %s", monitor_place.Name, company.Name)
				glog.Errorf("%s %s", prefix, errmsg)
				response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
				return
			}
		}

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

			//insert a new row into TodaySummary
			timeNow := time.Now()
			todayStr := fmt.Sprintf("%d%d%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
			shortForm := "20160102"
			todayTime, _ := time.Parse(shortForm, todayStr)
			todaySummary := TodaySummary{Day: todayTime, CompanyId: company.ID, CompanyName: company.Name, MonitorPlaceId: monitor_place.ID, MonitorPlaceName: monitor_place.Name, IsUpload: "F", Corrective: "F", EverCorrective: "F"}
			glog.Info("%s try to create todaySummary for company with id %d succesfully", prefix, company.ID)
			db.Create(&todaySummary)

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
