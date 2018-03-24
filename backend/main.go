package main

import (
	"flag"
	"fmt"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"net/http"
)

var dbuser string
var dbpass string
var dbip string
var dbport string
var dbname string
var debug bool
var imgRepo string
var domain string
var wxport string
var restport string

const RESTAPIVERSION = "/api/v1"

func main() {
	flag.StringVar(&dbuser, "user", "root", "database user")
	flag.StringVar(&dbpass, "pass", "123456", "database password")
	flag.StringVar(&dbip, "ip", "127.0.0.1", "database ip address")
	flag.StringVar(&dbport, "port", "3306", "database port")
	flag.StringVar(&dbname, "db", "mpmanager", "database name")
	flag.StringVar(&imgRepo, "imgRepo", "/opt/static", "image save Path")
	flag.StringVar(&domain, "domain", "www.juntengshoes.cn", "domain name")
	flag.BoolVar(&debug, "debug", false, "debug mode, disable weixin init")
	flag.StringVar(&wxport, "wxport", "8001", "wx port")
	flag.StringVar(&restport, "restport", "8000", "rest port")
	flag.Set("logtostderr", "true")
	flag.Parse()

	InitializeDB(dbuser, dbpass, dbip, dbport, dbname)

	if debug {
		glog.Info("DEBUG MODE")
	}

	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})

	town := Town{}
	town.Register(wsContainer)

	country := Country{}
	country.Register(wsContainer)

	company := Company{}
	company.Register(wsContainer)

	user := User{}
	user.Register(wsContainer)

	monitor_type := MonitorType{}
	monitor_type.Register(wsContainer)

	monitor_place := MonitorPlace{}
	monitor_place.Register(wsContainer)

	picture := Picture{}
	picture.Register(wsContainer)

	summary := Summary{}
	summary.Register(wsContainer)

	todaySummary := TodaySummary{}
	todaySummary.Register(wsContainer)

	menu := Menu{}
	menu.Register(wsContainer)

	materialPicture := MaterialPicture{}
	materialPicture.Register(wsContainer)

	materialVideo := MaterialVideo{}
	materialVideo.Register(wsContainer)

	mediaPicture := MediaPicture{}
	mediaPicture.Register(wsContainer)

	chapter := Chapter{}
	chapter.Register(wsContainer)

	news := News{}
	news.Register(wsContainer)

	groupsend := GroupSend{}
	groupsend.Register(wsContainer)

	templatePage := TemplatePage{}
	templatePage.Register(wsContainer)

	go func() {
		glog.Infof("starting cronjob system")
		jobWorker()
		glog.Infof("cronjob system end")
	}()

	refreshTodaySummary()
	refreshSummary()
	refreshSummaryStat()

	go func() {
		glog.Infof("starting restful webserver on localhost:%s", restport)
		server := &http.Server{Addr: fmt.Sprintf(":%s", restport), Handler: wsContainer}
		glog.Infof(server.ListenAndServe().Error())
	}()

	if !debug {
		WechatBackendInit()
	}

	glog.Infof("starting wechat backend webserver on localhost:%s", wxport)
	http.HandleFunc("/backend/wx_callback", wxCallbackHandler)
	http.HandleFunc("/backend/binding", bindingHandler)
	http.HandleFunc("/backend/confirm", confirmHandler)
	http.HandleFunc("/backend/scanqrcode", scanqrcodeHandler)
	http.HandleFunc("/backend/companystat", companystatHandler)
	http.HandleFunc("/backend/photo", photoHandler)
	http.HandleFunc("/backend/download", downloadHandler)
	http.HandleFunc("/backend/excel", excelHandler)
	glog.Info(http.ListenAndServe(fmt.Sprintf(":%s", wxport), nil))

}
