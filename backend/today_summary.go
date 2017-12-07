package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"net/http"
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
	CompanyID               int                   `json:"company_Id"`
	CompanyName             string                `json:"company_name"`
	MonitorPlaceSummaryList []MonitorPlaceSummary `json:"monitor_place_summary"`
}

func (s TodaySummary) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/today_summary").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{day}/").To(s.findTodaySummary))
	ws.Route(ws.GET("/{day}/{company_id}/").To(s.findTodaySummary))
	container.Add(ws)
}

func (s TodaySummary) findTodaySummary(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findTodaySummary]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)

	day := request.PathParameter("day")
	company_id := request.PathParameter("company_id")
	company := Company{}

	var searchDB *gorm.DB
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

		condition := fmt.Sprintf("day = str_to_date(%s, '%%%%Y%%m%%%%d')", day)
		glog.Infof("%s find today_summary on day %s", prefix, day)
		searchDB.Where(condition)
	} else {
		errmsg := "day is not provied"
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if company_id != "" {
		db.Where("id = " + company_id).First(&company)
		if company.ID == 0 {
			errmsg := fmt.Sprintf("%s company id %d does not exist", company_id)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		}
		glog.Infof("%s find today_summary with company id %s", prefix, company_id)
		searchDB.Where("id = ?", company.ID)
		//todaySummaryList with company id
		todaySummaryList := make([]TodaySummary, 0)
		searchDB.Find(&todaySummaryList)

		companySummary := CompanySummary{}
		companySummary.CompanyID = company.ID
		companySummary.CompanyName = company.Name
		monitorPlaceSummaryList := make([]MonitorPlaceSummary, 0)
		companySummary.MonitorPlaceSummaryList = monitorPlaceSummaryList

		for _, t := range todaySummaryList {
			m := MonitorPlaceSummary{MonitorPlaceID: t.MonitorPlaceId, MonitorPlaceName: t.MonitorPlaceName, IsUpload: t.IsUpload, Corrective: t.Corrective, EverCorrective: t.EverCorrective}
			monitorPlaceSummaryList = append(monitorPlaceSummaryList, m)
		}

		response.WriteHeaderAndEntity(http.StatusOK, companySummary)
		glog.Infof("%s return company_summary with company id %s on day %s", prefix, company_id, day)
		return
	}

	glog.Infof("%s return all today_summary list", prefix)
	errmsg := "company_id not provided"
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}
