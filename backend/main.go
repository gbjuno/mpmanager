package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-ini/ini"

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
var previewuser string
var recreateMenu bool
var configFile string

const RESTAPIVERSION = "/api/v1"

func EmptyStringError(name string, key string) {
	if key == "" {
		glog.Fatalf("%s is empty string, please provide correct vars", name)
	}
}

func main() {
	flag.StringVar(&dbuser, "user", "root", "database user")
	flag.StringVar(&dbpass, "pass", "123456", "database password")
	flag.StringVar(&dbip, "ip", "127.0.0.1", "database ip address")
	flag.StringVar(&dbport, "port", "3306", "database port")
	flag.StringVar(&dbname, "db", "mpmanager", "database name")
	flag.StringVar(&imgRepo, "imgRepo", "/opt/static", "image save Path")
	flag.StringVar(&domain, "domain", "www.sdlcaj.cn", "domain name")
	flag.BoolVar(&debug, "debug", false, "debug mode, disable weixin init")
	flag.StringVar(&wxport, "wxport", "8001", "wx port")
	flag.StringVar(&restport, "restport", "8000", "rest port")
	flag.StringVar(&previewuser, "previewuser", "", "previewuser")
	flag.StringVar(&configFile, "config", "", "config file path")
	flag.BoolVar(&recreateMenu, "recreatemenu", false, "recreate menu")
	flag.Set("logtostderr", "true")
	flag.Parse()

	InitializeDB(dbuser, dbpass, dbip, dbport, dbname)

	if configFile != "" {
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			glog.Fatalf("cannot open config file %s, config file does not exist!", configFile)
		}

		cfg, err := ini.Load(configFile)
		if err != nil {
			glog.Fatalf("cannot open config file %s, err %s", configFile, err)
		}

		wxAppId = cfg.Section("wechat").Key("appid").MustString("")
		wxAppSecret = cfg.Section("wechat").Key("secret").MustString("")
		wxOriId = cfg.Section("wechat").Key("id").MustString("")
		wxToken = cfg.Section("wechat").Key("token").MustString("")
		wxEncodedAESKey = cfg.Section("wechat").Key("encodedaeskey").MustString("")
		wxTemplateId = cfg.Section("wechat").Key("templateid").MustString("")

		CFGDBUSER := cfg.Section("database").Key("user").MustString("")
		if dbuser == "root" && CFGDBUSER != "" {
			dbuser = CFGDBUSER
		}
		CFGDBPASS := cfg.Section("database").Key("password").MustString("")
		if dbpass == "123456" && CFGDBPASS != "" {
			dbpass = CFGDBPASS
		}
		CFGDBHOST := cfg.Section("database").Key("host").MustString("127.0.0.1")
		if dbip == "127.0.0.1" && CFGDBHOST != "127.0.0.1" {
			dbip = CFGDBHOST
		}
		CFGDBPORT := cfg.Section("database").Key("port").MustString("3306")
		if dbport == "3306" && CFGDBPORT != "3306" {
			dbport = CFGDBPORT
		}
		CFGDBNAME := cfg.Section("database").Key("database").MustString("mpmanager")
		if dbname == "mpmanager" && CFGDBNAME != "mpmanager" {
			dbname = CFGDBNAME
		}

		CFGBASERESTPORT := cfg.Section("base").Key("restport").MustString("8000")
		if restport == "8000" && CFGBASERESTPORT != "8000" {
			restport = CFGBASERESTPORT
		}
		CFGWXPORT := cfg.Section("base").Key("wxport").MustString("8001")
		if wxport == "8001" && CFGWXPORT != "8001" {
			wxport = CFGWXPORT
		}
		CFGBASEIMGREPO := cfg.Section("base").Key("imgrepo").MustString("/opt/static")
		if imgRepo == "/opt/static" && CFGBASEIMGREPO != "/opt/static" {
			imgRepo = CFGBASEIMGREPO
		}
		CFGBASEDOMAIN := cfg.Section("base").Key("domain").MustString("www.sdlcaj.cn")
		if domain == "www.sdlcaj.cn" && CFGBASEDOMAIN != "www.sdlcaj.cn" {
			domain = CFGBASEDOMAIN
		}
		CFGBASERECREATEMENU := cfg.Section("base").Key("recreatemenu").MustBool(false)
		if !recreateMenu {
			recreateMenu = CFGBASERECREATEMENU
		}
	} else {
		glog.Errorf("please provide config file!")
		return
	}

	EmptyStringError("appid", wxAppId)
	EmptyStringError("secret", wxAppSecret)
	EmptyStringError("id", wxOriId)
	EmptyStringError("token", wxToken)
	EmptyStringError("encodedaeskey", wxEncodedAESKey)
	EmptyStringError("templateid", wxTemplateId)

	glog.Infof("%s=%v", "wxAppId", wxAppId)
	glog.Infof("%s=%v", "wxAppSecret", wxAppSecret)
	glog.Infof("%s=%v", "wxOriId", wxOriId)
	glog.Infof("%s=%v", "wxToken", wxToken)
	glog.Infof("%s=%v", "wxEncodedAESKey", wxEncodedAESKey)
	glog.Infof("%s=%v", "wxTemplateId", wxTemplateId)
	glog.Infof("%s=%v", "dbuser", dbuser)
	glog.Infof("%s=%v", "dbpass", dbpass)
	glog.Infof("%s=%v", "dbip", dbip)
	glog.Infof("%s=%v", "dbport", dbport)
	glog.Infof("%s=%v", "restport", restport)
	glog.Infof("%s=%v", "wxport", wxport)
	glog.Infof("%s=%v", "imgRepo", imgRepo)
	glog.Infof("%s=%v", "domain", domain)
	glog.Infof("%s=%v", "recreateMenu", recreateMenu)

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

	materialAudio := MaterialAudio{}
	materialAudio.Register(wsContainer)

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

	companyRelaxPeriod := CompanyRelaxPeriod{}
	companyRelaxPeriod.Register(wsContainer)

	globalRelaxPeriod := GlobalRelaxPeriod{}
	globalRelaxPeriod.Register(wsContainer)

	go func() {
		refreshCompanyFinishStat()
		refreshTodaySummary()
		refreshSummary()
		refreshSummaryStat()

		glog.Infof("starting cronjob system")
		jobWorker()
		glog.Infof("cronjob system end")
	}()

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
	http.HandleFunc("/backend/photolist", photoListHandler)
	http.HandleFunc("/backend/download", downloadHandler)
	http.HandleFunc("/backend/excel", excelHandler)
	glog.Info(http.ListenAndServe(fmt.Sprintf(":%s", wxport), nil))
}
