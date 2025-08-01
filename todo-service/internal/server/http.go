package server

import (
	todoService "todo-service"
	v1 "todo-service/api/helloworld/v1"
	"todo-service/internal/conf"
	"todo-service/internal/service"

	netHttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, todo *service.TodoService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterTodoServiceHTTPServer(srv, todo)

	// srv.HandlePrefix("/", netHttp.FileServer(netHttp.Dir("./front/dist")))
	fileRoute := mux.NewRouter()
	fileRoute.PathPrefix("/dist").Handler(netHttp.FileServer(netHttp.FS(todoService.FrontDist)))

	srv.HandlePrefix("/", fileRoute)

	return srv
}
