package main

import (
	"flag"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"net/http"
)

func main() {
	flag.Parse()

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

	glog.Infof("starting webserver on localhost:8000")

	mux := core.NewServeMux()

	server := &http.Server{Addr: ":8000", Handler: wsContainer}
	server.ListenAndServe()
}

func RestfulServer() {
}
