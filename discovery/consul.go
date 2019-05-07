package discovery

import (
	"github.com/hashicorp/consul/api"
	"github.com/wanghongfei/gogate/conf"
)

var consulClient *api.Client

func InitConsulClient()  {
	cfg := &api.Config{}
	cfg.Address = conf.App.ConsulConfig.Address
	cfg.Scheme = "http"

	c, err := api.NewClient(cfg)
	if nil != err {
		panic(err)
	}

	consulClient = c
}
