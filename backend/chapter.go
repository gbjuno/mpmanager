package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"gopkg.in/chanxuehong/wechat.v2/mp/material"
)

type ChapterList struct {
	Count    int       `json:"count"`
	Chapters []Chapter `json:"chapters"`
}

func (c Chapter) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/chapter").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(c.findChapter))
	ws.Route(ws.GET("/?pageNo={pageNo}&pageSize={pageSize}&order={order}").To(c.findChapter))
	ws.Route(ws.GET("/{chapter_id}").To(c.findChapter))
	ws.Route(ws.POST("").To(c.createChapter))
	ws.Route(ws.PUT("/{chapter_id}").To(c.updateChapter))
	ws.Route(ws.DELETE("/{chapter_id}").To(c.deleteChapter))
	container.Add(ws)
}

func (c Chapter) findChapter(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findChapter]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	chapter_id := request.PathParameter("chapter_id")
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")

	var searchChapter *gorm.DB = db.Debug()

	if order != "asc" && order != "desc" {
		errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
		glog.Errorf("%s %s", prefix, errmsg)
		order = "desc"
	}

	if order == "" {
		order = "desc"
	}

	glog.Infof("%s find chapter with order %s", prefix, order)

	chapters := make([]Chapter, 0)
	count := 0
	searchChapter.Find(&chapters).Count(&count)
	searchChapter = searchChapter.Order("id " + order)

	if chapter_id == "" {
		isPageSizeOk := true
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			isPageSizeOk = false
			errmsg := fmt.Sprintf("cannot find chapter with pageSize %s, err %s, ignore", pageSize, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		//pageNo depends on pageSize
		isPageNoOk := true
		pageNoInt, err := strconv.Atoi(pageNo)
		if err != nil {
			isPageNoOk = false
			errmsg := fmt.Sprintf("cannot find chapter with pageNo %s, err %s, ignore", pageNo, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		if isPageSizeOk && isPageNoOk {
			limit := pageSizeInt
			offset := (pageNoInt - 1) * limit
			glog.Infof("%s set find chapter db with pageSize %s, pageNo %s(limit %d, offset %d)", prefix, pageSize, pageNo, limit, offset)
			searchChapter = searchChapter.Offset(offset).Limit(limit)
		}

		chapterList := ChapterList{}
		chapterList.Chapters = make([]Chapter, 0)
		searchChapter.Find(&chapterList.Chapters)

		response.WriteHeaderAndEntity(http.StatusOK, &chapterList)
		glog.Infof("%s return chapter list", prefix)
		return
	}

	id, err := strconv.Atoi(chapter_id)
	//fail to parse chapter id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get chapter, chapter_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	chapter := Chapter{}
	db.Debug().First(&chapter, id)
	//cannot find chapter
	if chapter.ID == 0 {
		errmsg := fmt.Sprintf("cannot find chapter with id %s", chapter_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, &chapter)
	glog.Infof("%s find chapter with id %d", prefix, chapter.ID)
	return
}

func (c Chapter) createChapter(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createChapter]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	chapter := Chapter{}
	err := request.ReadEntity(&chapter)
	if err != nil {
		errmsg := fmt.Sprintf("cannot create chapter, err %s", err)
		returnmsg := fmt.Sprintf("无法创建文章，提供的文章信息错误,请联系管理员")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	article := material.Article{
		Title:            chapter.Title,
		ThumbMediaId:     chapter.ThumbMediaId,
		Author:           chapter.Author,
		Digest:           chapter.Digest,
		ShowCoverPic:     chapter.ShowCoverPic,
		Content:          chapter.Content,
		ContentSourceURL: chapter.ContentSourceUrl,
	}

	wxNews := material.News{}
	wxNews.Articles = []material.Article{article}
	media_id, err := material.AddNews(wechatClient, &wxNews)
	if err != nil {
		errmsg := fmt.Sprintf("cannot create news, err %s", err)
		returnmsg := fmt.Sprintf("无法创建文章,与微信通讯失联,请稍后重试")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	tx := db.Begin()

	news := News{}
	news.Name = chapter.Title
	news.MediaId = media_id
	tx.Debug().Create(&news)
	tx.Debug().Create(&chapter)
	news.ChapterIds = fmt.Sprintf("%d", chapter.ID)
	tx.Debug().Save(&news)
	glog.Infof("%s create news with id %d related to chapter %d succesfully", prefix, news.ID, chapter.ID)
	tx.Commit()

	glog.Infof("%s create chapter with id %d succesfully", prefix, chapter.ID)
	response.WriteHeader(http.StatusOK)
	return
}

func (c Chapter) updateChapter(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateChapter]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	chapter_id := request.PathParameter("chapter_id")
	chapter := Chapter{}
	err := request.ReadEntity(&chapter)

	//fail to parse chapter entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update chapter, err %s", err)
		returnmsg := fmt.Sprintf("无法更新公司信息，提供的公司信息解析失败")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//fail to parse chapter id
	id, err := strconv.Atoi(chapter_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update chapter, path chapter_id is %s, err %s", chapter_id, err)
		returnmsg := fmt.Sprintf("无法更新文章，提供的文章id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	if id != chapter.ID {
		errmsg := fmt.Sprintf("cannot update chapter, path chapter_id %d is not equal to id %d in body content", id, chapter.ID)
		returnmsg := fmt.Sprintf("无法更新公司信息，提供的文章id与URL中的文章id不匹配")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	realChapter := Chapter{}
	db.Debug().First(&realChapter, chapter.ID)

	//cannot find chapter
	if realChapter.ID == 0 {
		errmsg := fmt.Sprintf("cannot update chapter, chapter_id %d does not exist", chapter.ID)
		returnmsg := fmt.Sprintf("文章已被删除,无法更新")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	chapter.NewsId = realChapter.NewsId

	article := material.Article{
		Title:            chapter.Title,
		ThumbMediaId:     chapter.ThumbMediaId,
		Author:           chapter.Author,
		Digest:           chapter.Digest,
		ShowCoverPic:     chapter.ShowCoverPic,
		Content:          chapter.Content,
		ContentSourceURL: chapter.ContentSourceUrl,
	}

	news := News{}
	db.Debug().First(&news, chapter.NewsId)
	if news.ID == 0 {
		wxNews := material.News{}
		wxNews.Articles = []material.Article{article}
		media_id, err := material.AddNews(wechatClient, &wxNews)
		if err != nil {
			errmsg := fmt.Sprintf("cannot create news, err %s", err)
			returnmsg := fmt.Sprintf("无法更新文章,与微信通讯失联,请稍后重试")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}

		tx := db.Begin()
		news.Name = chapter.Title
		news.MediaId = media_id
		news.ChapterIds = fmt.Sprintf("%d", chapter.ID)
		tx.Debug().Create(&news)

		chapter.NewsId = news.ID
		tx.Debug().Model(&realChapter).Update(chapter)
		tx.Commit()

		glog.Infof("%s update chapter with id %d successfully and return", prefix, realChapter.ID)
		response.WriteHeader(http.StatusOK)
		return
	}

	err = material.UpdateNews(wechatClient, news.MediaId, 0, &article)
	if err != nil {
		errmsg := fmt.Sprintf("cannot create news, err %s", err)
		returnmsg := fmt.Sprintf("无法更新文章,与微信通讯失联,请稍后重试")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	db.Debug().Model(&realChapter).Update(chapter)
	glog.Infof("%s update chapter with id %d successfully and return", prefix, realChapter.ID)
	response.WriteHeader(http.StatusOK)
	return
}

func (c Chapter) deleteChapter(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteChapter]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	chapter_id := request.PathParameter("chapter_id")
	id, err := strconv.Atoi(chapter_id)
	//fail to parse chapter id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete chapter, chapter_id is not integer, err %s", err)
		returnmsg := fmt.Sprintf("无法删除文章，提供的文章id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	chapter := Chapter{}
	db.Debug().First(&chapter, id)
	if chapter.ID == 0 {
		//chapter with id doesn't exist, return ok
		glog.Infof("%s chapter with id %s doesn't exist, return ok", prefix, chapter_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	news := News{}
	db.Debug().First(&news, chapter.NewsId)

	if news.ID != 0 {
		err = material.Delete(wechatClient, news.MediaId)
		if err != nil {
			errmsg := fmt.Sprintf("cannot create news, err %s", err)
			returnmsg := fmt.Sprintf("无法删除文章,与微信通讯失联,请稍后重试")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}
	}

	db.Debug().Delete(&chapter)
	//delete chapter successfully
	glog.Infof("%s delete chapter with id %s successfully", prefix, chapter_id)
	response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	return
}
