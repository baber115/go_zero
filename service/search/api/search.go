package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"go_zero_demo/service/search/api/internal/config"
	"go_zero_demo/service/search/api/internal/handler"
	"go_zero_demo/service/search/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/search-api.yaml", "the config file")

var port = flag.Int("port", 3333, "the port to listen")

type (
	AnotherService struct{}
	Request        struct {
		User string `json:"user"`
	}
)

func (s *AnotherService) GetToken() string {
	return stringx.Rand()
}

func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("X-Middleware", "static-middleware")
		next(writer, request)
	}
}

func middlewareWithAnotherService(s *AnotherService) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Add("X-Middleware", s.GetToken())
			next(writer, request)
		}
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := httpx.Parse(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	httpx.OkJson(w, "hello, "+req.User)
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 中间件
	server.Use(middleware)
	server.Use(middlewareWithAnotherService(new(AnotherService)))
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/search/greet",
		Handler: handle,
	})
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
