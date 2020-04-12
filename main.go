package main

import (
	. "github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/perr"
	serv "github.com/wanghongfei/gogate/server"
	"os"
)

func main() {
	// 初始化
	serv.InitGogate("gogate.yml")

	// 构造gogate对象
	server, err := serv.NewGatewayServer(
		App.ServerConfig.Host,
		App.ServerConfig.Port,
		App.EurekaConfig.RouteFile,
		App.ServerConfig.MaxConnection,
	)
	checkErrorExit(err, true)

	Log.Infof("pre filters: %v", server.ExportAllPreFilters())
	Log.Infof("post filters: %v", server.ExportAllPostFilters())


	// 启动服务器
	err = server.Start()
	checkErrorExit(err, true)
	Log.Info("listener has been closed")

	// 等待优雅关闭
	err = server.Shutdown()
	checkErrorExit(err, false)
}

func checkErrorExit(err error, exit bool) {
	if nil != err {
		Log.Error(perr.EnvMsg(err))

		if exit {
			os.Exit(1)
		}
	}
}
