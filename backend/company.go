package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"strconv"
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
	ws.Path("/company").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(c.findCompany))
	ws.Route(ws.GET("/{company_id}").To(c.findCompany))
	ws.Route(ws.GET("/{company_id}/{scope}").To(c.findCompany))
	ws.Route(ws.POST("/{company_id}").To(c.updateCompany))
	ws.Route(ws.PUT("").To(c.createCompany))
	ws.Route(ws.DELETE("/{company_id}").To(c.deleteCompany))
	container.Add(ws)
}

func (c Company) findCompany(request *restful.Request, response *restful.Response) {
	glog.Infof("GET %s", request.Request.URL)
	company_id := request.PathParameter("company_id")
	scope := request.PathParameter("scope")

	if company_id == "" {
		companyList := CompanyList{}
		companyList.Companies = make([]Company, 0)
		db.Find(&companyList.Companies)
		companyList.Count = len(companyList.Companies)
		response.WriteHeaderAndEntity(http.StatusOK, companyList)
		return
	}

	id, err := strconv.Atoi(company_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get company, company_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	company := Company{}
	db.First(&company, id)
	if company.ID == 0 {
		errmsg := fmt.Sprintf("cannot find company with id %s", company_id)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	if scope == "" {
		response.WriteHeaderAndEntity(http.StatusOK, company)
		return
	}

	if scope == "user" {
		userList := UserListWithCompany{}
		userList.CompanyId = company.ID
		userList.Users = make([]User, 0)
		db.Model(&company).Related(&userList.Users)
		userList.Count = len(userList.Users)
		response.WriteHeaderAndEntity(http.StatusOK, userList)
		return
	}

	if scope == "monitorplace" {
		monitorPlaceList := MonitorPlaceWithCompany{}
		monitorPlaceList.CompanyId = company.ID
		monitorPlaceList.MonitorPlaces = make([]MonitorPlace, 0)
		db.Model(&company).Related(&monitorPlaceList)
		monitorPlaceList.Count = len(monitorPlaceList.MonitorPlaces)
		response.WriteHeaderAndEntity(http.StatusOK, monitorPlaceList)
		return
	}

	errmsg := fmt.Sprintf("cannot find object with scope %s", scope)
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}

func (c Company) createCompany(request *restful.Request, response *restful.Response) {
	glog.Infof("PUT %s", request.Request.URL)
	company := Company{}
	err := request.ReadEntity(&company)
	if err == nil {
		db.Create(&company)
	} else {
		errmsg := fmt.Sprintf("cannot create company, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (c Company) updateCompany(request *restful.Request, response *restful.Response) {
	glog.Infof("POST %s", request.Request.URL)
	company_id := request.PathParameter("company_id")
	company := Company{}
	err := request.ReadEntity(&company)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update company, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(company_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update company, path company_id is %s, err %s", company_id, err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != company.ID {
		errmsg := fmt.Sprintf("cannot update company, path company_id %d is not equal to id %d in body content", id, company.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realCompany := Company{}
	db.First(&realCompany, company.ID)
	if realCompany.ID == 0 {
		errmsg := fmt.Sprintf("cannot update company, company_id %d is not exist", company.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	db.Model(&realCompany).Update(company)
	response.WriteHeaderAndEntity(http.StatusCreated, &realCompany)
}

func (c Company) deleteCompany(request *restful.Request, response *restful.Response) {
	glog.Infof("DELETE %s", request.Request.URL)
	company_id := request.PathParameter("company_id")
	id, err := strconv.Atoi(company_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete company, company_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	company := Company{}
	db.First(&company, id)
	if company.ID == 0 {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&company)

	realCompany := Company{}
	db.First(&realCompany, id)

	if realCompany.ID != 0 {
		errmsg := fmt.Sprintf("cannot delete company,some of other object is referencing")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	}
}