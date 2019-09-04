package discovery

import (
	"fmt"
	log "github.com/alecthomas/log4go"
	"github.com/hashicorp/consul/api"
	"strconv"
	"strings"
)

func QueryEureka() ([]*InstanceInfo, error) {
	apps, err := euClient.GetApplications()
	if nil != err {
		return nil, fmt.Errorf("%w => failed to query eureka", err)
	}

	var instances []*InstanceInfo
	for _, app := range apps.Applications {
		// 服务名
		servName := app.Name

		// 遍历每一个实例
		for _, ins := range app.Instances {
			// 跳过无效实例
			if nil == ins.Port || ins.Status != "UP" {
				continue
			}

			addr := ins.HostName + ":" + strconv.Itoa(ins.Port.Port)
			var meta map[string]string
			if nil != ins.Metadata {
				meta = ins.Metadata.Map
			}

			instances = append(
				instances,
				&InstanceInfo{
					ServiceName: servName,
					Addr: addr,
					Meta: meta,
				},
			)
		}
	}

	return instances, nil
}

func QueryConsul() ([]*InstanceInfo, error) {
	// 查出所有实例
	servMap, err := consulClient.Agent().Services()
	if nil != err {
		return nil, err
	}

	// 查出所有健康实例
	healthList, _, err := consulClient.Health().State("passing", &api.QueryOptions{})
	if nil != err {
		return nil, fmt.Errorf("failed to query consul => %w", err)
	}

	instances := make([]*InstanceInfo, 0, 10)
	for _, servInfo := range servMap {
		servName := servInfo.Service
		servId := servInfo.ID

		// 查查在healthList中有没有
		isHealth := false
		for _, healthInfo := range healthList {
			if healthInfo.ServiceName == servName && healthInfo.ServiceID == servId {
				isHealth = true
				break
			}
		}

		if !isHealth {
			log.Warn("following instance is not health, skip; service name: %v, service id: %v", servName, servId)
			continue
		}

		instances = append(
			instances,
			&InstanceInfo{
				ServiceName: strings.ToUpper(servInfo.Service),
				Addr: servInfo.Address + ":" + strconv.Itoa(servInfo.Port),
				Meta: servInfo.Meta,
			},
		)
	}

	return instances, nil

}

