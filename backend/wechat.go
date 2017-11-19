package main

import (
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
)

type WxResponse struct {
}

func (w WxResponse) Register(container *restful.Container) {
	core.NewServeMux()
	ws := new(restful.WebService)
	ws.Path("/town").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(t.findTown))
	ws.Route(ws.GET("/{town_id}").To(t.findTown))
	ws.Route(ws.GET("/{town_id}/{scope}").To(t.findTown))
	ws.Route(ws.POST("/{town_id}").To(t.updateTown))
	ws.Route(ws.PUT("").To(t.createTown))
	ws.Route(ws.DELETE("/{town_id}").To(t.deleteTown))
	container.Add(ws)
}
