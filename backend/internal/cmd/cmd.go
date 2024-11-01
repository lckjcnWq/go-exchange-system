package cmd

import (
	"backend/internal/controller"
	"backend/internal/service/ethereum"
	"backend/internal/service/ws"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"backend/internal/controller/hello"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 初始化区块链服务
			if err := ethereum.InitBlockchainService(ctx); err != nil {
				g.Log().Fatal(ctx, "Failed to initialize blockchain services:", err)
			}
			defer ethereum.CloseBlockchainServices()
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewV1(),
				)
			})
			priceController, err := controller.NewPriceController()
			if err != nil {
				g.Log().Fatal(context.Background(), "Failed to create price controller:", err)
			}

			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.Bind(
					priceController.GetPrice,
				)
			})
			// 初始化WebSocket管理器
			ws.Init()

			// WebSocket路由
			s.BindHandler("/ws", func(r *ghttp.Request) {
				ws.Get().HandleConnection(r)
			})

			// REST API路由
			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.CORS)
				group.Bind(
					controller.Trade,
					controller.Price,
				)
			})
			s.Run()
			return nil
		},
	}
)
