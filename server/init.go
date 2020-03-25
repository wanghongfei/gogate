package server

import (
	. "github.com/wanghongfei/gogate/conf"
)

func InitGogate(gogateConfigFile string) {
	LoadConfig(gogateConfigFile)
	InitLog()
}
