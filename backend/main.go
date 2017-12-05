package main

import (
	"flag"
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
var debug bool
var imgRepo string

const RESTAPIVERSION = "/api/v1"

func main() {
	flag.StringVar(&dbuser, "user", "root", "database user")
	flag.StringVar(&dbpass, "pass", "123456", "database password")
	flag.StringVar(&dbip, "ip", "127.0.0.1", "database ip address")
	flag.StringVar(&dbport, "port", "3306", "database port")
	flag.StringVar(&imgRepo, "imgRepo", "/opt/static", "image save Path")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Set("logtostderr", "true")
	flag.Parse()

	InitializeDB()

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

	go func() {
		glog.Infof("starting restful webserver on localhost:8000")
		server := &http.Server{Addr: ":8000", Handler: wsContainer}
		glog.Infof(server.ListenAndServe().Error())
	}()

	glog.Infof("starting wechat backend webserver on localhost:8001")
	http.HandleFunc("/backend/wx_callback", wxCallbackHandler)
	http.HandleFunc("/backend/session", sessionHandler)
	http.HandleFunc("/backend/bind", bindingHandler)
	http.HandleFunc("/backend/confirm", confirmHandler)
	http.HandleFunc("/backend/picture", pictureHandler)
	http.HandleFunc("/backend/photo", photoHandler)
	http.HandleFunc("/backend/download", downloadHandler)
	glog.Info(http.ListenAndServe(":8001", nil))
}
