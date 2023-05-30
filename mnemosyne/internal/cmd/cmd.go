package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"mnemosyne/internal/controller"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/v1/log", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					controller.Hello,
					controller.CBase.FieldOpt,
					controller.Log.QueryLog,
					controller.Event.FieldOpts,
				)
			})
			s.Group("/v1/workload", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse) // 必须项
				//group.Bind(
				//	controller.Workload.Aggregations, controller.Workload.QueryWorkload, controller.Log.QueryLog)
				group.Map(g.Map{
					"/aggregations": controller.Workload.Aggregations,
					"/query":        controller.Workload.QueryWorkload,
					"/fields":       controller.Workload.FieldsOpt,
				})

				//
				//group.ALLMap(g.Map{
				//	"": controller.Log.QueryLog,
				//})
			})
			s.Run()
			return nil
		},
	}
)
