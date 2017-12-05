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
	prefix := fmt.Sprintf("[%s] [findSummary]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	summary_id := request.PathParameter("summary_id")

	//get summary list
	if summary_id == "" {
		summaryList := SummaryList{}
		summaryList.Summarys = make([]Summary, 0)
		db.Find(&summaryList.Summarys)
		summaryList.Count = len(summaryList.Summarys)
		glog.Infof("%s return summary list", prefix)
		response.WriteHeaderAndEntity(http.StatusOK, summaryList)
		return
	}

	id, err := strconv.Atoi(summary_id)
	//fail to parse summary id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get summary, summary_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	summary := Summary{}
	db.First(&summary, id)
	//cannot find summary
	if summary.ID == 0 {
		errmsg := fmt.Sprintf("cannot find summary with id %s", summary_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	//find summary
	glog.Infof("%s return summary with id %d", prefix, summary.ID)
	response.WriteHeaderAndEntity(http.StatusOK, summary)
	return
}

func (s Summary) createSummary(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createSummary]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	summary := Summary{}
	err := request.ReadEntity(&summary)
	if err == nil {
		db.Create(&summary)
		if summary.ID == 0 {
			//fail to create summary on database
			errmsg := fmt.Sprintf("cannot create summary on database")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
			return
		} else {
			//create summary on database
			glog.Info("%s create summary with id %d succesfully", prefix, summary.ID)
			response.WriteHeaderAndEntity(http.StatusOK, summary)
			return
		}
		return
	} else {
		//fail to parse summary entity
		errmsg := fmt.Sprintf("cannot create summary, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}
}

func (s Summary) updateSummary(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateSummary]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	summary_id := request.PathParameter("summary_id")
	summary := Summary{}
	err := request.ReadEntity(&summary)

	//fail to parse summary entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update summary, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	//fail to parse summary id
	id, err := strconv.Atoi(summary_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update summary, path summary_id is %s, err %s", summary_id, err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	if id != summary.ID {
		errmsg := fmt.Sprintf("cannot update summary, path summary_id %d is not equal to id %d in body content", id, summary.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	realSummary := Summary{}
	db.First(&realSummary, summary.ID)
	//cannot find summary
	if realSummary.ID == 0 {
		errmsg := fmt.Sprintf("cannot update summary, summary_id %d is not exist", summary.ID)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	//find summary and update
	db.Model(&realSummary).Update(summary)
	glog.Infof("%s update summary with id %d successfully and return", prefix, summary.ID)
	response.WriteHeaderAndEntity(http.StatusCreated, realSummary)
	return
}

func (s Summary) deleteSummary(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteSummary]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	summary_id := request.PathParameter("summary_id")
	id, err := strconv.Atoi(summary_id)
	//fail to parse summary id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete summary, summary_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	}

	summary := Summary{}
	db.First(&summary, id)
	if summary.ID == 0 {
		//summary with id doesn't exist
		glog.Infof("%s summary with id %s doesn't exist, return ok", prefix, summary_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Delete(&summary)

	realSummary := Summary{}
	db.First(&realSummary, id)

	if realSummary.ID != 0 {
		//fail to delete summary
		errmsg := fmt.Sprintf("cannot delete summary,some of other object is referencing")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: errmsg})
		return
	} else {
		//delete summary successfully
		glog.Infof("%s delete summary with id %s successfully", prefix, company_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}
}
