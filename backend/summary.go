package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"strconv"
)

type SummaryList struct {
	Count    int       `json:"count"`
	Summarys []Summary `json:"summarys"`
}

func (s Summary) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/summary").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(s.findSummary))
	ws.Route(ws.GET("/{summary_id}").To(s.findSummary))
	ws.Route(ws.POST("").To(s.createSummary))
	ws.Route(ws.PUT("/{summary_id}").To(s.updateSummary))
	ws.Route(ws.DELETE("/{summary_id}").To(s.deleteSummary))
	container.Add(ws)
}

func (s Summary) findSummary(request *restful.Request, response *restful.Response) {
	glog.Infof("GET %s", request.Request.URL)
	summary_id := request.PathParameter("summary_id")

	if summary_id == "" {
		summaryList := SummaryList{}
		summaryList.Summarys = make([]Summary, 0)
		db.Find(&summaryList.Summarys)
		summaryList.Count = len(summaryList.Summarys)
		response.WriteHeaderAndEntity(http.StatusOK, summaryList)
		return
	}

	id, err := strconv.Atoi(summary_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get summary, summary_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	summary := Summary{}
	db.First(&summary, id)
	if summary.ID == 0 {
		errmsg := fmt.Sprintf("cannot find summary with id %s", summary_id)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, summary)
		return
	}
}

func (s Summary) createSummary(request *restful.Request, response *restful.Response) {
	glog.Infof("POST %s", request.Request.URL)
	summary := Summary{}
	err := request.ReadEntity(&summary)
	if err == nil {
		db.Create(&summary)
		response.WriteHeaderAndEntity(http.StatusCreated, summary)
		return
	} else {
		errmsg := fmt.Sprintf("cannot create summary, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (s Summary) updateSummary(request *restful.Request, response *restful.Response) {
	glog.Infof("PUT %s", request.Request.URL)
	summary_id := request.PathParameter("summary_id")
	summary := Summary{}
	err := request.ReadEntity(&summary)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update summary, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	id, err := strconv.Atoi(summary_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update summary, path summary_id is %s, err %s", summary_id, err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != summary.ID {
		errmsg := fmt.Sprintf("cannot update summary, path summary_id %d is not equal to id %d in body content", id, summary.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realSummary := Summary{}
	db.First(&realSummary, summary.ID)
	if realSummary.ID == 0 {
		errmsg := fmt.Sprintf("cannot update summary, summary_id %d is not exist", summary.ID)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	db.Model(&realSummary).Update(summary)
	response.WriteHeaderAndEntity(http.StatusCreated, &realSummary)
	return
}

func (s Summary) deleteSummary(request *restful.Request, response *restful.Response) {
	glog.Infof("DELETE %s", request.Request.URL)
	summary_id := request.PathParameter("summary_id")
	id, err := strconv.Atoi(summary_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete summary, summary_id is not integer, err %s", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	summary := Summary{}
	db.First(&summary, id)
	if summary.ID == 0 {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&summary)

	realSummary := Summary{}
	db.First(&realSummary, id)

	if realSummary.ID != 0 {
		errmsg := fmt.Sprintf("cannot delete summary,some of other object is referencing")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
