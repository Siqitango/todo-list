//go:build wireinject
// +build wireinject

package server

import (
	"todo-service/internal/biz"
	"todo-service/internal/conf"
	"todo-service/internal/data"
	"todo-service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp 生成应用程序依赖
func wireApp(conf *conf.Server, dataConf *conf.Data, logger log.Logger) (*App, func(), error) {
	panic(wire.Build(
		NewApp,
		NewHTTPServer,
		NewGRPCServer,
		service.NewGreeterService,
		service.NewTodoService,
		biz.NewGreeterUsecase,
		biz.NewTodoUsecase,
		data.NewData,
		data.NewGreeterRepo,
		data.NewTodoRepo,
	))
}

// App 应用程序结构体
type App struct {
	HTTPServer *HTTPServer
	GRPCServer *GRPCServer
}

// NewApp 创建应用程序
func NewApp(hs *HTTPServer, gs *GRPCServer) *App {
	return &App{
		HTTPServer: hs,
		GRPCServer: gs,
	}
}