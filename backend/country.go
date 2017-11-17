package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"strconv"
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
	ws.Path("/country/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{country_id}/{scope}").To(c.findCountry))
	ws.Route(ws.POST("/{country_id}").To(c.updateCountry))
	ws.Route(ws.PUT("").To(c.createCountry))
	ws.Route(ws.DELETE("/{country_id}").To(c.deleteCountry))
	container.Add(ws)
}

func (c Country) findCountry(request *restful.Request, response *restful.Response) {
	glog.Infof("GET %s", request.Request.URL)
	country_id := request.PathParameter("country_id")
	scope := request.PathParameter("scope")

	if country_id == "" {
		countryList := CountryList{}
		countryList.Countries = make([]Country, 0)
		db.Find(&countryList.Countries)
		countryList.Count = len(countryList.Countries)
		response.WriteHeaderAndEntity(http.StatusOK, countryList)
		return
	}

	id, err := strconv.Atoi(country_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get country, country_id is not integer, err %", err)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	country := Country{}
	db.First(&country, id)
	if country.ID == 0 {
		errmsg := fmt.Sprintf("cannot find country with id %s", country.ID)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	if scope == "" {
		response.WriteHeaderAndEntity(http.StatusOK, country)
		return
	}

	if scope == "company" {
		companyList := CompanyListWithCountryID{}
		companyList.CountryId = country.ID
		companyList.Companies = make([]Company, 0)
		db.Model(&country).Related(&companyList.Companies)
		companyList.Count = len(companyList.Companies)
		response.WriteHeaderAndEntity(http.StatusOK, companyList)
		return
	}

	errmsg := fmt.Sprintf("cannot find object with scope %s", scope)
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}

func (c Country) createCountry(request *restful.Request, response *restful.Response) {
	country := Country{}
	err := request.ReadEntity(&country)
	if err == nil {
		db.Create(&country)
	} else {
		errmsg := fmt.Sprintf("cannot create country, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (c Country) updateCountry(request *restful.Request, response *restful.Response) {
	country_id := request.PathParameter("country_id")
	country := Country{}
	err := request.ReadEntity(&country)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update country, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(country_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update country, path country_id is %s, err %s", country_id, err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != country.ID {
		errmsg := fmt.Sprintf("cannot update country, path country_id %d is not equal to id %d in body content", id, country.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realCountry := Country{}
	db.First(&realCountry, country.ID)
	if realCountry.ID == 0 {
		errmsg := fmt.Sprintf("cannot update country, country_id %d is not exist", country.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	db.Model(&realCountry).Update(country)
	response.WriteHeaderAndEntity(http.StatusCreated, &realCountry)
}

func (c Country) deleteCountry(request *restful.Request, response *restful.Response) {
	country_id := request.PathParameter("country_id")
	id, err := strconv.Atoi(country_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete country, country_id is not integer, err %", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	country := Country{}
	db.First(&country, id)
	if country.ID == 0 {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&country)

	realCountry := Country{}
	db.First(&realCountry, id)

	if realCountry.ID != 0 {
		errmsg := fmt.Sprintf("cannot delete country,some of other object is referencing")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	}
}
