package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/robfig/cron"
	"time"
)

func refreshTodaySummary() {
	prefix := fmt.Sprintf("[%s]", "refreshTodaySummary")
	timeNow := time.Now()
	todayStr := fmt.Sprintf("%d%d%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	shortForm := "20160102"
	todayTime, _ := time.Parse(shortForm, todayStr)
	glog.Infof("%s start at time %s %d-%d", prefix, todayStr, timeNow.Hour(), timeNow.Minute())

	/*
		todaySummaries := make([]TodaySummary, 0)
		condition := fmt.Sprintf("date = str_to_date(%s,'%%Y%%m%%d'", todayStr)
		glog.Infof("%s search today_summary created today, condition %s", condition)
		db.Where(condition).Find(&todaySummaries)
	*/

	companies := make([]Company, 0)
	db.Where("eanble = 'T'").Find(&companies)
	for _, company := range companies {
		monitor_places := make([]MonitorPlace, 0)
		db.Where("company_id = ?", company.ID).Find(&monitor_places)
		for _, monitor_place := range monitor_places {
			todaySummary := TodaySummary{Day: todayTime, CompanyId: company.ID, MonitorPlaceId: monitor_place.ID, IsUpload: "F", Corrective: "F", EverCorrective: "F"}
			glog.Infof("%s try to insert todaySummary for company %d, monitor_place %d, should ignore conflict", prefix, company.ID, monitor_place.ID)
			db.Create(&todaySummary)
		}
	}
}

func refreshSummary() {
	prefix := fmt.Sprintf("[%s]", "refreshSummary")
	timeNow := time.Now()
	todayStr := fmt.Sprintf("%d%d%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	shortForm := "20160102"
	todayTime, _ := time.Parse(shortForm, todayStr)
	glog.Infof("%s start at time %s %d-%d", prefix, todayStr, timeNow.Hour(), timeNow.Minute())

	companies := make([]Company, 0)
	db.Where("enable = 'T'").Find(&companies)
	for _, company := range companies {
		summary := Summary{Day: todayTime, CompanyId: company.ID, IsFinish: "F"}
		glog.Infof("%s try to insert summary for company %d, should ignore conflict", prefix, company.ID)
		db.Create(&summary)
	}
}

func refreshSummaryStat() {
	prefix := fmt.Sprintf("[%s]", "refreshSummaryStat")
	timeNow := time.Now()
	todayStr := fmt.Sprintf("%d%d%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	shortForm := "20160102"
	todayTime, _ := time.Parse(shortForm, todayStr)
	glog.Infof("%s start at time %s %d-%d", prefix, todayStr, timeNow.Hour(), timeNow.Minute())

	todaySummaries := make([]TodaySummary, 0)
	condition := fmt.Sprintf("date = str_to_date(%s,'%%Y%%m%%d'", todayStr)
	glog.Infof("%s search today_summary created today, condition %s", condition)
	db.Where(condition).Find(&todaySummaries)

	var companyMap map[int]string = make(map[int]string)

	for _, todaySummary := range todaySummaries {
		if todaySummary.IsUpload == "F" {

			companyStat := companyMap[todaySummary.CompanyId]
			monitorPlace := MonitorPlace{}
			db.First(&monitorPlace, todaySummary.MonitorPlaceId)

			if companyStat == "" {
				companyMap[todaySummary.CompanyId] = monitorPlace.Name
			} else {
				companyMap[todaySummary.CompanyId] = companyStat + "," + monitorPlace.Name
			}
		}
	}

	for companyId, companyStat := range companyMap {
		company := Company{}
		db.First(&company, companyId)
		if companyStat == "" {
			db.Model(&company).Update("is_finish", "T", "unfinish_ids", "")
		} else {
			db.Model(&company).Update("is_finish", "F", "unfinish_ids", companyStat)
		}
	}
}

func jobWorker() {
	c := cron.New()
	c.AddFunc("0 0 0 * * *", refreshTodaySummary)
	c.AddFunc("0 5 0 * * *", refreshSummary)
	c.AddFunc("0 0 * * * *", refreshSummary)
	c.AddFunc("0 10 * * * *", refreshSummaryStat)
	c.Start()
}
