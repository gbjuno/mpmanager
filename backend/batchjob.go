package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/robfig/cron"
	"time"
)

func refreshTodaySummary() {
	prefix := fmt.Sprintf("[%s]", "refreshTodaySummary")
	loc, _ := time.LoadLocation("Local")
	timeNow := time.Now()
	todayStr := fmt.Sprintf("%d%02d%02d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	shortForm := "20060102"
	todayTime, _ := time.ParseInLocation(shortForm, todayStr, loc)
	glog.Infof("%s start at time %s %d:%d, %s", prefix, todayStr, timeNow.Hour(), timeNow.Minute(), todayTime)

	/*
		todaySummaries := make([]TodaySummary, 0)
		condition := fmt.Sprintf("date = str_to_date(%s,'%%Y%%m%%d')", todayStr)
		glog.Infof("%s search today_summary created today, condition %s", condition)
		db.Debug().Where(condition).Find(&todaySummaries)
	*/

	companies := make([]Company, 0)
	db.Debug().Where("enable = 'T'").Find(&companies)
	for _, company := range companies {
		monitor_places := make([]MonitorPlace, 0)
		db.Debug().Where("company_id = ?", company.ID).Find(&monitor_places)
		for _, monitor_place := range monitor_places {
			todaySummary := TodaySummary{Day: todayTime, CompanyId: company.ID, CompanyName: company.Name, MonitorPlaceId: monitor_place.ID, MonitorPlaceName: monitor_place.Name, IsUpload: "F", Corrective: "F", EverCorrective: "F"}
			glog.Infof("%s try to insert todaySummary for company %d, monitor_place %d, should ignore conflict", prefix, company.ID, monitor_place.ID)
			db.Debug().Create(&todaySummary)
		}
	}
}

func refreshSummary() {
	prefix := fmt.Sprintf("[%s]", "refreshSummary")
	loc, _ := time.LoadLocation("Local")
	timeNow := time.Now()
	todayStr := fmt.Sprintf("%d%02d%02d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	shortForm := "20060102"
	todayTime, _ := time.ParseInLocation(shortForm, todayStr, loc)
	glog.Infof("%s start at time %s %d:%d, %s", prefix, todayStr, timeNow.Hour(), timeNow.Minute(), todayTime)

	companies := make([]Company, 0)
	db.Debug().Where("enable = 'T'").Find(&companies)
	for _, company := range companies {
		summary := Summary{Day: todayTime, CompanyId: company.ID, IsFinish: "F"}
		glog.Infof("%s try to insert summary for company %d, should ignore conflict", prefix, company.ID)
		db.Debug().Create(&summary)
	}
}

func refreshSummaryStat() {
	prefix := fmt.Sprintf("[%s]", "refreshSummaryStat")
	timeNow := time.Now()
	todayStr := fmt.Sprintf("%d%02d%02d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	glog.Infof("%s start at time %s %d-%d", prefix, todayStr, timeNow.Hour(), timeNow.Minute())

	todaySummaries := make([]TodaySummary, 0)
	condition := fmt.Sprintf("day = str_to_date(%s,'%%Y%%m%%d')", todayStr)
	glog.Infof("%s search today_summary created today, condition %s", prefix, condition)
	db.Debug().Where(condition).Find(&todaySummaries)

	var companyMap map[int]string = make(map[int]string)

	for _, todaySummary := range todaySummaries {
		if todaySummary.IsUpload == "F" {

			companyStat := companyMap[todaySummary.CompanyId]
			monitorPlace := MonitorPlace{}
			db.Debug().First(&monitorPlace, todaySummary.MonitorPlaceId)

			if companyStat == "" {
				companyMap[todaySummary.CompanyId] = monitorPlace.Name
			} else {
				companyMap[todaySummary.CompanyId] = companyStat + "," + monitorPlace.Name
			}
		}
	}

	for companyId, companyStat := range companyMap {
		company := Company{}
		db.Debug().First(&company, companyId)
		if companyStat == "" {
			db.Debug().Model(&company).Update("is_finish", "T", "unfinish_ids", "")
		} else {
			db.Debug().Model(&company).Update("is_finish", "F", "unfinish_ids", companyStat)
		}
	}
}

func jobWorker() {
	c := cron.New()
	c.AddFunc("0 0 0 * * *", refreshTodaySummary)
	c.AddFunc("0 2 * * * *", refreshSummary)
	c.AddFunc("0 5 * * * *", refreshSummaryStat)
	c.Start()
}
