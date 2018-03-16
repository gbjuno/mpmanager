package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"gopkg.in/chanxuehong/wechat.v2/mp/material"
)

type NewsList struct {
	Count  int    `json:"count"`
	Newses []News `json:"newses"`
}

func (c News) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path(RESTAPIVERSION + "/news").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON).Filter(PasswordAuthenticate)
	ws.Route(ws.GET("").To(c.findNews))
	ws.Route(ws.GET("/?pageNo={pageNo}&pageSize={pageSize}&order={order}").To(c.findNews))
	ws.Route(ws.GET("/{news_id}").To(c.findNews))
	ws.Route(ws.POST("").To(c.createNews))
	ws.Route(ws.PUT("/{news_id}").To(c.updateNews))
	ws.Route(ws.DELETE("/{news_id}").To(c.deleteNews))
	container.Add(ws)
}

func (c News) findNews(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [findNews]", request.Request.RemoteAddr)
	glog.Infof("%s GET %s", prefix, request.Request.URL)
	news_id := request.PathParameter("news_id")
	pageSize := request.QueryParameter("pageSize")
	pageNo := request.QueryParameter("pageNo")
	order := request.QueryParameter("order")

	var searchNews *gorm.DB = db.Debug()

	if order != "asc" && order != "desc" {
		errmsg := fmt.Sprintf("order %s is not asc or desc, ignore", order)
		glog.Errorf("%s %s", prefix, errmsg)
		order = "desc"
	}

	if order == "" {
		order = "desc"
	}

	glog.Infof("%s find news with order %s", prefix, order)

	newss := make([]News, 0)
	count := 0
	searchNews.Find(&newss).Count(&count)
	searchNews = searchNews.Order("id " + order)

	if news_id == "" {
		isPageSizeOk := true
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			isPageSizeOk = false
			errmsg := fmt.Sprintf("cannot find news with pageSize %s, err %s, ignore", pageSize, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		//pageNo depends on pageSize
		isPageNoOk := true
		pageNoInt, err := strconv.Atoi(pageNo)
		if err != nil {
			isPageNoOk = false
			errmsg := fmt.Sprintf("cannot find news with pageNo %s, err %s, ignore", pageNo, err)
			glog.Errorf("%s %s", prefix, errmsg)
		}

		if isPageSizeOk && isPageNoOk {
			limit := pageSizeInt
			offset := (pageNoInt - 1) * limit
			glog.Infof("%s set find news db with pageSize %s, pageNo %s(limit %d, offset %d)", prefix, pageSize, pageNo, limit, offset)
			searchNews = searchNews.Offset(offset).Limit(limit)
		}

		newsList := NewsList{}
		newsList.Newses = make([]News, 0)
		searchNews.Find(&newsList.Newses)

		for index := range newsList.Newses {
			newsList.Newses[index].ChapterList = make([]Chapter, 0)
			for _, id := range strings.Split(newsList.Newses[index].ChapterIds, ",") {
				temp_chapter := Chapter{}
				db.Debug().Where("id=?", id).Find(&temp_chapter)
				if temp_chapter.ID == 0 {
					temp_chapter.Title = "文章已被删除"
				}
				newsList.Newses[index].ChapterList = append(newsList.Newses[index].ChapterList, temp_chapter)
			}
		}

		response.WriteHeaderAndEntity(http.StatusOK, &newsList)
		glog.Infof("%s return news list", prefix)
		return
	}

	id, err := strconv.Atoi(news_id)
	//fail to parse news id
	if err != nil {
		errmsg := fmt.Sprintf("cannot get news, news_id is not integer, err %s", err)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	news := News{}
	db.Debug().First(&news, id)
	//cannot find news
	if news.ID == 0 {
		errmsg := fmt.Sprintf("cannot find news with id %s", news_id)
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusNotFound, Response{Status: "error", Error: errmsg})
		return
	}

	for _, id := range strings.Split(news.ChapterIds, ",") {
		temp_chapter := Chapter{}
		db.Debug().Where("id=?", id).Find(&temp_chapter)
		if temp_chapter.ID == 0 {
			temp_chapter.Title = "文章已被删除"
		}
		news.ChapterList = append(news.ChapterList, temp_chapter)
	}

	response.WriteHeaderAndEntity(http.StatusOK, &news)
	glog.Infof("%s find news with id %d", prefix, news.ID)
	return
}

func (c News) createNews(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [createNews]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s POST %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	news := News{}
	err := request.ReadEntity(&news)
	if err != nil {
		errmsg := fmt.Sprintf("cannot create news, err %s", err)
		returnmsg := fmt.Sprintf("无法创建消息，提供的信息错误,请联系管理员")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	wxNews := material.News{}
	wxNews.Articles = make([]material.Article, 0)

	for _, id := range strings.Split(news.ChapterIds, ",") {
		temp_chapter := Chapter{}
		db.Debug().Where("id=?", id).Find(&temp_chapter)
		if temp_chapter.ID == 0 {
			errmsg := fmt.Sprintf("cannot create news, err %s", err)
			returnmsg := fmt.Sprintf("无法创建消息,id为%s的文章已被删除,请联系管理员", id)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}
		wxNews.Articles = append(wxNews.Articles, material.Article{
			Title:            temp_chapter.Title,
			ThumbMediaId:     temp_chapter.ThumbMediaId,
			Author:           temp_chapter.Author,
			Digest:           temp_chapter.Digest,
			ShowCoverPic:     temp_chapter.ShowCoverPic,
			Content:          temp_chapter.Content,
			ContentSourceURL: temp_chapter.ContentSourceUrl,
		})
	}

	media_id, err := material.AddNews(wechatClient, &wxNews)
	if err != nil {
		errmsg := fmt.Sprintf("cannot create news, err %s", err)
		returnmsg := fmt.Sprintf("无法创建消息,连接微信服务器失败,请重试")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	returnWxNews, err := material.GetNews(wechatClient, media_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get news, err %s", err)
		returnmsg := fmt.Sprintf("无法创建消息,连接微信服务器失败,请重试")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	for index, id := range strings.Split(news.ChapterIds, ",") {
		temp_chapter := Chapter{}
		db.Debug().Where("id=?", id).Find(&temp_chapter)

		temp_chapter.Url = string(bytes.Replace([]byte(returnWxNews.Articles[index].URL), []byte("\\u0026"), []byte("&"), -1))
		db.Debug().Save(&temp_chapter)
	}

	news.MediaId = media_id
	db.Debug().Create(&news)

	if news.ID == 0 {
		//fail to create news on database
		errmsg := fmt.Sprintf("cannot create news on database")
		returnmsg := fmt.Sprintf("无法创建消息，请联系管理员")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	glog.Infof("%s create news with id %d succesfully", prefix, news.ID)
	response.WriteHeader(http.StatusOK)
	return
}

func (c News) updateNews(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [updateNews]", request.Request.RemoteAddr)
	content, _ := ioutil.ReadAll(request.Request.Body)
	glog.Infof("%s PUT %s, content %s", prefix, request.Request.URL, content)
	newContent := ioutil.NopCloser(bytes.NewBuffer(content))
	request.Request.Body = newContent
	news_id := request.PathParameter("news_id")
	news := News{}
	err := request.ReadEntity(&news)

	//fail to parse news entity
	if err != nil {
		errmsg := fmt.Sprintf("cannot update news, err %s", err)
		returnmsg := fmt.Sprintf("无法更新信息，提供的信息解析失败")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	//fail to parse news id
	id, err := strconv.Atoi(news_id)
	if err != nil {
		errmsg := fmt.Sprintf("cannot update news, path news_id is %s, err %s", news_id, err)
		returnmsg := fmt.Sprintf("无法更新消息，提供的消息id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	if id != news.ID {
		errmsg := fmt.Sprintf("cannot update news, path news_id %d is not equal to id %d in body content", id, news.ID)
		returnmsg := fmt.Sprintf("无法更新公司信息，提供的消息id与URL中的消息id不匹配")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	realNews := News{}
	db.Debug().First(&realNews, news.ID)

	//cannot find news
	if realNews.ID == 0 {
		errmsg := fmt.Sprintf("cannot update news, news_id %d does not exist", news.ID)
		returnmsg := fmt.Sprintf("消息已被删除,无法更新")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	articles := make([]material.Article, 0)

	for _, id := range strings.Split(news.ChapterIds, ",") {
		temp_chapter := Chapter{}
		db.Debug().Where("id=?", id).Find(&temp_chapter)
		if temp_chapter.ID == 0 {
			errmsg := fmt.Sprintf("cannot create news, err %s", err)
			returnmsg := fmt.Sprintf("无法更新消息,id为%s的文章已被删除,请联系管理员", id)
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}
		articles = append(articles, material.Article{
			Title:            temp_chapter.Title,
			ThumbMediaId:     temp_chapter.ThumbMediaId,
			Author:           temp_chapter.Author,
			Digest:           temp_chapter.Digest,
			ShowCoverPic:     temp_chapter.ShowCoverPic,
			Content:          temp_chapter.Content,
			ContentSourceURL: temp_chapter.ContentSourceUrl,
		})
	}

	for index := range articles {
		err := material.UpdateNews(wechatClient, news.MediaId, index, &articles[index])
		if err != nil {
			errmsg := fmt.Sprintf("cannot create news, err %s", err)
			returnmsg := fmt.Sprintf("无法更新消息,请稍后重试")
			glog.Errorf("%s %s", prefix, errmsg)
			response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
			return
		}
	}

	returnWxNews, err := material.GetNews(wechatClient, news.MediaId)
	if err != nil {
		errmsg := fmt.Sprintf("cannot get news, err %s", err)
		returnmsg := fmt.Sprintf("无法更新消息,连接微信服务器失败,请重试")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	for index, id := range strings.Split(news.ChapterIds, ",") {
		temp_chapter := Chapter{}
		db.Debug().Where("id=?", id).Find(&temp_chapter)
		temp_chapter.Url = string(bytes.Replace([]byte(returnWxNews.Articles[index].URL), []byte("\\u0026"), []byte("&"), -1))
		db.Debug().Save(&temp_chapter)
	}

	db.Debug().Model(&realNews).Update(news)
	glog.Infof("%s update news with id %d successfully and return", prefix, realNews.ID)
	response.WriteHeaderAndEntity(http.StatusOK, realNews)
	return
}

func (c News) deleteNews(request *restful.Request, response *restful.Response) {
	prefix := fmt.Sprintf("[%s] [deleteNews]", request.Request.RemoteAddr)
	glog.Infof("%s DELETE %s", prefix, request.Request.URL)
	news_id := request.PathParameter("news_id")
	id, err := strconv.Atoi(news_id)
	//fail to parse news id
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete news, news_id is not integer, err %s", err)
		returnmsg := fmt.Sprintf("无法删除消息，提供的消息id不是整数")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	news := News{}
	db.Debug().First(&news, id)
	if news.ID == 0 {
		//news with id doesn't exist, return ok
		glog.Infof("%s news with id %s doesn't exist, return ok", prefix, news_id)
		response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
		return
	}

	err = material.Delete(wechatClient, news.MediaId)
	if err != nil {
		errmsg := fmt.Sprintf("cannot delete news,err %s", err)
		returnmsg := fmt.Sprintf("无法删除消息,与微信通讯失联,请重试")
		glog.Errorf("%s %s", prefix, errmsg)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, Response{Status: "error", Error: returnmsg})
		return
	}

	db.Debug().Delete(&news)
	//delete news successfully
	glog.Infof("%s delete news with id %s successfully", prefix, news_id)
	response.WriteHeaderAndEntity(http.StatusOK, Response{Status: "success"})
	return
}
