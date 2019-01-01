package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

type GlobalRelaxPeriodList struct {
	Count              int                 `json:"count"`
	GlobalRelaxPeriods []GlobalRelaxPeriod `json:"global_relax_periods"`
}

func (c GlobalRelaxPeriod) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/global_relax_period").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(c.findGlobalRelaxPeriod))
	ws.Route(ws.GET("/?pageNo={pageNo}&pageSize={pageSize}&order={order}").To(c.findGlobalRelaxPeriod))
	ws.Route(ws.GET("/{global_relax_period_id}").To(c.findGlobalRelaxPeriod))
	ws.Route(ws.POST("").To(c.createGlobalRelaxPeriod))
	ws.Route(ws.PUT("/{global_relax_period_id}").To(c.updateGlobalRelaxPeriod))
	ws.Route(ws.DELETE("/{global_relax_period_id}").To(c.deleteGlobalRelaxPeriod))
	container.Add(ws)
}

func (c GlobalRelaxPeriod) findGlobalRelaxPeriod(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findGlobalRelaxPeriod]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	global_relax_period_id := request.PathParameter("global_relax_period_id")
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")

	var searchGlobalRelaxPeriod *gorm.DB = db.Debug()

	if order != "asc" && order != "desc" {
		errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
		glog.Errorf("%s %s", prefix, errmsg)
		order = "asc"
	}

	if order == "" {
		order = "asc"
	}

	glog.Infof("%s find global_relax_period with order %s", prefix, order)

	global_relax_period_relax_periods := make([]GlobalRelaxPeriod, 0)
	count := 0
	searchGlobalRelaxPeriod.Find(&global_relax_period_relax_periods).Count(&count)
	searchGlobalRelaxPeriod = searchGlobalRelaxPeriod.Order("id " + order)

	if global_relax_period_id == "" {
		isPageSizeOk := true
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			isPageSizeOk = false
			errmsg := fmt.Sprintf("cannot find object with pageSize %s, err %s, ignore", pageSize, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		//pageNo depends on pageSize
		isPageNoOk := true
		pageNoInt, err := strconv.Atoi(pageNo)
		if err != nil {
			isPageNoOk = false
			errmsg := fmt.Sprintf("cannot find object with pageNo %s, err %s, ignore", pageNo, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		if isPageSizeOk && isPageNoOk {
			limit := pageSizeInt
			offset := (pageNoInt - 1) * limit
			glog.Infof("%s set find global_relax_period db with pageSize %s, pageNo %s(limit %d, offset %d)", prefix, pageSize, pageNo, limit, offset)
			searchGlobalRelaxPeriod = searchGlobalRelaxPeriod.Offset(offset).Limit(limit)
		}

		global_relax_period_list := GlobalRelaxPeriodList{}
		global_relax_period_list.GlobalRelaxPeriods = make([]GlobalRelaxPeriod, 0)
		searchGlobalRelaxPeriod.Find(&global_relax_period_list.GlobalRelaxPeriods)

		global_relax_period_list.Count = count
		response.WriteHeaderAndEntity(http.StatusOK, global_relax_period_list)
		glog.Infof("%s return global_relax_period list", prefix)
		return
	}

	id, err := strconv.Atoi(global_relax_period_id)
	//fail to parse global_relax_period id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get global_relax_period, global_relax_period_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	global_relax_period := GlobalRelaxPeriod{}
	db.Debug().First(&global_relax_period, id)
	//cannot find global_relax_period
	if global_relax_period.ID == 0 {
		errmsg := fmt.Sprintf("cannot find global_relax_period with id %s", global_relax_period_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	glog.Infof("%s return global_relax_period with id %d", prefix, global_relax_period.ID)
	response.WriteHeaderAndEntity(http.StatusOK, global_relax_period)
	return
}

func (c GlobalRelaxPeriod) createGlobalRelaxPeriod(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createGlobalRelaxPeriod]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	global_relax_period := GlobalRelaxPeriod{}
	err := request.ReadEntity(&global_relax_period)
	if err != nil {
		//fail to parse global_relax_period entity
		errmsg := fmt.Sprintf("cannot create global_relax_period, err %s", err)
		returnmsg := fmt.Sprintf("无法创建信息，提供的信息错误")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	db.Debug().Create(&global_relax_period)
	if global_relax_period.ID == 0 {
		//fail to create global_relax_period on database
		errmsg := fmt.Sprintf("cannot create global_relax_period on database")
		returnmsg := fmt.Sprintf("无法创建全局假期，请联系管理员")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//create global_relax_period on database
	glog.Infof("%s create global_relax_period with id %d succesfully", prefix, global_relax_period.ID)
	response.WriteHeaderAndEntity(http.StatusOK, global_relax_period)
	//insert a new row into Summary

	startAt := getDateStr(global_relax_period.StartAt)
	endAt := getDateStr(global_relax_period.EndAt)

	summaries := make([]Summary, 0)
	db.Debug().Find(&summaries, "(day between '?' and '?')", startAt, endAt)
	for _, s := range summaries {
		db.Model(&s).Update("relax_day", 'T')
	}
	return
}

func (c GlobalRelaxPeriod) updateGlobalRelaxPeriod(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateGlobalRelaxPeriod]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	global_relax_period_id := request.PathParameter("global_relax_period_id")
	global_relax_period := GlobalRelaxPeriod{}
	err := request.ReadEntity(&global_relax_period)

	//fail to parse global_relax_period entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update global_relax_period, err %s", err)
		returnmsg := fmt.Sprintf("无法更新全局休假信息，提供的信息解析失败")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//fail to parse global_relax_period id
	id, err := strconv.Atoi(global_relax_period_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update global_relax_period, path global_relax_period_id is %s, err %s", global_relax_period_id, err)
		returnmsg := fmt.Sprintf("无法更新全局休假信息，提供的休假id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	if id != global_relax_period.ID {
		errmsg := fmt.Sprintf("cannot update global_relax_period, path global_relax_period_id %d is not equal to id %d in body content", id, global_relax_period.ID)
		returnmsg := fmt.Sprintf("无法更新全局休假信息，提供的休假id与URL中的休假id不匹配")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	realGlobalRelaxPeriod := GlobalRelaxPeriod{}
	db.Debug().First(&realGlobalRelaxPeriod, global_relax_period.ID)

	//cannot find global_relax_period
	if realGlobalRelaxPeriod.ID == 0 {
		errmsg := fmt.Sprintf("cannot update global_relax_period, global_relax_period_id %d does not exist", global_relax_period.ID)
		returnmsg := fmt.Sprintf("该项公司休假信息已被删除")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	summaries := make([]Summary, 0)
	startAt := getDateStr(realGlobalRelaxPeriod.StartAt)
	endAt := getDateStr(realGlobalRelaxPeriod.EndAt)
	db.Debug().Find(&summaries, "(day between '?' and '?')", startAt, endAt)
	for _, s := range summaries {
		db.Model(&s).Update("relax_day", 'F')
	}

	//find global_relax_period and update
	db.Debug().Model(&realGlobalRelaxPeriod).Update(global_relax_period)
	glog.Infof("%s update global_relax_period with id %d successfully and return", prefix, realGlobalRelaxPeriod.ID)

	summaries = summaries[:0]
	startAt = getDateStr(global_relax_period.StartAt)
	endAt = getDateStr(global_relax_period.EndAt)
	db.Debug().Find(&summaries, "(day between '?' and '?')", startAt, endAt)
	for _, s := range summaries {
		db.Model(&s).Update("relax_day", 'T')
	}

	response.WriteHeaderAndEntity(http.StatusOK, realGlobalRelaxPeriod)
	return
}

func (c GlobalRelaxPeriod) deleteGlobalRelaxPeriod(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteGlobalRelaxPeriod]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	global_relax_period_id := request.PathParameter("global_relax_period_id")
	id, err := strconv.Atoi(global_relax_period_id)
	//fail to parse global_relax_period id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete global_relax_period, global_relax_period_id is not integer, err %s", err)
		returnmsg := fmt.Sprintf("无法删除公司休假信息，提供的公司id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	global_relax_period := GlobalRelaxPeriod{}
	db.Debug().First(&global_relax_period, id)
	if global_relax_period.ID == 0 {
		//global_relax_period with id doesn't exist, return ok
		glog.Infof("%s global_relax_period with id %s doesn't exist, return ok", prefix, global_relax_period_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	summaries := make([]Summary, 0)
	startAt := getDateStr(global_relax_period.StartAt)
	endAt := getDateStr(global_relax_period.EndAt)
	db.Debug().Find(&summaries, "(day between '?' and '?')", startAt, endAt)
	for _, s := range summaries {
		db.Model(&s).Update("relax_day", 'F')
	}
	db.Debug().Delete(&global_relax_period)

	//delete global_relax_period successfully
	glog.Infof("%s delete global_relax_period with id %s successfully", prefix, global_relax_period_id)
	response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	return
}
