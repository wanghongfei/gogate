package main

import (
	"os"
	"time"

	log "github.com/alecthomas/log4go"
	"github.com/wanghongfei/gogate/conf"
	serv "github.com/wanghongfei/gogate/server"
)

func main() {
	// 初始化
	serv.InitGogate("gogate.yml", "log.xml")

	// 构造gogate对象
	server, err := serv.NewGatewayServer(
		conf.App.ServerConfig.Host,
		conf.App.ServerConfig.Port,
		conf.App.EurekaConfig.RouteFile,
		conf.App.ServerConfig.MaxConnection,
		// 是否启用优雅关闭
		true,
		// 优雅关闭最大等待时间, 上一个参数为true时有效
		time.Second*1,
	)
	checkErrorExit(err, true)

	log.Info("pre filters: %v", server.ExportAllPreFilters())
	log.Info("post filters: %v", server.ExportAllPostFilters())

	// deferClose(server, time.Second * 5)

	// 启动服务器
	err = server.Start()
	checkErrorExit(err, true)
	log.Info("listener has been closed")

	// 等待优雅关闭
	err = server.WaitForGracefullyClose()
	checkErrorExit(err, false)
	time.Sleep(500 * time.Millisecond)

}

func checkErrorExit(err error, exit bool) {
	if nil != err {
		log.Error(err)
		time.Sleep(time.Second)
		if exit {
			os.Exit(1)
		}
	}
}
