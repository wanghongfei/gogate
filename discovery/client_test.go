package discovery

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestQueryAll(t *testing.T) {
	// QueryAll()
}

func TestQueryConsul(t *testing.T) {
	cfg := &api.Config{}
	cfg.Address = "127.0.0.1:8500"
	cfg.Scheme = "http"

	client, err := api.NewClient(cfg)

	checkList, _, err := client.Health().State("passing", &api.QueryOptions{})
	if nil != err {
		t.Error(err)
		return
	}

	for _, info := range checkList {
		fmt.Println(info.ServiceName)
	}
}
