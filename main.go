package main

import (
	"fmt"
	"os"

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
