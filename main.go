package main

import (
	"fmt"
	"os"

	"code.google.com/p/log4go"
	"github.com/wanghongfei/gogate/serv"
)

func main() {
	log4go.Info("start gogate at :8080")

	server, err := serv.NewGatewayServer("127.0.0.1", 8080, "route.yml", 10000)
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
