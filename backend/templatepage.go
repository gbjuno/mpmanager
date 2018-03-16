package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
	myTemplate "github.com/gbjuno/mpmanager/backend/templates"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

type TemplatePageList struct {
	Count         int            `json:"count"`
	TemplatePages []TemplatePage `json:"templatePages"`
}

type TemplatePageHtml struct {
	Name         string
	HtmlChapters []HtmlChapter
}

type HtmlChapter struct {
	HtmlUrl    string
	PictureUrl string
	Title      string
	Digest     string
}

func (t TemplatePage) HtmlPage(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [HtmlPage]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	templatePage_id := request.PathParameter("templatePage_id")

	id, err := strconv.Atoi(templatePage_id)
	//fail to parse templatePage id
	if err != nil {
		errmsg := fmt.Sprintf("cannot htmlPage templatePage, templatePage_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	templatePage := TemplatePage{}
	db.Debug().First(&templatePage, id)
	//cannot find templatePage
	if templatePage.ID == 0 {
		errmsg := fmt.Sprintf("cannot htmlPage templatePage with id %s", templatePage_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	html := TemplatePageHtml{}
	html.Name = templatePage.Name
	html.HtmlChapters = make([]HtmlChapter, 0)

	for _, chapterid := range strings.Split(templatePage.ChapterIds, ",") {
		chapter := Chapter{}
		db.Debug().Where("id = ?", chapterid).Find(&chapter)
		if chapter.ID == 0 {
			continue
		}
		m := MaterialPicture{}
		db.Debug().Where("media_id = ?", chapter.ThumbMediaId).Find(&m)
		hc := HtmlChapter{HtmlUrl: chapter.Url, PictureUrl: m.Url, Title: chapter.Title, Digest: chapter.Digest}
		html.HtmlChapters = append(html.HtmlChapters, hc)
	}

	templatePageTmpl := template.Must(template.New("templatepagehtml").Parse(myTemplate.TEMPLATEPAGE))
	templatePageTmpl.Execute(response.ResponseWriter, html)

	glog.Infof("%s return template page html successfully", prefix)
	return
}

func (t TemplatePage) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/templatepage").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(t.findTemplatePage))
	ws.Route(ws.GET("/").To(t.findTemplatePage))
	ws.Route(ws.GET("/{templatePage_id}").To(t.findTemplatePage))
	ws.Route(ws.GET("/{templatePage_id}/html").To(t.HtmlPage))
	ws.Route(ws.GET("/?pageNo={pageNo}&pageSize={pageSize}&order={order}").To(t.findTemplatePage))
	ws.Route(ws.POST("").To(t.createTemplatePage))
	ws.Route(ws.PUT("/{templatePage_id}").To(t.updateTemplatePage))
	ws.Route(ws.DELETE("/{templatePage_id}").To(t.deleteTemplatePage))
	container.Add(ws)
}

func (t TemplatePage) findTemplatePage(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findTemplatePage]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	templatePage_id := request.PathParameter("templatePage_id")
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")

	var searchTemplatePage *gorm.DB = db.Debug()

	if order != "asc" && order != "desc" {
		errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
		glog.Errorf("%s %s", prefix, errmsg)
		order = "desc"
	}

	if order == "" {
		order = "desc"
	}

	glog.Infof("%s find templatePage with order %s", prefix, order)

	templatePages := make([]TemplatePage, 0)
	count := 0
	searchTemplatePage.Find(&templatePages).Count(&count)
	searchTemplatePage = searchTemplatePage.Order("id " + order)

	if templatePage_id == "" {
		isPageSizeOk := true
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			isPageSizeOk = false
			errmsg := fmt.Sprintf("cannot find templatePage with pageSize %s, err %s, ignore", pageSize, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		//pageNo depends on pageSize
		isPageNoOk := true
		pageNoInt, err := strconv.Atoi(pageNo)
		if err != nil {
			isPageNoOk = false
			errmsg := fmt.Sprintf("cannot find templatePage with pageNo %s, err %s, ignore", pageNo, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		if isPageSizeOk && isPageNoOk {
			limit := pageSizeInt
			offset := (pageNoInt - 1) * limit
			glog.Infof("%s set find templatePage db with pageSize %s, pageNo %s(limit %d, offset %d)", prefix, pageSize, pageNo, limit, offset)
			searchTemplatePage = searchTemplatePage.Offset(offset).Limit(limit)
		}

		templatePageList := TemplatePageList{}
		templatePageList.TemplatePages = make([]TemplatePage, 0)
		searchTemplatePage.Find(&templatePageList.TemplatePages)

		for index := range templatePageList.TemplatePages {
			templatePageList.TemplatePages[index].ChapterList = make([]Chapter, 0)
			for _, chapterid := range strings.Split(templatePageList.TemplatePages[index].ChapterIds, ",") {
				chapter := Chapter{}
				db.Debug().Where("id = ?", chapterid).Find(&chapter)
				if chapter.ID != 0 {
					templatePageList.TemplatePages[index].ChapterList = append(templatePageList.TemplatePages[index].ChapterList, chapter)
				}
			}
		}

		response.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(response.ResponseWriter)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		enc.Encode(&templatePageList)
		glog.Infof("%s return templatePage list", prefix)
		return
	}

	id, err := strconv.Atoi(templatePage_id)
	//fail to parse templatePage id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get templatePage, templatePage_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	templatePage := TemplatePage{}
	db.Debug().First(&templatePage, id)
	//cannot find templatePage
	if templatePage.ID == 0 {
		errmsg := fmt.Sprintf("cannot find templatePage with id %s", templatePage_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	templatePage.ChapterList = make([]Chapter, 0)
	for _, chapterid := range strings.Split(templatePage.ChapterIds, ",") {
		chapter := Chapter{}
		db.Debug().Where("id = ?", chapterid).Find(&chapter)
		if chapter.ID != 0 {
			templatePage.ChapterList = append(templatePage.ChapterList, chapter)
		}
	}

	response.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(response.ResponseWriter)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(&templatePage)
	glog.Infof("%s find templatePage with id %d", prefix, templatePage.ID)
	return
}

func (c TemplatePage) createTemplatePage(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createTemplatePage]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	templatePage := TemplatePage{}
	err := request.ReadEntity(&templatePage)
	if err != nil {
		errmsg := fmt.Sprintf("cannot create templatePage, err %s", err)
		returnmsg := fmt.Sprintf("无法创建页面模板，提供的信息错误,请联系管理员")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	templatePage.URL = request.Request.URL.Path + "/html"
	glog.Infof("%s create templatePage with id %d succesfully", prefix, templatePage.ID)
	db.Debug().Create(&templatePage)
	response.WriteHeader(http.StatusOK)
	return
}

func (c TemplatePage) updateTemplatePage(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateTemplatePage]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	templatePage_id := request.PathParameter("templatePage_id")
	templatePage := TemplatePage{}
	err := request.ReadEntity(&templatePage)

	//fail to parse templatePage entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update templatePage, err %s", err)
		returnmsg := fmt.Sprintf("无法更新页面模板，提供的信息解析失败")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//fail to parse templatePage id
	id, err := strconv.Atoi(templatePage_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update templatePage, path templatePage_id is %s, err %s", templatePage_id, err)
		returnmsg := fmt.Sprintf("无法更新页面模板，提供的文章id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	if id != templatePage.ID {
		errmsg := fmt.Sprintf("cannot update templatePage, path templatePage_id %d is not equal to id %d in body content", id, templatePage.ID)
		returnmsg := fmt.Sprintf("无法更新页面模板信息，提供的页面模板id与URL中的页面模板id不匹配")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	realTemplatePage := TemplatePage{}
	db.Debug().First(&realTemplatePage, templatePage.ID)

	//cannot find templatePage
	if realTemplatePage.ID == 0 {
		errmsg := fmt.Sprintf("cannot update templatePage, templatePage_id %d does not exist", templatePage.ID)
		returnmsg := fmt.Sprintf("页面模板已被删除,无法更新")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	templatePage.URL = realTemplatePage.URL
	db.Debug().Model(&realTemplatePage).Update(templatePage)
	glog.Infof("%s update templatePage with id %d successfully and return", prefix, realTemplatePage.ID)
	response.WriteHeader(http.StatusOK)
	return
}

func (c TemplatePage) deleteTemplatePage(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteTemplatePage]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	templatePage_id := request.PathParameter("templatePage_id")
	id, err := strconv.Atoi(templatePage_id)
	//fail to parse templatePage id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete templatePage, templatePage_id is not integer, err %s", err)
		returnmsg := fmt.Sprintf("无法删除页面模板，提供的页面模板id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	templatePage := TemplatePage{}
	db.Debug().First(&templatePage, id)
	if templatePage.ID == 0 {
		//templatePage with id doesn't exist, return ok
		glog.Infof("%s templatePage with id %s doesn't exist, return ok", prefix, templatePage_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	db.Debug().Delete(&templatePage)
	//delete templatePage successfully
	glog.Infof("%s delete templatePage with id %s successfully", prefix, templatePage_id)
	response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	return
}
