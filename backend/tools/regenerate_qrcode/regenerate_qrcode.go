package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gbjuno/mpmanager/backend/tools"
	"github.com/gbjuno/mpmanager/backend/utils"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

var dbuser string
var dbpass string
var dbip string
var dbport string
var dbname string
var imgRepo string
var domain string
var companyid string

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbuser, dbpass, dbip, dbport, dbname))
	if err != nil {
		glog.Errorf("cannot initialize database connection, err %s", err)
		return
	}
}

func main() {
	flag.StringVar(&dbuser, "user", "root", "database user")
	flag.StringVar(&dbpass, "pass", "123456", "database password")
	flag.StringVar(&dbip, "ip", "127.0.0.1", "database ip address")
	flag.StringVar(&dbport, "port", "3306", "database port")
	flag.StringVar(&dbname, "db", "mpmanager", "database name")
	flag.StringVar(&imgRepo, "imgRepo", "/opt/static", "image save Path")
	flag.StringVar(&domain, "domain", "www.juntengshoes.cn", "domain name")
	flag.StringVar(&companyid, "companyid", "all", "company id")
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	initDB()

	if companyid != "all" {
		RegenerateForCompany(companyid)
	} else {
		RegenerateForAllCompany()
	}

	return
}

func RegenerateForAllCompany() {
	glog.Infof("regenerateing qrcode pictures for all company")
	companyList := make([]tools.Company, 0)
	db.Find(&companyList)
	total := 0
	for _, company := range companyList {
		glog.Infof("regenerateing qrcode pictures for company %s", company.Name)
		//create monitor_place qrcode image
		monitorPlaceList := make([]tools.MonitorPlace, 0)
		db.Where("company_id = ?", company.ID).Find(&monitorPlaceList)
		for _, place := range monitorPlaceList {
			qrcodePath := fmt.Sprintf("https://%s/backend/photo?place=%d", domain, place.ID)
			imagePath := (imgRepo + place.QrcodePath)
			imageBaseDir := imagePath[:strings.LastIndex(imagePath, "/")+1]
			if _, err := os.Stat(imageBaseDir); os.IsNotExist(err) {
				os.MkdirAll(imageBaseDir, 0755)
			}
			if err := utils.GenerateQrcodeImage(qrcodePath, company.Name+place.Name, imagePath); err != nil {
				errmsg := fmt.Sprintf("cannot update company %s, monitor_place %s, generate qrcode failed, err %s", company.Name, place.Name, err)
				glog.Errorf("%s", errmsg)
				return
			} else {
				glog.Infof("regenerateing qrcode for company %s monitor place %s on %s", company.Name, place.Name, imagePath)
			}
		}
		total += len(monitorPlaceList)
		glog.Infof("regenerating %d qrcode pictures for company %s", len(monitorPlaceList), company.Name)
	}

	glog.Infof("summary: regenerating %d qrcode pictures for all company", total)
	return
}

func RegenerateForCompany(companyid string) {
	glog.Infof("regenerateing qrcode pictures for company id %s", companyid)
	company := tools.Company{}
	db.Where("id = ?", companyid).First(&company)
	if company.ID == 0 {
		errmsg := fmt.Sprintf("cannot update qrcode image, company with id %s not found", companyid)
		glog.Errorf("%s", errmsg)
		return
	}

	glog.Infof("regenerateing qrcode pictures for company %s", company.Name)

	//create monitor_place qrcode image
	monitorPlaceList := make([]tools.MonitorPlace, 0)
	db.Where("company_id = ?", companyid).Find(&monitorPlaceList)
	for _, place := range monitorPlaceList {
		qrcodePath := fmt.Sprintf("https://%s/backend/photo?place=%d", domain, place.ID)
		imagePath := (imgRepo + place.QrcodePath)
		imageBaseDir := imagePath[:strings.LastIndex(imagePath, "/")+1]
		if _, err := os.Stat(imageBaseDir); os.IsNotExist(err) {
			os.MkdirAll(imageBaseDir, 0755)
		}
		if err := utils.GenerateQrcodeImage(qrcodePath, company.Name+place.Name, imagePath); err != nil {
			errmsg := fmt.Sprintf("cannot update company %s, monitor_place %s, generate qrcode failed, err %s", company.Name, place.Name, err)
			glog.Errorf("%s", errmsg)
			return
		} else {
			glog.Infof("regenerateing qrcode for company %s monitor place %s on %s", company.Name, place.Name, imagePath)
		}
	}
	glog.Infof("regenerating %d qrcode pictures for company %s", len(monitorPlaceList), company.Name)
	return
}
