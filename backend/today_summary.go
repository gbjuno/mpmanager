package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type MonitorPlaceSummary struct {
	MonitorPlaceID   int    `json:"monitor_place_id"`
	MonitorPlaceName string `json:"monitor_place_name"`
	IsUpload         string `json:"is_upload"`
	Corrective       string `json:"corrective"`
	EverCorrective   string `json:"ever_corrective"`
}

type CompanySummary struct {
	CompanyID               int                    `json:"company_Id"`
	CompanyName             string                 `json:"company_name"`
	MonitorPlaceSummaryList []*MonitorPlaceSummary `json:"monitor_place_summary"`
}

type CompanySummaryList struct {
	Count          int               `json:"count"`
	CompanySummary []*CompanySummary `json:"company_summary"`
}

func (s TodaySummary) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/today_summary").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{day}/").To(s.findTodaySummary))
	ws.Route(ws.GET("/{day}/?pageNo={pageNo}&pageSize={pageSize}&order={order}").To(s.findTodaySummary))
	ws.Route(ws.GET("/{day}/{company_id}").To(s.findTodaySummary))
	container.Add(ws)
}

func (s TodaySummary) findTodaySummary(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findTodaySummary]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)

	day := request.PathParameter("day")
	company_id := request.PathParameter("company_id")
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")
	company := Company{}

	var searchTodaySummary *gorm.DB = db.Debug()
	var err error

	if day != "" {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		const shortFormat = "20060102"
		_, err = time.ParseInLocation(shortFormat, day, loc)
		if err != nil {
			errmsg := fmt.Sprintf("cannot find object on day %s, err %s", day, err)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		}

		condition := fmt.Sprintf("day = str_to_date(%s, '%%Y%%m%%d')", day)
		glog.Infof("%s find today_summary on day %s", prefix, day)
		searchTodaySummary = searchTodaySummary.Where(condition)
	} else {
		errmsg := "day is not provied"
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if company_id != "" {
		db.Debug().Where("id = " + company_id).First(&company)
		if company.ID == 0 {
			errmsg := fmt.Sprintf("%s company id %d does not exist", company_id)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		}
		glog.Infof("%s find today_summary with company id %s", prefix, company_id)
		searchTodaySummary = searchTodaySummary.Where("company_id = ?", company.ID)
		//todaySummaryList with company id
		todaySummaryList := make([]TodaySummary, 0)
		searchTodaySummary.Find(&todaySummaryList)

		companySummary := CompanySummary{}
		companySummary.CompanyID = company.ID
		companySummary.CompanyName = company.Name
		monitorPlaceSummaryList := make([]*MonitorPlaceSummary, 0)
		for _, t := range todaySummaryList {
			m := MonitorPlaceSummary{MonitorPlaceID: t.MonitorPlaceId, MonitorPlaceName: t.MonitorPlaceName, IsUpload: t.IsUpload, Corrective: t.Corrective, EverCorrective: t.EverCorrective}
			monitorPlaceSummaryList = append(monitorPlaceSummaryList, &m)
		}
		companySummary.MonitorPlaceSummaryList = monitorPlaceSummaryList
		response.WriteHeaderAndEntity(http.StatusOK, companySummary)
		glog.Infof("%s return company_summary with company id %s on day %s", prefix, company_id, day)
		return
	}

	var searchCompany *gorm.DB = db.Debug()
	if order != "asc" && order != "desc" {
		errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
		glog.Errorf("%s %s", prefix, errmsg)
		order = "asc"
	}

	if order == "" {
		order = "asc"
	}

	glog.Infof("%s find company with order %s", prefix, order)
	searchCompany = searchCompany.Order("id " + order)

	var company_CompanySummaryMap = make(map[int]*CompanySummary)
	companies := make([]Company, 0)

	var usePage bool = false
	if pageSize != "" && pageNo != "" {
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
			glog.Infof("%s set find company with pageSize %s, pageNo %s(limit %d, offset %d)", prefix, pageSize, pageNo, limit, offset)
			searchCompany = searchCompany.Offset(offset).Limit(limit)
			usePage = true
		}
	}

	searchCompany.Find(&companies)
	companyIdList := make([]int, 0)
	for _, c := range companies {
		companyIdList = append(companyIdList, c.ID)

		cs := CompanySummary{CompanyName: c.Name, CompanyID: c.ID}
		cs.MonitorPlaceSummaryList = make([]*MonitorPlaceSummary, 0)
		company_CompanySummaryMap[c.ID] = &cs
	}
	//find all today_summary on day
	todaySummaryList := make([]TodaySummary, 0)
	if usePage {
		searchTodaySummary = searchTodaySummary.Where("company_id in (?)", companyIdList)
	}

	searchTodaySummary.Find(&todaySummaryList)
	for _, t := range todaySummaryList {
		cs := company_CompanySummaryMap[t.CompanyId]
		m := MonitorPlaceSummary{MonitorPlaceID: t.MonitorPlaceId, MonitorPlaceName: t.MonitorPlaceName, IsUpload: t.IsUpload, Corrective: t.Corrective, EverCorrective: t.EverCorrective}
		cs.MonitorPlaceSummaryList = append(cs.MonitorPlaceSummaryList, &m)
	}

	csl := CompanySummaryList{}
	csl.Count = len(company_CompanySummaryMap)
	csl.CompanySummary = make([]*CompanySummary, 0)
	for _, c := range companies {
		csl.CompanySummary = append(csl.CompanySummary, company_CompanySummaryMap[c.ID])
	}

	glog.Infof("%s return all today_summary list", prefix)
	response.WriteHeaderAndEntity(http.StatusOK, &csl)
	return
}
