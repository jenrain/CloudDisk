package main

import (
	"core/errorx"
	"core/internal/config"
	"core/internal/handler"
	"core/internal/svc"
	"core/mq"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/apiService-api.yaml", "the config file")

func main() {
	flag.Parse()

	//var c config.Config
	conf.MustLoad(*configFile, &config.Conf)

	server := rest.MustNewServer(config.Conf.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(config.Conf)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx.CodeError:
			return http.StatusOK, e.Data()
		default:
			return http.StatusInternalServerError, nil
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", config.Conf.Host, config.Conf.Port)
	// 监听节点心跳
	go mq.ListenHeartbeat()
	// 监听Mysql的binlog日志
	go mq.WatchBinLog(config.Conf)

	server.Start()

}
