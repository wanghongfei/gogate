package server

import (
	"fmt"
	"github.com/wanghongfei/gogate/discovery"
	"testing"
)

func TestRoundRobinLoadBalancer_Choose(t *testing.T) {
	lb := &RoundRobinLoadBalancer{}

	instances := make([]*discovery.InstanceInfo, 0)
	instances = append(instances, &discovery.InstanceInfo{
		Addr: "1",
	})
	instances = append(instances, &discovery.InstanceInfo{
		Addr: "2",
	})
	instances = append(instances, &discovery.InstanceInfo{
		Addr: "3",
	})

	for ix := 0; ix < 10; ix++ {
		target := lb.Choose(instances)
		fmt.Println(target.Addr)
	}
}
