package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/robfig/cron"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/template"
)

var monthDay [13]int = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func getDaysOfMonth(month time.Month, isLeapYear bool) int {
	if int(month) == 2 {
		if isLeapYear {
			return 29
		} else {
			return 28
		}
	}
	return monthDay[int(month)]
}

//全局：连续拍照完成天数，完成拍照天数，拍照完成率
//今年：连续拍照完成天数，完成拍照天数，拍照完成率
//本月：连续拍照完成天数，完成拍照天数，拍照完成率
//最近30天：连续拍照完成天数，完成拍照天数，拍照完成率
func refreshCompanyFinishStat() {
	prefix := fmt.Sprintf("[%s]", "refreshCompanyFinishStat")

	companies := make([]Company, 0)
	companyMap := make(map[int]*Company)
	db.Debug().Find(&companies)
	for index := range companies {
		companyMap[companies[index].ID] = &companies[index]
		companies[index].ContinuousFinishDaysAll = 0
		companies[index].ContinuousFinishDaysInLast365days = 0
		companies[index].ContinuousFinishDaysInLast182days = 0
		companies[index].ContinuousFinishDaysInLast90days = 0
		companies[index].ContinuousFinishDaysInLast30days = 0
		companies[index].MaxContinuousFinishDaysAll = 0
		companies[index].MaxContinuousFinishDaysInLast365days = 0
		companies[index].MaxContinuousFinishDaysInLast182days = 0
		companies[index].MaxContinuousFinishDaysInLast90days = 0
		companies[index].MaxContinuousFinishDaysInLast30days = 0
		companies[index].FinishDaysAll = 0
		companies[index].FinishDaysInLast365days = 0
		companies[index].FinishDaysInLast182days = 0
		companies[index].FinishDaysInLast90days = 0
		companies[index].FinishDaysInLast30days = 0
	}

	timeNow := time.Now()
	thisYear, thisMonth, thisDay := timeNow.Date()
	today := time.Date(thisYear, thisMonth, thisDay, 0, 0, 0, 0, time.Local)
	firstOfThisYear := time.Date(thisYear, time.January, 1, 0, 0, 0, 0, time.Local)
	firstOfThisMonth := time.Date(thisYear, thisMonth, 1, 0, 0, 0, 0, time.Local)
	firstOfLast365days := today.Add(-365 * 24 * time.Duration(time.Second*3600))
	firstOfLast182days := today.Add(-182 * 24 * time.Duration(time.Second*3600))
	firstOfLast90days := today.Add(-90 * 24 * time.Duration(time.Second*3600))
	firstOfLast30days := today.Add(-30 * 24 * time.Duration(time.Second*3600))

	totalDaysInLast365Days := 365.0
	totalDaysInLast182Days := 182.0
	totalDaysInLast90Days := 90.0
	totalDaysInLast30Days := 30.0

	companyMonthStatMap := make(map[int]map[string]*CompanyMonthStat)
	companyYearStatMap := make(map[int]map[string]*CompanyYearStat)

	summaries := make([]Summary, 0)
	db.Debug().Find(&summaries, "(day < ?)", getDateStr(timeNow))
	glog.Infof("%s summary must match (day < %s)", prefix, getDateStr(timeNow))
	glog.Infof("%s get %d matched summary", prefix, len(summaries))
	for _, s := range summaries {
		sYear, sMonth, _ := s.Day.Date()
		sYearStr := fmt.Sprintf("%d", sYear)
		sMonthStr := fmt.Sprintf("%d-%02d", sYear, sMonth)
		if _, ok := companyYearStatMap[s.CompanyId]; !ok {
			companyYearStatMap[s.CompanyId] = make(map[string]*CompanyYearStat)
		}
		sYearLeap := false
		if (sYear%400 == 0) || (sYear%4 == 0 && sYear%100 != 0) {
			sYearLeap = true
		}
		companyCreateAt := companyMap[s.CompanyId].CreateAt
		companyCreateDay := time.Date(companyCreateAt.Year(), companyCreateAt.Month(), companyCreateAt.Day(), 0, 0, 0, 0, time.Local)
		if _, ok := companyYearStatMap[s.CompanyId][sYearStr]; !ok {
			companyYearStat := CompanyYearStat{}
			companyYearStat.CompanyID = s.CompanyId
			companyYearStat.Date = time.Date(s.Day.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
			//4种情况
			if companyCreateDay.Year() == thisYear {
				//1. 今年创建, 统计今年的情况
				companyYearStat.TotalDays = int(today.Sub(companyCreateDay).Seconds() / 3600 / 24)
			} else if thisYear == sYear {
				//2. 非今年创建，统计今年的情况
				companyYearStat.TotalDays = int(today.Sub(firstOfThisYear).Seconds() / 3600 / 24)
			} else if companyCreateDay.Year() == sYear {
				//3. 非今年创建，如果统计的数据和创建的日期在同一年
				if sYearLeap {
					companyYearStat.TotalDays = 366 - companyCreateDay.YearDay() + 1
				} else {
					companyYearStat.TotalDays = 365 - companyCreateDay.YearDay() + 1
				}
			} else {
				//4. 非今年创建，如果统计的数据和创建的日期不在同一年
				if sYearLeap {
					companyYearStat.TotalDays = 366
				} else {
					companyYearStat.TotalDays = 365
				}
			}
			companyYearStatMap[s.CompanyId][sYearStr] = &companyYearStat
		}
		companyYearStat := companyYearStatMap[s.CompanyId][sYearStr]

		if _, ok := companyMonthStatMap[s.CompanyId]; !ok {
			companyMonthStatMap[s.CompanyId] = make(map[string]*CompanyMonthStat)
		}
		if _, ok := companyMonthStatMap[s.CompanyId][sMonthStr]; !ok {
			companyMonthStat := CompanyMonthStat{}
			companyMonthStat.CompanyID = s.CompanyId
			companyMonthStat.Date = time.Date(s.Day.Year(), s.Day.Month(), 1, 0, 0, 0, 0, time.Local)
			//4种情况
			if companyCreateDay.Year() == thisYear && companyCreateDay.Month() == thisMonth {
				//1. 公司本月创建, 统计本月的情况
				companyMonthStat.TotalDays = int(today.Sub(companyCreateDay).Seconds() / 3600 / 24)
			} else if companyCreateDay.Year() == sYear && companyCreateDay.Month() == sMonth {
				//2. 统计公司创建的月份的情况
				companyMonthStat.TotalDays = getDaysOfMonth(sMonth, sYearLeap) - companyCreateDay.Day() + 1
			} else if thisYear == sYear && thisMonth == sMonth {
				//3. 统计本月(未结束)
				companyMonthStat.TotalDays = int(today.Sub(firstOfThisMonth).Seconds() / 3600 / 24)
			} else {
				//4. 统计普通月份（统计的月份不是公司创建的月，也不是本月)
				companyMonthStat.TotalDays = getDaysOfMonth(sMonth, sYearLeap)
			}
			companyMonthStatMap[s.CompanyId][sMonthStr] = &companyMonthStat
		}
		companyMonthStat := companyMonthStatMap[s.CompanyId][sMonthStr]

		if s.RelaxDay == "T" || s.IsFinish == "T" {
			if s.RelaxDay == "T" {
				companyYearStat.RelaxDays++
			}
			companyYearStat.FinishDaysThisYear++
			companyYearStat.ContinuousFinishDaysThisYear++
			if companyYearStat.MaxContinuousFinishDaysThisYear < companyYearStat.ContinuousFinishDaysThisYear {
				companyYearStat.MaxContinuousFinishDaysThisYear = companyYearStat.ContinuousFinishDaysThisYear
			}

			if s.RelaxDay == "T" {
				companyMonthStat.RelaxDays++
			}
			companyMonthStat.FinishDaysThisMonth++
			companyMonthStat.ContinuousFinishDaysThisMonth++
			if companyMonthStat.MaxContinuousFinishDaysThisMonth < companyMonthStat.ContinuousFinishDaysThisMonth {
				companyMonthStat.MaxContinuousFinishDaysThisMonth = companyMonthStat.ContinuousFinishDaysThisMonth
			}

			if s.RelaxDay == "T" {
				companyMap[s.CompanyId].RelaxDaysAll++
			}
			companyMap[s.CompanyId].FinishDaysAll++
			companyMap[s.CompanyId].ContinuousFinishDaysAll++
			if companyMap[s.CompanyId].MaxContinuousFinishDaysAll < companyMap[s.CompanyId].ContinuousFinishDaysAll {
				companyMap[s.CompanyId].MaxContinuousFinishDaysAll = companyMap[s.CompanyId].ContinuousFinishDaysAll
			}

			if !s.Day.Before(firstOfLast365days) {
				if s.RelaxDay == "T" {
					companyMap[s.CompanyId].RelaxDaysInLast365days++
				}
				companyMap[s.CompanyId].FinishDaysInLast365days++
				companyMap[s.CompanyId].ContinuousFinishDaysInLast365days++
				if companyMap[s.CompanyId].MaxContinuousFinishDaysInLast365days < companyMap[s.CompanyId].ContinuousFinishDaysInLast365days {
					companyMap[s.CompanyId].MaxContinuousFinishDaysInLast365days = companyMap[s.CompanyId].ContinuousFinishDaysInLast365days
				}
				if !s.Day.Before(firstOfLast182days) {
					if s.RelaxDay == "T" {
						companyMap[s.CompanyId].RelaxDaysInLast182days++
					}
					companyMap[s.CompanyId].FinishDaysInLast182days++
					companyMap[s.CompanyId].ContinuousFinishDaysInLast182days++
					if companyMap[s.CompanyId].MaxContinuousFinishDaysInLast182days < companyMap[s.CompanyId].ContinuousFinishDaysInLast182days {
						companyMap[s.CompanyId].MaxContinuousFinishDaysInLast182days = companyMap[s.CompanyId].ContinuousFinishDaysInLast182days
					}
					if !s.Day.Before(firstOfLast90days) {
						if s.RelaxDay == "T" {
							companyMap[s.CompanyId].RelaxDaysInLast90days++
						}
						companyMap[s.CompanyId].FinishDaysInLast90days++
						companyMap[s.CompanyId].ContinuousFinishDaysInLast90days++
						if companyMap[s.CompanyId].MaxContinuousFinishDaysInLast90days < companyMap[s.CompanyId].ContinuousFinishDaysInLast90days {
							companyMap[s.CompanyId].MaxContinuousFinishDaysInLast90days = companyMap[s.CompanyId].ContinuousFinishDaysInLast90days
						}
						if !s.Day.Before(firstOfLast30days) {
							if s.RelaxDay == "T" {
								companyMap[s.CompanyId].RelaxDaysInLast30days++
							}
							companyMap[s.CompanyId].FinishDaysInLast30days++
							companyMap[s.CompanyId].ContinuousFinishDaysInLast30days++
							if companyMap[s.CompanyId].MaxContinuousFinishDaysInLast30days < companyMap[s.CompanyId].ContinuousFinishDaysInLast30days {
								companyMap[s.CompanyId].MaxContinuousFinishDaysInLast30days = companyMap[s.CompanyId].ContinuousFinishDaysInLast30days
							}
						}
					}
				}
			}
		} else {
			companyYearStat.ContinuousFinishDaysThisYear = 0
			companyMonthStat.ContinuousFinishDaysThisMonth = 0
			companyMap[s.CompanyId].ContinuousFinishDaysAll = 0
			companyMap[s.CompanyId].ContinuousFinishDaysInLast365days = 0
			companyMap[s.CompanyId].ContinuousFinishDaysInLast182days = 0
			companyMap[s.CompanyId].ContinuousFinishDaysInLast90days = 0
			companyMap[s.CompanyId].ContinuousFinishDaysInLast30days = 0
		}
	}
	glog.Infof("%s scan summary finish", prefix)

	totalMonthRecord := 0
	for companyID := range companyMonthStatMap {
		glog.Infof("%s monthStat companyID=%d start calcuating", prefix, companyID)
		for _, companyMonthStat := range companyMonthStatMap[companyID] {
			totalMonthRecord++
			companyMonthStat.FinishPercentageThisMonth = float64(companyMonthStat.FinishDaysThisMonth) / float64(companyMonthStat.TotalDays)
			glog.Infof("%s monthStat companyID=%d month=%d-%02d totalDays=%d relaxDays=%d finishDays=%d maxContinuousFinishDays=%d continuousFinishDays=%d finishPercentage=%v", prefix, companyMonthStat.CompanyID, companyMonthStat.Date.Year(), companyMonthStat.Date.Month(), companyMonthStat.TotalDays, companyMonthStat.RelaxDays, companyMonthStat.FinishDaysThisMonth, companyMonthStat.MaxContinuousFinishDaysThisMonth, companyMonthStat.ContinuousFinishDaysThisMonth, companyMonthStat.FinishPercentageThisMonth)
			db.Debug().Save(companyMonthStat)
		}
	}
	glog.Infof("%s finish calculating companyMonthStatMap stat, update %d records", prefix, totalMonthRecord)

	totalYearRecord := 0
	for companyID := range companyYearStatMap {
		glog.Infof("%s yearStat companyID=%d start calcuating", prefix, companyID)
		for _, companyYearStat := range companyYearStatMap[companyID] {
			totalYearRecord++
			companyYearStat.FinishPercentageThisYear = float64(companyYearStat.FinishDaysThisYear) / float64(companyYearStat.TotalDays)
			glog.Infof("%s yearStat companyID=%d year=%d totalDays=%d relaxDays=%d finishDays=%d maxContinuousFinishDays=%d continuousFinishDays=%d finishPercentage=%v", prefix, companyYearStat.CompanyID, companyYearStat.Date.Year(), companyYearStat.TotalDays, companyYearStat.RelaxDays, companyYearStat.FinishDaysThisYear, companyYearStat.MaxContinuousFinishDaysThisYear, companyYearStat.ContinuousFinishDaysThisYear, companyYearStat.FinishPercentageThisYear)
			db.Debug().Save(companyYearStat)
		}
	}
	glog.Infof("%s finish calculating companyYearStatMap stat, update %d records", prefix, totalYearRecord)

	bar := int(len(companies) / 10)
	for index := range companies {
		if index%bar == 0 {
			glog.Infof("%s calculating company progressing %d %%", prefix, 10*(index/bar))
		}
		companyCreateAt := companies[index].CreateAt
		companyCreateDay := time.Date(companyCreateAt.Year(), companyCreateAt.Month(), companyCreateAt.Day(), 0, 0, 0, 0, time.Local)

		totalDaysSinceCreate := today.Sub(companyCreateDay).Seconds() / 3600 / 24
		glog.Infof("%s id=%d, finish_days_all=%v, total_days_all=%v", prefix, companies[index].ID, companies[index].FinishDaysAll, totalDaysSinceCreate)
		companies[index].FinishPercentageAll = float64(companies[index].FinishDaysAll) / float64(int(totalDaysSinceCreate))
		companies[index].TotalDaysAll = int(totalDaysSinceCreate)

		if totalDaysSinceCreate < totalDaysInLast365Days {
			totalDaysInLast365Days = totalDaysSinceCreate
		}
		glog.Infof("%s id=%d, finish_days_last_365_days=%v, total_days_last_365_days=%v", prefix, companies[index].ID, companies[index].FinishDaysInLast365days, totalDaysInLast365Days)
		companies[index].FinishPercentageInLast365days = float64(companies[index].FinishDaysInLast365days) / totalDaysInLast365Days
		companies[index].TotalDaysInLast365days = int(totalDaysInLast365Days)

		if totalDaysSinceCreate < totalDaysInLast182Days {
			totalDaysInLast182Days = totalDaysSinceCreate
		}
		glog.Infof("%s id=%d, finish_days_last_182_days=%v, total_days_last_182_days=%v", prefix, companies[index].ID, companies[index].FinishDaysInLast182days, totalDaysInLast182Days)
		companies[index].FinishPercentageInLast182days = float64(companies[index].FinishDaysInLast182days) / totalDaysInLast182Days
		companies[index].TotalDaysInLast182days = int(totalDaysInLast182Days)

		if totalDaysSinceCreate < totalDaysInLast90Days {
			totalDaysInLast90Days = totalDaysSinceCreate
		}
		glog.Infof("%s id=%d, finish_days_last_90_days=%v, total_days_last_90_days=%v", prefix, companies[index].ID, companies[index].FinishDaysInLast90days, totalDaysInLast90Days)
		companies[index].FinishPercentageInLast90days = float64(companies[index].FinishDaysInLast90days) / totalDaysInLast90Days
		companies[index].TotalDaysInLast90days = int(totalDaysInLast90Days)

		if totalDaysSinceCreate < totalDaysInLast30Days {
			totalDaysInLast30Days = totalDaysSinceCreate
		}
		glog.Infof("%s id=%d, finish_days_last_30_days=%v, total_days_last_30_days=%v", prefix, companies[index].ID, companies[index].FinishDaysInLast30days, totalDaysInLast30Days)
		companies[index].FinishPercentageInLast30days = float64(companies[index].FinishDaysInLast30days) / totalDaysInLast30Days
		companies[index].TotalDaysInLast30days = int(totalDaysInLast30Days)

		db.Debug().Model(&companies[index]).Select([]string{
			"relax_days_all", "total_days_all", "max_continuous_finish_days_all", "continuous_finish_days_all", "finish_days_all", "finish_percentage_all",
			"relax_days_last_365_days", "total_days_last_365_days", "max_continuous_finish_days_last_365_days", "continuous_finish_days_last_365_days", "finish_days_last_365_days", "finish_percentage_last_365_days",
			"relax_days_last_182_days", "total_days_last_182_days", "max_continuous_finish_days_last_182_days", "continuous_finish_days_last_182_days", "finish_days_last_182_days", "finish_percentage_last_182_days",
			"relax_days_last_90_days", "total_days_last_90_days", "max_continuous_finish_days_last_90_days", "continuous_finish_days_last_90_days", "finish_days_last_90_days", "finish_percentage_last_90_days",
			"relax_days_last_30_days", "total_days_last_30_days", "max_continuous_finish_days_last_30_days", "continuous_finish_days_last_30_days", "finish_days_last_30_days", "finish_percentage_last_30_days",
		}).Update(&companies[index])
	}
	glog.Infof("%s finish calculating company stat", prefix)
}

//假期更新确保唯一和准确性
func refreshRelaxPeriod() {
	global_relax_period_list := make([]GlobalRelaxPeriod, 0)
	db.Debug().Find(&global_relax_period_list)
	for _, global_relax_period := range global_relax_period_list {
		summaries := make([]Summary, 0)
		startAt := getDateStr(global_relax_period.StartAt)
		endAt := getDateStr(global_relax_period.EndAt)
		db.Debug().Find(&summaries, "(day between '?' and '?')", startAt, endAt)
		for _, s := range summaries {
			db.Model(&s).Update("relax_day", 'T')
		}
	}

	company_relax_period_list := make([]CompanyRelaxPeriod, 0)
	db.Debug().Find(&company_relax_period_list)
	for _, company_relax_period := range company_relax_period_list {
		summaries := make([]Summary, 0)
		startAt := getDateStr(company_relax_period.StartAt)
		endAt := getDateStr(company_relax_period.EndAt)
		db.Debug().Find(&summaries, "company_id = ? AND (day between '?' and '?')", company_relax_period.CompanyId, startAt, endAt)
		for _, s := range summaries {
			db.Model(&s).Update("relax_day", 'T')
		}
	}
}

func refreshTodaySummary() {
	prefix := fmt.Sprintf("[%s]", "refreshTodaySummary")
	loc, _ := time.LoadLocation("Local")
	timeNow := time.Now()
	todayStr := fmt.Sprintf("%d%02d%02d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	shortForm := "20060102"
	todayTime, _ := time.ParseInLocation(shortForm, todayStr, loc)
	condition := fmt.Sprintf("day = str_to_date(%s,'%%Y%%m%%d')", todayStr)
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
			t := TodaySummary{}
			db.Debug().Where(condition).Where("company_id = ?", company.ID).Where("monitor_place_id = ?", monitor_place.ID).First(&t)
			//create (day,companyid,monitor_place_id) new TodaySummary row
			if t.ID == 0 {
				todaySummary := TodaySummary{Day: todayTime, CompanyId: company.ID, CompanyName: company.Name, MonitorPlaceId: monitor_place.ID, MonitorPlaceName: monitor_place.Name, IsUpload: "F", Judgement: "T", EverJudge: "F"}
				glog.Infof("%s try to insert todaySummary for company %d, monitor_place %d, should ignore conflict", prefix, company.ID, monitor_place.ID)
				db.Debug().Create(&todaySummary)
			}
		}
	}
}

//判断是否有新公司，新公司会插入一行
func refreshSummary() {
	prefix := fmt.Sprintf("[%s]", "refreshSummary")
	loc, _ := time.LoadLocation("Local")
	timeNow := time.Now()
	todayStr := fmt.Sprintf("%d%02d%02d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	shortForm := "20060102"
	todayTime, _ := time.ParseInLocation(shortForm, todayStr, loc)
	condition := fmt.Sprintf("day = str_to_date(%s,'%%Y%%m%%d')", todayStr)
	glog.Infof("%s start at time %s %d:%d, %s", prefix, todayStr, timeNow.Hour(), timeNow.Minute(), todayTime)

	companies := make([]Company, 0)
	db.Debug().Where("enable = 'T'").Find(&companies)
	for _, company := range companies {
		s := Summary{}
		db.Debug().Where(condition).Where("company_id = ?", company.ID).First(&s)
		//create (day,companyid,monitor_place_id) new TodaySummary row
		if s.ID == 0 {
			summary := Summary{Day: todayTime, CompanyId: company.ID, CompanyName: company.Name, IsFinish: "F"}
			glog.Infof("%s try to insert summary for company %d, should ignore conflict", prefix, company.ID)
			db.Debug().Create(&summary)
		}
	}
}

//更新公司的上传照片情况
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
		companyStat, ok := companyMap[todaySummary.CompanyId]
		if !ok {
			companyMap[todaySummary.CompanyId] = ""
		}
		if todaySummary.IsUpload == "F" {
			monitorPlace := MonitorPlace{}
			db.Debug().First(&monitorPlace, todaySummary.MonitorPlaceId)
			glog.Infof("%s today summary company id %d, monitor_place %s picture not uploaded", prefix, todaySummary.CompanyId, todaySummary.MonitorPlaceName)
			if companyStat == "" {
				companyMap[todaySummary.CompanyId] = monitorPlace.Name
			} else {
				companyMap[todaySummary.CompanyId] = companyStat + "," + monitorPlace.Name
			}
		}
	}

	for companyId, companyStat := range companyMap {
		summary := Summary{}
		db.Debug().Where(condition).Where("company_id = ?", companyId).First(&summary)
		glog.Infof("%s update summary today for company %d, unfinish_ids %s.", prefix, companyId, companyStat)
		if companyStat == "" {
			summary.IsFinish = "T"
			summary.UnfinishIds = ""
			db.Save(&summary)
		} else {
			summary.IsFinish = "F"
			summary.UnfinishIds = companyStat
			db.Save(&summary)
		}
	}

	summaryList := make([]Summary, 0)
	db.Debug().Where(condition).Find(&summaryList)
	for _, s := range summaryList {
		if s.UnfinishIds == "" {
			s.IsFinish = "T"
			db.Save(&s)
		}
	}
}

type Keyword struct {
	Value string `json:"value"`
}

type TMsgData struct {
	First    Keyword `json:"first"`
	Keyword1 Keyword `json:"keyword1"`
	Keyword2 Keyword `json:"keyword2"`
	Keyword3 Keyword `json:"keyword3"`
	Remark   Keyword `json:"remark"`
}

// send template message
func sendTemplateMsg() {
	prefix := fmt.Sprintf("[%s]", "sendTemplateMsg")
	loc, _ := time.LoadLocation("Local")
	timeNow := time.Now()
	todayStr := fmt.Sprintf("%d%02d%02d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	todayStrCN := fmt.Sprintf("%d年%d月%d日", timeNow.Year(), timeNow.Month(), timeNow.Day())
	shortForm := "20060102"
	todayTime, _ := time.ParseInLocation(shortForm, todayStr, loc)
	glog.Infof("%s start at time %s %d:%d, %s", prefix, todayStr, timeNow.Hour(), timeNow.Minute(), todayTime)

	summaries := make([]Summary, 0)
	condition := fmt.Sprintf("day = str_to_date(%s, '%%Y%%m%%d')", todayStr)
	db.Debug().Where(condition).Find(&summaries)
	for _, s := range summaries {
		if s.IsFinish == "F" {
			glog.Infof("%s company id %d job is not finish today", prefix, s.CompanyId)
			company := Company{}
			db.Debug().First(&company, s.CompanyId)
			users := make([]User, 0)
			db.Debug().Where("company_id = ?", company.ID).Where("enable = 'T'").Find(&users)
			k1 := Keyword{Value: "监控地点拍照"}
			k2 := Keyword{Value: todayStrCN}
			k3 := Keyword{Value: fmt.Sprintf("未完成的拍照地点为%s", s.UnfinishIds)}
			first := Keyword{Value: fmt.Sprintf("您好，本日贵企业%s尚需处理如下安全检查任务", company.Name)}
			remark := Keyword{Value: "请您尽快处理,谢谢!"}
			msg := TMsgData{First: first, Keyword1: k1, Keyword2: k2, Keyword3: k3, Remark: remark}
			t := template.TemplateMessage2{TemplateId: wxTemplateId, Data: msg}
			for _, u := range users {
				if u.WxOpenId != nil {
					t.ToUser = *u.WxOpenId
					tStr, _ := json.Marshal(t)
					msgId, err := template.Send(wechatClient, json.RawMessage(tStr))
					if err != nil {
						glog.Errorf("%s failed to send message to user %s openid %s, message %s,  err %s", prefix, u.Name, t.ToUser, string(tStr), err)
					} else {
						glog.Infof("%s ok to send message to user %s openid %s, msgid %d", prefix, u.Name, t.ToUser, msgId)
					}
				}
			}
		}
	}
	glog.Infof("%s end send message job", prefix)
}

func jobWorker() {
	c := cron.New()
	c.AddFunc("0 0 * * * *", refreshTodaySummary)
	c.AddFunc("0 */2 * * * *", refreshRelaxPeriod)
	c.AddFunc("0 */5 * * * *", refreshCompanyFinishStat)
	c.AddFunc("0 */30 * * * *", refreshSummary)
	c.AddFunc("0 */2 * * * *", refreshSummaryStat)
	c.AddFunc("0 0 14 * * *", sendTemplateMsg)
	c.AddFunc("0 0 16 * * *", sendTemplateMsg)
	c.AddFunc("0 0 18 * * *", sendTemplateMsg)
	c.AddFunc("0 0 20 * * *", sendTemplateMsg)
	c.AddFunc("0 0 22 * * *", sendTemplateMsg)
	c.Start()
}
