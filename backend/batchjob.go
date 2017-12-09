package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/robfig/cron"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/template"
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
		summary := Summary{Day: todayTime, CompanyId: company.ID, CompanyName: company.Name, IsFinish: "F"}
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
			db.Debug().Model(&summary).Update("is_finish", "T")
		} else {
			s := Summary{IsFinish: "F", UnfinishIds: companyStat}
			db.Debug().Model(&summary).Update(s)
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
	db.Debug().Find(&summaries)
	for _, s := range summaries {
		if s.IsFinish == "F" {
			glog.Infof("%s company id %d job is not finish today", prefix, s.CompanyId)
			company := Company{}
			db.Debug().First(&company, s.CompanyId)
			users := make([]User, 0)
			db.Debug().Where("company_id = ?", company.ID).Where("enable = 'T'").Find(&users)
			k1 := Keyword{Value: company.Name}
			k2 := Keyword{Value: todayStrCN}
			k3 := Keyword{Value: s.UnfinishIds}
			first := Keyword{Value: "您好，本日贵企业尚需处理如下拍照任务"}
			remark := Keyword{Value: "请您尽快处理"}
			msg := TMsgData{First: first, Keyword1: k1, Keyword2: k2, Keyword3: k3, Remark: remark}
			t := template.TemplateMessage2{TemplateId: wxTemplateId, Data: msg}
			for _, u := range users {
				t.ToUser = u.WxOpenId
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
	glog.Infof("%s end send message job", prefix)
}

func jobWorker() {
	c := cron.New()
	c.AddFunc("0 0 0 * * *", refreshTodaySummary)
	c.AddFunc("0 2 * * * *", refreshSummary)
	c.AddFunc("0 */30 * * * *", refreshSummaryStat)
	c.AddFunc("0 0 12 * * *", sendTemplateMsg)
	c.AddFunc("0 0 16 * * *", sendTemplateMsg)
	c.AddFunc("0 35 19 * * *", sendTemplateMsg)
	c.AddFunc("0 0 18 * * *", sendTemplateMsg)
	c.AddFunc("0 0 20 * * *", sendTemplateMsg)
	c.AddFunc("0 0 22 * * *", sendTemplateMsg)
	c.Start()
}
