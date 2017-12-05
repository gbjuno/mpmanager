package main

import (
	"bytes"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"io/ioutil"
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
	ws.Path(RESTAPIVERSION + "/town").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(t.findTown))
	ws.Route(ws.GET("/{town_id}").To(t.findTown))
	ws.Route(ws.GET("/{town_id}/{scope}").To(t.findTown))
	ws.Route(ws.POST("").To(t.createTown))
	ws.Route(ws.PUT("/{town_id}").To(t.updateTown))
	ws.Route(ws.DELETE("/{town_id}").To(t.deleteTown))
	container.Add(ws)
}

func (t Town) findTown(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findTown]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	town_id := request.PathParameter("town_id")
	scope := request.PathParameter("scope")

	//get town list
	if town_id == "" {
		townList := TownList{}
		townList.Towns = make([]Town, 0)
		db.Find(&townList.Towns)
		townList.Count = len(townList.Towns)
		response.WriteHeaderAndEntity(http.StatusOK, townList)
		glog.Infof("%s return town list", prefix)
		return
	}

	id, err := strconv.Atoi(town_id)
	//fail to parse town id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get town, town_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	town := Town{}
	db.First(&town, id)
	//cannot find town
	if town.ID == 0 {
		errmsg := fmt.Sprintf("cannot find town with id %d", town_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	//find town
	if scope == "" {
		glog.Infof("%s return town with id %d ", prefix, town.ID)
		response.WriteHeaderAndEntity(http.StatusOK, town)
		return
	}

	//find countries related to town
	if scope == "country" {
		countryList := CountryListWithTownID{}
		countryList.TownId = town.ID
		countryList.Countries = make([]Country, 0)
		db.Model(&town).Related(&countryList.Countries)
		countryList.Count = len(countryList.Countries)
		response.WriteHeaderAndEntity(http.StatusOK, countryList)
		glog.Infof("%s return countries related to town id %d", prefix, town.ID)
		return
	}

	//find companies related to town
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
		glog.Infof("%s return companies related to town id %d", prefix, town.ID)
		response.WriteHeaderAndEntity(http.StatusOK, companyList)
		return
	}

	errmsg := fmt.Sprintf("cannot find object with scope %s", scope)
	glog.Errorf("%s %s", prefix, errmsg)
	response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
	return
}

func (t Town) createTown(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createTown]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	town := Town{}
	err := request.ReadEntity(&town)
	if err == nil {
		db.Create(&town)
		if town.ID == 0 {
			//fail to create town on database
			errmsg := fmt.Sprintf("cannot create town on database")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		} else {
			//create town on database
			glog.Info("%s create town with id %d succesfully", prefix, town.ID)
			response.WriteHeaderAndEntity(http.StatusOK, town)
			return
		}
	} else {
		//fail to parse town entity
		errmsg := fmt.Sprintf("cannot create town, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (t Town) updateTown(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateTown]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	town_id := request.PathParameter("town_id")
	town := Town{}
	err := request.ReadEntity(&town)

	//fail to parse town entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update town, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	//fail to parse town id
	id, err := strconv.Atoi(town_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update town, path town_id is %s, err %s", town_id, err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != town.ID {
		errmsg := fmt.Sprintf("cannot update town, path town_id %d is not equal to id %d in body content", id, town.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realTown := Town{}
	db.First(&realTown, town.ID)

	//cannot find town
	if realTown.ID == 0 {
		errmsg := fmt.Sprintf("cannot update town, town_id %d is not exist", town.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
	//find town and update
	db.Model(&realTown).Update(town)
	glog.Infof("%s update town with id %d successfully and return", prefix, realTown.ID)
	response.WriteHeaderAndEntity(http.StatusCreated, realTown)
	return
}

func (t Town) deleteTown(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteTown]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	town_id := request.PathParameter("town_id")
	id, err := strconv.Atoi(town_id)
	//fail to parse town id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete town, town_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	town := Town{}
	db.First(&town, id)
	if town.ID == 0 {
		//town with id doesn't exist, return ok
		glog.Infof("%s town with id %s doesn't exist, return ok", prefix, town_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&town)

	realTown := Town{}
	db.First(&realTown, id)

	if realTown.ID != 0 {
		//fail to delete town
		errmsg := fmt.Sprintf("cannot delete town,some of other object is referencing")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	} else {
		//delete town successfully
		glog.Infof("%s delete town with id %s successfully", prefix, town_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
