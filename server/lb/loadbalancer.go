package lb

import (
	"github.com/wanghongfei/gogate/discovery"
	"sync/atomic"
)

// 负载均衡接口
type LoadBalancer interface {
	// 从instance中选一个对象返回
	Choose(instances []*discovery.InstanceInfo) *discovery.InstanceInfo

	ChooseByAddresses(addrs []string) string
}

// 轮询均衡器实现
type RoundRobinLoadBalancer struct {
	index	int64
}

func (lb *RoundRobinLoadBalancer) Choose(instances []*discovery.InstanceInfo) *discovery.InstanceInfo {
	total := len(instances)
	next := lb.nextIndex(total)

	return instances[next]
}

func (lb *RoundRobinLoadBalancer) ChooseByAddresses(addrs []string) string {
	total := len(addrs)
	next := lb.nextIndex(total)

	return addrs[next]
}

func (lb *RoundRobinLoadBalancer) nextIndex(total int) int64 {
	next := lb.index % int64(total)
	if next < 0 {
		next = next * -1
	}

	atomic.AddInt64(&lb.index, 1)

	return next
}

