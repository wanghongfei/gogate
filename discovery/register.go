package discovery

import "github.com/ArthurHlt/go-eureka-client/eureka"

var euClient *eureka.Client

func init() {
	euClient = eureka.NewClient([]string{
		"http://10.150.186.11:8761/eureka",
	})
}

