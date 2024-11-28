package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"meishiedu.com/classin/internal/config"
	"meishiedu.com/classin/internal/handler"
	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile1 = flag.String("f", "etc/classin-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile1, &c)
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()
	httpx.SetErrorHandler(func(err error) (int, any) {
		responseBody := types.ErrorResponse{}
		responseBody.ErrorInfo.ErrorCode = 102
		responseBody.ErrorInfo.ErrorMsg = err.Error()
		return http.StatusOK, responseBody
	})
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
