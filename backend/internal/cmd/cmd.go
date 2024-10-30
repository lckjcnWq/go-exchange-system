package cmd

import (
	"backend/internal/controller"
	"backend/internal/service/ethereum"
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
			s.Run()
			return nil
		},
	}
)
