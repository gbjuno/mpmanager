package main

import (
	"bytes"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
)

type SummaryList struct {
	Count    int       `json:"count"`
	Summarys []Summary `json:"summarys"`
}

func (s Summary) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/summary").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(s.findSummary))
	ws.Route(ws.GET("/?companyid={companyid}&after={after}&limit={limit}&order={order}").To(s.findSummary))
	container.Add(ws)
}

func (s Summary) findSummary(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findSummary]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	company_id := request.PathParameter("companyid")

	company := Company{}
	var limitInt int
	var afterOk bool
	var limitOk bool
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

	//get summary list
	if company_id == "" {
		summaryList := SummaryList{}
		summaryList.Summarys = make([]Summary, 0)
		db.Find(&summaryList.Summarys)
		summaryList.Count = len(summaryList.Summarys)
		glog.Infof("%s return all summary list", prefix)
		response.WriteHeaderAndEntity(http.StatusOK, summaryList)
		return
	}

	//find summary
	glog.Infof("%s return summary with id %d", prefix, summary.ID)
	response.WriteHeaderAndEntity(http.StatusOK, summary)
	return
}
