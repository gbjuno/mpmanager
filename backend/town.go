package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"strconv"
)

type TownList struct {
	Count int    `json:"count"`
	Towns []Town `json:"towns"`
}

type CountryListWithTownID struct {
	TownId    int       `json:"town_id"`
	Count     int       `json:"count"`
	Countries []Country `json:"countries"`
}

type CompanyListWithTownID struct {
	TownId    int       `json:"town_id"`
	Count     int       `json:"count"`
	Companies []Company `json:"companies"`
}

func (t Town) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/town").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(t.findTown))
	ws.Route(ws.GET("/{town_id}").To(t.findTown))
	ws.Route(ws.GET("/{town_id}/{scope}").To(t.findTown))
	ws.Route(ws.POST("").To(t.createTown))
	ws.Route(ws.PUT("/{town_id}").To(t.updateTown))
	ws.Route(ws.DELETE("/{town_id}").To(t.deleteTown))
	container.Add(ws)
}

func (t Town) findTown(request *restful.Request, response *restful.Response) {
	glog.Infof("GET %s", request.Request.URL)
	town_id := request.PathParameter("town_id")
	scope := request.PathParameter("scope")

	if town_id == "" {
		townList := TownList{}
		townList.Towns = make([]Town, 0)
		db.Find(&townList.Towns)
		townList.Count = len(townList.Towns)
		response.WriteHeaderAndEntity(http.StatusOK, townList)
		return
	}

	id, err := strconv.Atoi(town_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get town, town_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	town := Town{}
	db.First(&town, id)
	if town.ID == 0 {
		errmsg := fmt.Sprintf("cannot find town with id %d", town_id)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	if scope == "" {
		response.WriteHeaderAndEntity(http.StatusOK, town)
		return
	}

	if scope == "country" {
		countryList := CountryListWithTownID{}
		countryList.TownId = town.ID
		countryList.Countries = make([]Country, 0)
		db.Model(&town).Related(&countryList.Countries)
		countryList.Count = len(countryList.Countries)
		response.WriteHeaderAndEntity(http.StatusOK, countryList)
		return
	}

	if scope == "company" {
		companyList := CompanyListWithTownID{}
		companyList.TownId = town.ID
		companyList.Companies = make([]Company, 0)

		countryList := CountryListWithTownID{}
		countryList.TownId = town.ID
		countryList.Countries = make([]Country, 0)
		db.Model(&town).Related(&countryList.Countries)
		for _, country := range countryList.Countries {
			companyTempList := make([]Company, 0)
			db.Model(&country).Related(&companyTempList)
			companyList.Companies = append(companyList.Companies, companyTempList...)
		}
		companyList.Count = len(companyList.Companies)
		response.WriteHeaderAndEntity(http.StatusOK, companyList)
		return
	}

	errmsg := fmt.Sprintf("cannot find object with scope %s", scope)
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}

func (t Town) createTown(request *restful.Request, response *restful.Response) {
	glog.Infof("POST %s", request.Request.URL)
	town := Town{}
	err := request.ReadEntity(&town)
	if err == nil {
		db.Create(&town)
		response.WriteHeaderAndEntity(http.StatusCreated, town)
	} else {
		errmsg := fmt.Sprintf("cannot create town, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (t Town) updateTown(request *restful.Request, response *restful.Response) {
	glog.Infof("PUT %s", request.Request.URL)
	town_id := request.PathParameter("town_id")
	town := Town{}
	err := request.ReadEntity(&town)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update town, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(town_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update town, path town_id is %s, err %s", town_id, err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != town.ID {
		errmsg := fmt.Sprintf("cannot update town, path town_id %d is not equal to id %d in body content", id, town.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realTown := Town{}
	db.First(&realTown, town.ID)
	if realTown.ID == 0 {
		errmsg := fmt.Sprintf("cannot update town, town_id %d is not exist", town.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	db.Model(&realTown).Update(town)
	response.WriteHeaderAndEntity(http.StatusCreated, &realTown)
	return
}

func (t Town) deleteTown(request *restful.Request, response *restful.Response) {
	glog.Infof("DELETE %s", request.Request.URL)
	town_id := request.PathParameter("town_id")
	id, err := strconv.Atoi(town_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete town, town_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	town := Town{}
	db.First(&town, id)
	if town.ID == 0 {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&town)

	realTown := Town{}
	db.First(&realTown, id)

	if realTown.ID != 0 {
		errmsg := fmt.Sprintf("cannot delete town,some of other object is referencing")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	}
}
