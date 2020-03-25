package discovery

import (
	"github.com/hashicorp/consul/api"
	"github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/utils"
)

var consulClient *api.Client

func InitConsulClient() error {
	cfg := &api.Config{}
	cfg.Address = conf.App.ConsulConfig.Address
	cfg.Scheme = "http"

	c, err := api.NewClient(cfg)
	if nil != err {
		return utils.Errorf("failed to init consule client => %w", err)
	}

	consulClient = c

	return nil
}
