package server

import (
	"github.com/alecthomas/log4go"
	"github.com/wanghongfei/gogate/conf"
)

func InitGogate(gogateConfigFile, logConfigFile string) {
	conf.InitLog(logConfigFile)

	log4go.Info("initializing gogate config file")
	conf.LoadConfig(gogateConfigFile)
	log4go.Info("done initializing gogate config file")

}
