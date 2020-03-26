package main

import (
	"os"
	"time"

	. "github.com/wanghongfei/gogate/conf"
	serv "github.com/wanghongfei/gogate/server"
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
		// 是否启用优雅关闭
		true,
		// 优雅关闭最大等待时间, 上一个参数为true时有效
		time.Second*1,
	)
	checkErrorExit(err, true)

	Log.Infof("pre filters: %v", server.ExportAllPreFilters())
	Log.Infof("post filters: %v", server.ExportAllPostFilters())

	// deferClose(server, time.Second * 5)

	// 启动服务器
	err = server.Start()
	checkErrorExit(err, true)
	Log.Info("listener has been closed")

	// 等待优雅关闭
	err = server.WaitForGracefullyClose()
	checkErrorExit(err, false)
}

func checkErrorExit(err error, exit bool) {
	if nil != err {
		Log.Error(err)
		if exit {
			os.Exit(1)
		}
	}
}
