package server

import (
	"github.com/alecthomas/log4go"
	"github.com/wanghongfei/gogate/asynclog"
	"github.com/wanghongfei/gogate/conf"
)

func InitGogate(gogateConfigFile, logConfigFile string) {
	// conf.InitLog(logConfigFile)
	asynclog.InitAsyncLog(logConfigFile, 2000)

	log4go.Info("initializing gogate config file")
	conf.LoadConfig(gogateConfigFile)
	log4go.Info("done initializing gogate config file")

}
