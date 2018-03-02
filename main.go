package main

import (
	"fmt"
	"os"

	"github.com/wanghongfei/gogate/serv"
)

func main() {
	fmt.Println("start gogate at 8080")

	server, err := serv.NewGatewayServer("127.0.0.1", 8080, "route.yml")
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
