package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

type SummaryList struct {
	Count     int       `json:"count"`
	Finish    int       `json:"finish_num"`
	NotFinish int       `json:"not_finish_num"`
	Summaries []Summary `json:"summaries"`
}

func (s Summary) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/summary").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(s.findSummary))
	ws.Route(ws.GET("?from={from}&company_id={company_id}&day={day}&finish={finish}&pageSize={pageSize}&pageNo={pageNo}&order={order}&format={format}").To(s.findSummary))
	container.Add(ws)
}

func (s Summary) findSummary(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findSummary]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	company_id := request.QueryParameter("company_id")
	day := request.QueryParameter("day")
	from := request.QueryParameter("from")
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")
	finish := request.QueryParameter("finish")
	format := request.QueryParameter("format")

	company := Company{}
	var pageNoInt int
	var pageSizeInt int
	var err error
	var searchDB *gorm.DB = db
	var noPageSearchDB *gorm.DB = db.Debug().Model(&Summary{})

	if company_id != "" {
		db.Debug().Where("id = " + company_id).First(&company)
		if company.ID == 0 {
			errmsg := fmt.Sprintf("company id %s does not exist", company_id)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		}
		glog.Infof("%s find summary with company id %s", prefix, company_id)
		searchDB = searchDB.Where("company_id = ?", company.ID)
		noPageSearchDB = noPageSearchDB.Where("id = ?", company.ID)
	}

	if day != "" {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		const shortFormat = "20060102"
		_, err = time.ParseInLocation(shortFormat, day, loc)
		if err != nil {
			errmsg := fmt.Sprintf("cannot find object with after %s, err %s, ignore", day, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		var condition string
		if from == "true" {
			condition = fmt.Sprintf("day >= str_to_date(%s, '%%Y%%m%%d')", day)
		} else {
			condition = fmt.Sprintf("day = str_to_date(%s, '%%Y%%m%%d')", day)
		}
		glog.Infof("%s find summary with day %s", prefix, day)
		searchDB = searchDB.Where(condition)
		noPageSearchDB = noPageSearchDB.Where(condition)
	}

	if finish == "true" {
		glog.Infof("%s find summary where job is finished", prefix)
		condition := "is_finish = 'T'"
		searchDB = searchDB.Where(condition)
		noPageSearchDB = noPageSearchDB.Where(condition)
	} else if finish == "false" {
		glog.Infof("%s find summary where job is not finished", prefix)
		condition := "is_finish = 'F'"
		searchDB = searchDB.Where(condition)
		noPageSearchDB = noPageSearchDB.Where(condition)
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
	summaryList.Summaries = make([]Summary, 0)
	searchDB.Find(&summaryList.Summaries)
	//get count of summaries where day = day and company_id = company_id
	noPageSearchDB.Count(&summaryList.Count)
	noPageSearchDB.Where("is_finish = 'T'").Count(&summaryList.Finish)
	noPageSearchDB.Where("is_finish = 'F'").Count(&summaryList.NotFinish)
	glog.Infof("%s return all summary list", prefix)

	if format != "xlsx" {
		response.WriteHeaderAndEntity(http.StatusOK, summaryList)
	} else {
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "日期")
		f.SetCellValue("Sheet1", "B1", "镇")
		f.SetCellValue("Sheet1", "C1", "村")
		f.SetCellValue("Sheet1", "D1", "公司")
		f.SetCellValue("Sheet1", "E1", "是否完成")
		f.SetCellValue("Sheet1", "F1", "未完成地点")
		for index, s := range summaryList.Summaries {
			if s.IsFinish == "T" {
				f.SetCellValue("Sheet1", fmt.Sprintf("E%d", index+2), "是")
			} else {
				f.SetCellValue("Sheet1", fmt.Sprintf("E%d", index+2), "否")
				f.SetCellValue("Sheet1", fmt.Sprintf("F%d", index+2), s.UnfinishIds)
			}
			company := Company{}
			db.Debug().Where("id = ?", s.CompanyId).First(&company)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%d", index+2), company.Name)
			country := Country{}
			db.Debug().Where("id = ?", company.CountryId).First(&country)
			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", index+2), country.Name)
			town := Town{}
			db.Debug().Where("id = ?", country.TownId).First(&town)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", index+2), town.Name)
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", index+2), fmt.Sprintf("%d%02d%02d", s.Day.Year(), s.Day.Month(), s.Day.Day()))
		}

		t := time.Now()
		saveFileName := fmt.Sprintf("/tmp/summary_%d%02d%02d_%d.xlsx", t.Year(), t.Month(), t.Day(), t.Nanosecond())
		if err := f.SaveAs(saveFileName); err != nil {
			errmsg := fmt.Sprintf("cannot save file, err %s", err)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		xlsx, err := os.Open(saveFileName)
		if err != nil {
			errmsg := fmt.Sprintf("cannot open file, err %s", err)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer xlsx.Close()

		response.Header().Set("Content-Type", "application/octet-stream")
		response.Header().Set("Content-Disposition", "attachment; filename=summary.xlsx")
		response.WriteHeader(http.StatusOK)
		io.Copy(response, xlsx)
		return
	}
	return
}
