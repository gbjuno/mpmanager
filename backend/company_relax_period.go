package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

type CompanyRelaxPeriodList struct {
	Count               int                  `json:"count"`
	CompanyRelaxPeriods []CompanyRelaxPeriod `json:"company_relax_periods"`
}

func (c CompanyRelaxPeriod) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/company_relax_period").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(c.findCompanyRelaxPeriod))
	ws.Route(ws.GET("?company_id={company_id}&pageNo={pageNo}&pageSize={pageSize}&order={order}").To(c.findCompanyRelaxPeriod))
	ws.Route(ws.GET("/?company_id={company_id}&pageNo={pageNo}&pageSize={pageSize}&order={order}").To(c.findCompanyRelaxPeriod))
	ws.Route(ws.GET("/{company_relax_period_id}").To(c.findCompanyRelaxPeriod))
	ws.Route(ws.POST("").To(c.createCompanyRelaxPeriod))
	ws.Route(ws.PUT("/{company_relax_period_id}").To(c.updateCompanyRelaxPeriod))
	ws.Route(ws.DELETE("/{company_relax_period_id}").To(c.deleteCompanyRelaxPeriod))
	container.Add(ws)
}

func (c CompanyRelaxPeriod) findCompanyRelaxPeriod(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findCompanyRelaxPeriod]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	company_relax_period_id := request.PathParameter("company_relax_period_id")
	companyID := request.QueryParameter("company_id")
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")

	var searchCompanyRelaxPeriod *gorm.DB = db.Debug()

	if order != "asc" && order != "desc" {
		errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
		glog.Errorf("%s %s", prefix, errmsg)
		order = "asc"
	}

	if order == "" {
		order = "asc"
	}

	glog.Infof("%s find company_relax_period with order %s", prefix, order)

	company_relax_period_relax_periods := make([]CompanyRelaxPeriod, 0)
	count := 0
	glog.Infof("%s find company_relax_period with companyID %s", prefix, companyID)
	if companyID != "" {
		searchCompanyRelaxPeriod = searchCompanyRelaxPeriod.Where("company_id = ?", companyID)
	}
	searchCompanyRelaxPeriod.Find(&company_relax_period_relax_periods).Count(&count)
	searchCompanyRelaxPeriod = searchCompanyRelaxPeriod.Order("id " + order)

	if company_relax_period_id == "" {
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
			glog.Infof("%s set find company_relax_period db with pageSize %s, pageNo %s(limit %d, offset %d)", prefix, pageSize, pageNo, limit, offset)
			searchCompanyRelaxPeriod = searchCompanyRelaxPeriod.Offset(offset).Limit(limit)
		}

		company_relax_period_list := CompanyRelaxPeriodList{}
		company_relax_period_list.CompanyRelaxPeriods = make([]CompanyRelaxPeriod, 0)
		searchCompanyRelaxPeriod.Find(&company_relax_period_list.CompanyRelaxPeriods)
		company_relax_period_list.Count = count

		for i, _ := range company_relax_period_list.CompanyRelaxPeriods {
			company := Company{}
			db.Debug().First(&company, company_relax_period_list.CompanyRelaxPeriods[i].CompanyId)
			company_relax_period_list.CompanyRelaxPeriods[i].CompanyName = company.Name
		}

		response.WriteHeaderAndEntity(http.StatusOK, company_relax_period_list)
		glog.Infof("%s return company_relax_period list", prefix)
		return
	}

	id, err := strconv.Atoi(company_relax_period_id)
	//fail to parse company_relax_period id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get company_relax_period, company_relax_period_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	company_relax_period := CompanyRelaxPeriod{}
	db.Debug().First(&company_relax_period, id)
	//cannot find company_relax_period
	if company_relax_period.ID == 0 {
		errmsg := fmt.Sprintf("cannot find company_relax_period with id %s", company_relax_period_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	company := Company{}
	db.Debug().First(&company, company_relax_period.CompanyId)
	company_relax_period.CompanyName = company.Name

	glog.Infof("%s return company_relax_period with id %d", prefix, company_relax_period.ID)
	response.WriteHeaderAndEntity(http.StatusOK, company_relax_period)
	return
}

func (c CompanyRelaxPeriod) createCompanyRelaxPeriod(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createCompanyRelaxPeriod]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	company_relax_period := CompanyRelaxPeriod{}
	err := request.ReadEntity(&company_relax_period)
	if err != nil {
		//fail to parse company_relax_period entity
		errmsg := fmt.Sprintf("cannot create company_relax_period, err %s", err)
		returnmsg := fmt.Sprintf("无法创建信息，提供的信息错误")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	db.Debug().Create(&company_relax_period)
	if company_relax_period.ID == 0 {
		//fail to create company_relax_period on database
		errmsg := fmt.Sprintf("cannot create company_relax_period on database")
		returnmsg := fmt.Sprintf("无法创建公司假期，请联系管理员")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//create company_relax_period on database
	glog.Infof("%s create company_relax_period with id %d succesfully", prefix, company_relax_period.ID)
	response.WriteHeaderAndEntity(http.StatusOK, company_relax_period)
	//insert a new row into Summary

	startAt := getDateStr(company_relax_period.StartAt)
	endAt := getDateStr(company_relax_period.EndAt)

	summaries := make([]Summary, 0)
	db.Debug().Find(&summaries, "company_id = ? AND (day between '?' and '?')", company_relax_period.CompanyId, startAt, endAt)
	for _, s := range summaries {
		db.Model(&s).Update("relax_day", 'T')
	}
	return
}

func getDateStr(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d 00:00:00", year, month, day)
}

func (c CompanyRelaxPeriod) updateCompanyRelaxPeriod(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateCompanyRelaxPeriod]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	company_relax_period_id := request.PathParameter("company_relax_period_id")
	company_relax_period := CompanyRelaxPeriod{}
	err := request.ReadEntity(&company_relax_period)

	//fail to parse company_relax_period entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update company_relax_period, err %s", err)
		returnmsg := fmt.Sprintf("无法更新公司休假信息，提供的公司信息解析失败")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//fail to parse company_relax_period id
	id, err := strconv.Atoi(company_relax_period_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update company_relax_period, path company_relax_period_id is %s, err %s", company_relax_period_id, err)
		returnmsg := fmt.Sprintf("无法更新公司休假信息，提供的公司休假id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	if id != company_relax_period.ID {
		errmsg := fmt.Sprintf("cannot update company_relax_period, path company_relax_period_id %d is not equal to id %d in body content", id, company_relax_period.ID)
		returnmsg := fmt.Sprintf("无法更新公司休假信息，提供的公司休假id与URL中的公司休假id不匹配")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	realCompanyRelaxPeriod := CompanyRelaxPeriod{}
	db.Debug().First(&realCompanyRelaxPeriod, company_relax_period.ID)

	//cannot find company_relax_period
	if realCompanyRelaxPeriod.ID == 0 {
		errmsg := fmt.Sprintf("cannot update company_relax_period, company_relax_period_id %d does not exist", company_relax_period.ID)
		returnmsg := fmt.Sprintf("该项公司休假信息已被删除")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	summaries := make([]Summary, 0)
	startAt := getDateStr(realCompanyRelaxPeriod.StartAt)
	endAt := getDateStr(realCompanyRelaxPeriod.EndAt)
	db.Debug().Find(&summaries, "company_id = ? AND (day between '?' and '?')", realCompanyRelaxPeriod.CompanyId, startAt, endAt)
	for _, s := range summaries {
		db.Model(&s).Update("relax_day", 'F')
	}

	//find company_relax_period and update
	db.Debug().Model(&realCompanyRelaxPeriod).Update(company_relax_period)
	glog.Infof("%s update company_relax_period with id %d successfully and return", prefix, realCompanyRelaxPeriod.ID)

	summaries = summaries[:0]
	startAt = getDateStr(company_relax_period.StartAt)
	endAt = getDateStr(company_relax_period.EndAt)
	db.Debug().Find(&summaries, "company_id = ? AND (day between '?' and '?')", company_relax_period.CompanyId, startAt, endAt)
	for _, s := range summaries {
		db.Model(&s).Update("relax_day", 'T')
	}

	response.WriteHeaderAndEntity(http.StatusOK, realCompanyRelaxPeriod)
	return
}

func (c CompanyRelaxPeriod) deleteCompanyRelaxPeriod(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteCompanyRelaxPeriod]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	company_relax_period_id := request.PathParameter("company_relax_period_id")
	id, err := strconv.Atoi(company_relax_period_id)
	//fail to parse company_relax_period id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete company_relax_period, company_relax_period_id is not integer, err %s", err)
		returnmsg := fmt.Sprintf("无法删除公司休假信息，提供的公司id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	company_relax_period := CompanyRelaxPeriod{}
	db.Debug().First(&company_relax_period, id)
	if company_relax_period.ID == 0 {
		//company_relax_period with id doesn't exist, return ok
		glog.Infof("%s company_relax_period with id %s doesn't exist, return ok", prefix, company_relax_period_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	summaries := make([]Summary, 0)
	startAt := getDateStr(company_relax_period.StartAt)
	endAt := getDateStr(company_relax_period.EndAt)
	db.Debug().Find(&summaries, "company_id = ? AND (day between '?' and '?')", company_relax_period.CompanyId, startAt, endAt)
	for _, s := range summaries {
		db.Model(&s).Update("relax_day", 'F')
	}
	db.Debug().Delete(&company_relax_period)

	//delete company_relax_period successfully
	glog.Infof("%s delete company_relax_period with id %s successfully", prefix, company_relax_period_id)
	response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	return
}
