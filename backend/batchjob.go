package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/robfig/cron"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/template"
)

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
	c.AddFunc("0 */30 * * * *", refreshSummary)
	c.AddFunc("0 */2 * * * *", refreshSummaryStat)
	c.AddFunc("0 0 12 * * *", sendTemplateMsg)
	c.AddFunc("0 0 16 * * *", sendTemplateMsg)
	c.AddFunc("0 0 18 * * *", sendTemplateMsg)
	c.AddFunc("0 0 20 * * *", sendTemplateMsg)
	c.AddFunc("0 0 22 * * *", sendTemplateMsg)
	c.Start()
}
