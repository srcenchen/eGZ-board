package cmd

import (
	"context"
	"eGZ-Board/internal/controller/count_down"
	"eGZ-Board/internal/controller/device"
	"eGZ-Board/internal/controller/proxy"
	"eGZ-Board/internal/controller/schedule"
	"eGZ-Board/internal/controller/setting"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// 静态资源绑定
			s.SetIndexFolder(true)
			s.AddStaticPath("/resource", "./resource")
			s.SetServerRoot("public")
			s.BindStatusHandler(404, func(r *ghttp.Request) {
				r.Response.ServeFile("public/index.html")
			})
			api := s.Group("/api", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
			})
			api.Group("/setting", func(group *ghttp.RouterGroup) {
				group.Bind(setting.NewV1())
			})
			api.Group("/count-down", func(group *ghttp.RouterGroup) {
				group.Bind(count_down.NewV1())
			})
			api.Group("/schedule", func(group *ghttp.RouterGroup) {
				group.Bind(schedule.NewV1())
			})
			api.Group("/proxy", func(group *ghttp.RouterGroup) {
				group.Bind(proxy.NewV1())
			})
			api.Group("/device", func(group *ghttp.RouterGroup) {
				group.Bind(device.NewV1())
			})
			s.Run()
			return nil
		},
	}
)
