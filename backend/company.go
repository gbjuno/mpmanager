package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
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
	ws.Path(RESTAPIVERSION + "/company").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(c.findCompany))
	ws.Route(ws.GET("/?pageNo={pageNo}&pageSize={pageSize}&order={order}").To(c.findCompany))
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
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")

	var searchCompany *gorm.DB = db.Debug()

	if order != "asc" && order != "desc" {
		errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
		glog.Errorf("%s %s", prefix, errmsg)
		order = "asc"
	}

	if order == "" {
		order = "asc"
	}

	glog.Infof("%s find company with order %s", prefix, order)

	companies := make([]Company, 0)
	count := 0
	searchCompany.Find(&companies).Count(&count)
	searchCompany = searchCompany.Order("id " + order)

	if company_id == "" {
		isPageSizeOk := true
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			isPageSizeOk = false
			errmsg := fmt.Sprintf("cannot find object with pageSize %s, err %s, ignore", pageSize, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		//pageNo depends on pageSize
		isPageNoOk := true
		pageNoInt, err := strconv.Atoi(pageNo)
		if err != nil {
			isPageNoOk = false
			errmsg := fmt.Sprintf("cannot find object with pageNo %s, err %s, ignore", pageNo, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		if isPageSizeOk && isPageNoOk {
			limit := pageSizeInt
			offset := (pageNoInt - 1) * limit
			glog.Infof("%s set find company db with pageSize %s, pageNo %s(limit %d, offset %d)", prefix, pageSize, pageNo, limit, offset)
			searchCompany = searchCompany.Offset(offset).Limit(limit)
		}

		companyList := CompanyList{}
		companyList.Companies = make([]Company, 0)
		searchCompany.Find(&companyList.Companies)

		//companyList.Count = len(companyList.Companies)
		companyList.Count = count
		for i, _ := range companyList.Companies {
			country := Country{}
			db.Debug().First(&country, companyList.Companies[i].CountryId)
			companyList.Companies[i].CountryName = country.Name
		}

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
	db.Debug().First(&company, id)
	//cannot find company
	if company.ID == 0 {
		errmsg := fmt.Sprintf("cannot find company with id %s", company_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	country := Country{}
	db.Debug().First(&country, company.CountryId)
	company.CountryName = country.Name

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
		db.Debug().Model(&company).Related(&userList.Users)
		userList.Count = len(userList.Users)
		for i := range userList.Users {
			userList.Users[i].CompanyName = company.Name
		}
		response.WriteHeaderAndEntity(http.StatusOK, userList)
		glog.Infof("%s return users related company with id %d", prefix, company.ID)
		return
	}

	//find monitor_place related to company
	if scope == "monitorplace" {
		monitorPlaceList := MonitorPlaceWithCompany{}
		monitorPlaceList.CompanyId = company.ID
		monitorPlaceList.MonitorPlaces = make([]MonitorPlace, 0)
		db.Debug().Model(&company).Related(&monitorPlaceList.MonitorPlaces)
		monitorPlaceList.Count = len(monitorPlaceList.MonitorPlaces)
		for i := range monitorPlaceList.MonitorPlaces {
			m := MonitorType{}
			db.Debug().Where("id = ?", monitorPlaceList.MonitorPlaces[i].MonitorTypeId).First(&m)
			monitorPlaceList.MonitorPlaces[i].MonitorTypeName = m.Name
		}
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
		if company.Enable == "F" {
			company.Enable = "F"
		} else {
			company.Enable = "T"
		}
		searchCompany := Company{}
		db.Debug().Where("name = ?", company.Name).First(&searchCompany)
		if searchCompany.ID != 0 {
			errmsg := fmt.Sprintf("company %s already exists", company.Name)
			returnmsg := fmt.Sprintf("存在相同的公司名")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}
		db.Debug().Create(&company)

		if company.ID == 0 {
			//fail to create company on database
			errmsg := fmt.Sprintf("cannot create company on database")
			returnmsg := fmt.Sprintf("无法创建公司，清联系管理员")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		} else {
			//create company on database

			glog.Info("%s create company with id %d succesfully", prefix, company.ID)
			response.WriteHeaderAndEntity(http.StatusOK, company)

			//insert a new row into Summary
			loc, _ := time.LoadLocation("Local")
			timeNow := time.Now()
			todayStr := fmt.Sprintf("%d%02d%02d", timeNow.Year(), timeNow.Month(), timeNow.Day())
			shortForm := "20060102"
			todayTime, _ := time.ParseInLocation(shortForm, todayStr, loc)
			summary := Summary{Day: todayTime, CompanyId: company.ID, IsFinish: "F"}
			glog.Info("%s try to create summary for company with id %d succesfully", prefix, company.ID)
			db.Debug().Create(&summary)

			return
		}
	} else {
		//fail to parse company entity
		errmsg := fmt.Sprintf("cannot create company, err %s", err)
		returnmsg := fmt.Sprintf("无法创建信息，提供的公司信息错误")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
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
		returnmsg := fmt.Sprintf("无法更新公司信息，提供的公司信息解析失败")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//fail to parse company id
	id, err := strconv.Atoi(company_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update company, path company_id is %s, err %s", company_id, err)
		returnmsg := fmt.Sprintf("无法更新公司信息，提供的公司id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	if id != company.ID {
		errmsg := fmt.Sprintf("cannot update company, path company_id %d is not equal to id %d in body content", id, company.ID)
		returnmsg := fmt.Sprintf("无法更新公司信息，提供的公司id与URL中的公司id不匹配")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	realCompany := Company{}
	db.Debug().First(&realCompany, company.ID)

	//cannot find company
	if realCompany.ID == 0 {
		errmsg := fmt.Sprintf("cannot update company, company_id %d does not exist", company.ID)
		returnmsg := fmt.Sprintf("公司已被删除")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	if company.Enable == "F" {
		company.Enable = "F"
	} else {
		company.Enable = "T"
	}

	//find comopany and update
	db.Debug().Model(&realCompany).Update(company)
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
		returnmsg := fmt.Sprintf("无法删除公司，提供的公司id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	company := Company{}
	db.Debug().First(&company, id)
	if company.ID == 0 {
		//company with id doesn't exist, return ok
		glog.Infof("%s company with id %s doesn't exist, return ok", prefix, company_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Debug().Delete(&company)

	realCompany := Company{}
	db.Debug().First(&realCompany, id)

	if realCompany.ID != 0 {
		//fail to delete company
		errmsg := fmt.Sprintf("cannot delete company,some of other object is referencing")
		returnmsg := fmt.Sprintf("无法删除公司，公司仍被引用")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	} else {
		//delete company successfully
		glog.Infof("%s delete company with id %s successfully", prefix, company_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
