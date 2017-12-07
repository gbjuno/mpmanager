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
	ws.Route(ws.GET("/?company_id={company_id}&after={after}&limit={limit}&order={order}").To(s.findSummary))
	container.Add(ws)
}

func (s Summary) findSummary(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findSummary]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	company_id := request.PathParameter("company_id")
	after := request.QueryParameter("after")
	limit := request.QueryParameter("limit")
	order := request.QueryParameter("order")

	company := Company{}
	var limitInt int
	var err error
	var searchDB *gorm.DB

	if company_id != "" {
		db.Where("id = " + company_id).First(&company)
		if company.ID == 0 {
			errmsg := fmt.Sprintf("%s company id %d does not exist", company_id)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		}
		glog.Infof("%s find summary with company id %s", prefix, company_id)
		searchDB.Where("id = ?", company.ID)
	}

	if after != "" {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		const shortFormat = "20060102"
		_, err = time.ParseInLocation(shortFormat, after, loc)
		if err != nil {
			errmsg := fmt.Sprintf("cannot find object with after %s, err %s, ignore", after, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		condition := fmt.Sprintf("day >= str_to_date(%s, '%%%%Y%%%%m%%%%d')", after)
		glog.Infof("%s find summary with after %s", prefix, after)
		searchDB.Where(condition)
	}

	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			errmsg := fmt.Sprintf("cannot find object with limit %s, err %s, ignore", limit, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}
		glog.Infof("%s set find summary db with limit %s", prefix, limit)
		searchDB.Limit(limitInt)
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
	searchDB.Order("day " + order)

	//get summary list
	summaryList := SummaryList{}
	summaryList.Summarys = make([]Summary, 0)
	searchDB.Find(&summaryList.Summarys)
	summaryList.Count = len(summaryList.Summarys)
	glog.Infof("%s return all summary list", prefix)
	response.WriteHeaderAndEntity(http.StatusOK, summaryList)
	return
}
