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

type SummaryList struct {
	Count    int       `json:"count"`
	Summarys []Summary `json:"summarys"`
}

func (s Summary) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/summary").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(s.findSummary))
	ws.Route(ws.GET("?company_id={company_id}&day={day}&pageSize={pageSize}&pageNo={pageNo}&order={order}").To(s.findSummary))
	container.Add(ws)
}

func (s Summary) findSummary(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findSummary]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	company_id := request.QueryParameter("company_id")
	day := request.QueryParameter("day")
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")

	company := Company{}
	var pageNoInt int
	var pageSizeInt int
	var err error
	var searchDB *gorm.DB = db

	if company_id != "" {
		db.Debug().Where("id = " + company_id).First(&company)
		if company.ID == 0 {
			errmsg := fmt.Sprintf("%s company id %d does not exist", company_id)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		}
		glog.Infof("%s find summary with company id %s", prefix, company_id)
		searchDB = searchDB.Where("id = ?", company.ID)
	}

	if day != "" {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		const shortFormat = "20060102"
		_, err = time.ParseInLocation(shortFormat, day, loc)
		if err != nil {
			errmsg := fmt.Sprintf("cannot find object with after %s, err %s, ignore", day, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		condition := fmt.Sprintf("day = str_to_date(%s, '%%Y%%m%%d')", day)
		glog.Infof("%s find summary with day %s", prefix, day)
		searchDB = searchDB.Where(condition)
	}

	if pageSize != "" && pageNo != "" {
		isPageSizeOk := true
		pageSizeInt, err = strconv.Atoi(pageSize)
		if err != nil {
			isPageSizeOk = false
			errmsg := fmt.Sprintf("cannot find object with pageSize %s, err %s, ignore", pageSize, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		//pageNo depends on pageSize
		isPageNoOk := true
		pageNoInt, err = strconv.Atoi(pageNo)
		if err != nil {
			isPageNoOk = false
			errmsg := fmt.Sprintf("cannot find object with pageNo %s, err %s, ignore", pageNo, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		if isPageSizeOk && isPageNoOk {
			limit := pageSizeInt
			offset := (pageNoInt - 1) * limit
			glog.Infof("%s set find summary db with pageSize %s, pageNo %s(limit %d, offset %d)", prefix, pageSize, pageNo, limit, offset)
			searchDB = searchDB.Offset(offset).Limit(limit)
		}
	}

	if order != "asc" && order != "desc" {
		errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
		glog.Errorf("%s %s", prefix, errmsg)
		order = "desc"
	}

	if order == "" {
		order = "desc"
	}

	glog.Infof("%s find summary with order %s", prefix, order)
	searchDB = searchDB.Order("day " + order)

	//get summary list
	summaryList := SummaryList{}
	summaryList.Summarys = make([]Summary, 0)
	searchDB.Find(&summaryList.Summarys)
	summaryList.Count = len(summaryList.Summarys)
	glog.Infof("%s return all summary list", prefix)
	response.WriteHeaderAndEntity(http.StatusOK, summaryList)
	return
}
