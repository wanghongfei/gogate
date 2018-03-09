package main

import (
	"fmt"
	"os"

	"code.google.com/p/log4go"
	"github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/serv"
)

func main() {
	log4go.Info("start gogate at %s:%d", conf.App.Host, conf.App.Port)

	server, err := serv.NewGatewayServer(conf.App.Host, conf.App.Port, conf.App.RouteConfig, conf.App.MaxConnection)
	checkErrorExit(err)

	err = server.Start()
	checkErrorExit(err)
}

func checkErrorExit(err error) {
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
}
