package discovery

import "github.com/wanghongfei/go-eureka-client/eureka"

var euClient *eureka.Client

func init() {
	c, err := eureka.NewClientFromFile("eureka.json")
	if nil != err {
		panic(err)
	}

	euClient = c
}

