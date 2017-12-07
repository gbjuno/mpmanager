package main

import (
	"bytes"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type CompanyList struct {
	Count     int       `json:"count"`
	Companies []Company `json:"companies"`
}

type UserListWithCompany struct {
	CompanyId int    `json:"company_id"`
	Count     int    `json:"count"`
	Users     []User `json:"users"`
}

type MonitorPlaceWithCompany struct {
	CompanyId     int            `json:"company_id"`
	Count         int            `json:"count"`
	MonitorPlaces []MonitorPlace `json:"monitor_places"`
}

func (c Company) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/company").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(c.findCompany))
	ws.Route(ws.GET("/?limit={limit}&offset={offset}").To(c.findCompany))
	ws.Route(ws.GET("/{company_id}").To(c.findCompany))
	ws.Route(ws.GET("/{company_id}/{scope}").To(c.findCompany))
	ws.Route(ws.POST("").To(c.createCompany))
	ws.Route(ws.PUT("/{company_id}").To(c.updateCompany))
	ws.Route(ws.DELETE("/{company_id}").To(c.deleteCompany))
	container.Add(ws)
}

func (c Company) findCompany(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findCompany]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	company_id := request.PathParameter("company_id")
	scope := request.PathParameter("scope")
	offset := request.QueryParameter("offset")
	limit := request.QueryParameter("limit")

	//get company list
	var offsetOk = false
	var limitOk = false
	var offsetInt int
	var limitInt int
	var err error

	if company_id == "" {
		if offset != "" {
			offsetInt, err = strconv.Atoi(offset)
			if err != nil {
				glog.Errorf("%s offset %s is not integer, ignore", prefix, offset)
			}
			offsetOk = true
		}

		if limit != "" {
			limitInt, err = strconv.Atoi(limit)
			if err != nil {
				glog.Errorf("%s limit %s is not integer, ignore", prefix, limit)
			}
			limitOk = true
		}

		companyList := CompanyList{}
		companyList.Companies = make([]Company, 0)
		if offsetOk && limitOk {
			glog.Infof("%s get company list, offset %d limit %d", prefix, offsetInt, limitInt)
			db.Offset(offsetInt).Limit(limitInt).Find(&companyList.Companies)
		} else if offsetOk {
			limitInt = int(^uint(0) >> 1)
			glog.Infof("%s get company list, offset %d limit %d", prefix, offsetInt, limitInt)
			db.Offset(offsetInt).Limit(limitInt).Find(&companyList.Companies)
		} else if limitOk {
			glog.Infof("%s get company list, limit %d", prefix, limitInt)
			db.Limit(limitInt).Find(&companyList.Companies)
		} else {
			glog.Infof("%s get company list", prefix)
			db.Find(&companyList.Companies)
		}

		companyList.Count = len(companyList.Companies)
		response.WriteHeaderAndEntity(http.StatusOK, companyList)
		glog.Infof("%s return company list", prefix)
		return
	}

	id, err := strconv.Atoi(company_id)
	//fail to parse company id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get company, company_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	company := Company{}
	db.First(&company, id)
	//cannot find company
	if company.ID == 0 {
		errmsg := fmt.Sprintf("cannot find company with id %s", company_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	//find company
	if scope == "" {
		glog.Infof("%s return company with id %d", prefix, company.ID)
		response.WriteHeaderAndEntity(http.StatusOK, company)
		return
	}

	//find user related to company
	if scope == "user" {
		userList := UserListWithCompany{}
		userList.CompanyId = company.ID
		userList.Users = make([]User, 0)
		db.Model(&company).Related(&userList.Users)
		userList.Count = len(userList.Users)
		response.WriteHeaderAndEntity(http.StatusOK, userList)
		glog.Infof("%s return users related company with id %d", prefix, company.ID)
		return
	}

	//find monitor_place related to company
	if scope == "monitorplace" {
		monitorPlaceList := MonitorPlaceWithCompany{}
		monitorPlaceList.CompanyId = company.ID
		monitorPlaceList.MonitorPlaces = make([]MonitorPlace, 0)
		db.Model(&company).Related(&monitorPlaceList)
		monitorPlaceList.Count = len(monitorPlaceList.MonitorPlaces)
		response.WriteHeaderAndEntity(http.StatusOK, monitorPlaceList)
		glog.Infof("%s return monitor_places related company with id %d", prefix, company.ID)
		return
	}

	errmsg := fmt.Sprintf("cannot find object with scope %s", scope)
	glog.Errorf("%s %s", prefix, errmsg)
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}

func (c Company) createCompany(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createCompany]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	company := Company{}
	err := request.ReadEntity(&company)
	if err == nil {
		if company.Enable != "" && company.Enable != "F" {
			company.Enable = "T"
		} else {
			company.Enable = "F"
		}
		searchCompany := Company{}
		db.Where("name = ?", company.Name).First(&searchCompany)
		if searchCompany.ID != 0 {
			errmsg := fmt.Sprintf("company %s already exists", company.Name)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		}
		db.Create(&company)

		if company.ID == 0 {
			//fail to create company on database
			errmsg := fmt.Sprintf("cannot create company on database")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		} else {
			//create company on database

			glog.Info("%s create company with id %d succesfully", prefix, company.ID)
			response.WriteHeaderAndEntity(http.StatusOK, company)

			//insert a new row into Summary
			timeNow := time.Now()
			todayStr := fmt.Sprintf("%d%d%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
			shortForm := "20160102"
			todayTime, _ := time.Parse(shortForm, todayStr)
			companies := make([]Company, 0)
			summary := Summary{Day: todayTime, CompanyId: company.ID, IsFinish: "F"}
			glog.Info("%s try to create summary for company with id %d succesfully", prefix, company.ID)
			db.Create(&summary)

			return
		}
	} else {
		//fail to parse company entity
		errmsg := fmt.Sprintf("cannot create company, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

}

func (c Company) updateCompany(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateCompany]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	company_id := request.PathParameter("company_id")
	company := Company{}
	err := request.ReadEntity(&company)

	//fail to parse company entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update company, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	//fail to parse company id
	id, err := strconv.Atoi(company_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update company, path company_id is %s, err %s", company_id, err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != company.ID {
		errmsg := fmt.Sprintf("cannot update company, path company_id %d is not equal to id %d in body content", id, company.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realCompany := Company{}
	db.First(&realCompany, company.ID)

	//cannot find company
	if realCompany.ID == 0 {
		errmsg := fmt.Sprintf("cannot update company, company_id %d does not exist", company.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if company.Enable != "" && company.Enable != "F" {
		company.Enable = "T"
	} else {
		company.Enable = "F"
	}

	//find comopany and update
	db.Model(&realCompany).Update(company)
	glog.Infof("%s update company with id %d successfully and return", prefix, realCompany.ID)
	response.WriteHeaderAndEntity(http.StatusOK, realCompany)
	return
}

func (c Company) deleteCompany(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteCompany]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	company_id := request.PathParameter("company_id")
	id, err := strconv.Atoi(company_id)
	//fail to parse company id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete company, company_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	company := Company{}
	db.First(&company, id)
	if company.ID == 0 {
		//company with id doesn't exist, return ok
		glog.Infof("%s company with id %s doesn't exist, return ok", prefix, company_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&company)

	realCompany := Company{}
	db.First(&realCompany, id)

	if realCompany.ID != 0 {
		//fail to delete company
		errmsg := fmt.Sprintf("cannot delete company,some of other object is referencing")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	} else {
		//delete company successfully
		glog.Infof("%s delete company with id %s successfully", prefix, company_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
