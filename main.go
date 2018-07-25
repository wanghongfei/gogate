package main

import (
	"fmt"
	"os"
	"time"

	asynclog "github.com/alecthomas/log4go"
	"github.com/wanghongfei/gogate/conf"
	serv "github.com/wanghongfei/gogate/server"
)

func main() {
	serv.InitGogate("gogate.yml", "log.xml")

	server, err := serv.NewGatewayServer(
		conf.App.ServerConfig.Host,
		conf.App.ServerConfig.Port,
		conf.App.EurekaConfig.RouteFile,
		conf.App.ServerConfig.MaxConnection,
		// 是否启用优雅关闭
		true,
		// 优雅关闭最大等待时间, 上一个参数为true时有效
		time.Second * 30,
	)
	checkErrorExit(err)

	asynclog.Info("pre filters: %v", server.ExportAllPreFilters())
	asynclog.Info("post filters: %v", server.ExportAllPostFilters())

	asynclog.Info("started gogate at %s:%d", conf.App.ServerConfig.Host, conf.App.ServerConfig.Port)
	err = server.Start()

	checkErrorExit(err)
}

func checkErrorExit(err error) {
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
}
