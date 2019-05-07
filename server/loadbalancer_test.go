package server

import (
	"fmt"
	"testing"
)

func TestRoundRobinLoadBalancer_Choose(t *testing.T) {
	lb := &RoundRobinLoadBalancer{}

	instances := make([]*InstanceInfo, 0)
	instances = append(instances, &InstanceInfo{
		Addr: "1",
	})
	instances = append(instances, &InstanceInfo{
		Addr: "2",
	})
	instances = append(instances, &InstanceInfo{
		Addr: "3",
	})

	for ix := 0; ix < 10; ix++ {
		target := lb.Choose(instances)
		fmt.Println(target.Addr)
	}
}
