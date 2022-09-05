package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"trace2/app/apis"
)

func init() {
	s := g.Server()
	s.Group("/v1", func(group *ghttp.RouterGroup) {
		group.GET("/healthz", apis.Healthz)
		group.GET("/jobs", apis.QueryJob)
	})
	s.Group("/v1/results", func(group *ghttp.RouterGroup) {
		group.GET("/jobs", apis.QueryJob)
	})

}
