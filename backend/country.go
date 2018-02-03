package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

type CountryList struct {
	Count     int       `json:"count"`
	Countries []Country `json:"countries"`
}

type CompanyListWithCountryID struct {
	CountryId int       `json:"country_id"`
	Count     int       `json:"count"`
	Companies []Company `json:"companies"`
}

func (c Country) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/country").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(c.findCountry))
	ws.Route(ws.GET("/{country_id}").To(c.findCountry))
	ws.Route(ws.GET("/{country_id}/{scope}").To(c.findCountry))
	ws.Route(ws.POST("").To(c.createCountry))
	ws.Route(ws.PUT("/{country_id}").To(c.updateCountry))
	ws.Route(ws.DELETE("/{country_id}").To(c.deleteCountry))
	container.Add(ws)
}

func (c Country) findCountry(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findCountry]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	country_id := request.PathParameter("country_id")
	scope := request.PathParameter("scope")

	//get country list
	if country_id == "" {
		countryList := CountryList{}
		countryList.Countries = make([]Country, 0)
		db.Debug().Find(&countryList.Countries)
		countryList.Count = len(countryList.Countries)
		response.WriteHeaderAndEntity(http.StatusOK, countryList)
		glog.Infof("%s return country list", prefix)
		return
	}

	id, err := strconv.Atoi(country_id)
	//fail to parse country id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get country, country_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}
	country := Country{}
	db.Debug().First(&country, id)
	//cannot find country
	if country.ID == 0 {
		errmsg := fmt.Sprintf("cannot find country with id %s", country_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	//find country
	if scope == "" {
		glog.Infof("%s return country with id %d", prefix, country.ID)
		response.WriteHeaderAndEntity(http.StatusOK, country)
		return
	}

	//find company related to country
	if scope == "company" {
		companyList := CompanyListWithCountryID{}
		companyList.CountryId = country.ID
		companyList.Companies = make([]Company, 0)
		db.Debug().Model(&country).Related(&companyList.Companies)
		companyList.Count = len(companyList.Companies)
		for i, _ := range companyList.Companies {
			companyList.Companies[i].CountryName = country.Name
		}
		glog.Infof("%s return companies related country with id %d", prefix, country.ID)
		response.WriteHeaderAndEntity(http.StatusOK, companyList)
		return
	}

	errmsg := fmt.Sprintf("cannot find object with scope %s", scope)
	glog.Errorf("%s %s", prefix, errmsg)
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}

func (c Country) createCountry(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createCountry]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	country := Country{}
	err := request.ReadEntity(&country)
	if err == nil {
		town := Town{}
		db.Debug().First(&town, country.TownId)
		if town.ID == 0 {
			errmsg := fmt.Sprintf("town id %d not exists", country.TownId)
			returnmsg := fmt.Sprintf("没有找到id为%s的镇，请先提供/选择镇", town.ID)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}

		//whether country name is unique in the same town
		countries := make([]Country, 0)
		db.Debug().Where(" town_id = ?", town.ID).Find(&countries)
		for _, c := range countries {
			if c.Name == country.Name {
				errmsg := fmt.Sprintf("country %s already exists", country.Name)
				returnmsg := fmt.Sprintf("同名的村已存在")
				glog.Errorf("%s %s", prefix, errmsg)
				response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
				return
			}
		}

		db.Debug().Create(&country)
		if country.ID == 0 {
			//fail to create company on database
			errmsg := fmt.Sprintf("cannot create country on database")
			returnmsg := fmt.Sprintf("无法创建村，请联系管理员")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		} else {
			//create country on database
			glog.Info("%s create country with id %d succesfully", prefix, country.ID)
			response.WriteHeaderAndEntity(http.StatusOK, country)
			return
		}
	} else {
		//fail to parse company entity
		errmsg := fmt.Sprintf("cannot create country, err %s", err)
		returnmsg := fmt.Sprintf("无法创建村，提供的村信息无法解析")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}
}

func (c Country) updateCountry(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateCountry]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	country_id := request.PathParameter("country_id")
	country := Country{}
	err := request.ReadEntity(&country)

	//fail to parse country entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update country, err %s", err)
		returnmsg := fmt.Sprintf("无法更新村信息，提供的村信息无法解析")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//fail to parse country id
	id, err := strconv.Atoi(country_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update country, path country_id is %s, err %s", country_id, err)
		returnmsg := fmt.Sprintf("无法更新村信息，提供的村id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	if id != country.ID {
		errmsg := fmt.Sprintf("cannot update country, path country_id %d is not equal to id %d in body content", id, country.ID)
		returnmsg := fmt.Sprintf("无法更新村信息，提供的镇id与URL中的村id不匹配")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	realCountry := Country{}
	db.Debug().First(&realCountry, country.ID)

	//cannot find country
	if realCountry.ID == 0 {
		errmsg := fmt.Sprintf("cannot update country, country_id %d does not exist", country.ID)
		returnmsg := fmt.Sprintf("无法更新村信息，村已被删除")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//find country and update
	db.Debug().Model(&realCountry).Update(country)
	glog.Infof("%s update country with id %d successfully and return", prefix, realCountry.ID)
	response.WriteHeaderAndEntity(http.StatusOK, realCountry)
}

func (c Country) deleteCountry(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteCountry]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	country_id := request.PathParameter("country_id")
	id, err := strconv.Atoi(country_id)
	//fail to parse country id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete country, country_id is not integer, err %s", err)
		returnmsg := fmt.Sprintf("无法删除村信息，提供的村id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	country := Country{}
	db.Debug().First(&country, id)
	if country.ID == 0 {
		//country with id doesn't exist
		glog.Infof("%s country with id %s doesn't exist, return ok", prefix, country_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Debug().Delete(&country)

	realCountry := Country{}
	db.Debug().First(&realCountry, id)

	if realCountry.ID != 0 {
		//fail to delete country
		errmsg := fmt.Sprintf("cannot delete country,some of other object is referencing")
		returnmsg := fmt.Sprintf("无法删除村，村仍被引用")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	} else {
		//delete country successfully
		glog.Infof("%s delete country with id %s successfully", prefix, country_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
