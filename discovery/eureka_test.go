package discovery

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestStartRegister(t *testing.T) {
	// StartRegister()
	// time.Sleep(time.Second * 60)
}

func TestRegisterToConsul(t *testing.T) {
	client, err := api.NewClient(api.DefaultConfig())
	if nil != err {
		t.Error(err)
		return
	}

	reg := &api.AgentServiceRegistration{}
	reg.ID = "id"
	reg.Name = "go-unit-test"
	reg.Address = "127.0.0.1"
	reg.Port = 8080
	reg.Check = &api.AgentServiceCheck{}
	reg.Check.HTTP = "http://127.0.0.1:9000"
	reg.Check.Method = "GET"
	reg.Check.Interval = "10s"
	reg.Check.Timeout = "1s"
	// reg.Check.DeregisterCriticalServiceAfter = "2s"

	err = client.Agent().ServiceRegister(reg)
	if nil != err {
		t.Error(err)
		return
	}


	servMap, err := client.Agent().Services()
	if nil != err {
		t.Error(err)
		return
	}

	for name, serv := range servMap {
		fmt.Println(name)
		fmt.Println(serv)
	}
}