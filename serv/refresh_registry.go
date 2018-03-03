package serv

import (
	"log"
	"strconv"
	"sync"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/wanghongfei/gogate/discovery"
	"github.com/wanghongfei/gogate/utils"
)

// 向eureka查询注册列表, 刷新本地列表
func refreshRegistry(serv *Server) error {
	apps, err := discovery.QueryAll()
	if nil != err {
		return err
	}

	if nil == apps {
		log.Println("no service found")
		return nil
	}

	newRegistryMap := convertToMap(apps)
	refreshRegistryMap(serv, newRegistryMap)

	return nil
}

// 将新服务列表保存为map
func convertToMap(apps []eureka.Application) *sync.Map {
	newAppsMap := new(sync.Map)
	for _, app := range apps {
		// 服务名
		servName := app.Name

		// 遍历每一个实例
		var instances []string
		for _, ins := range app.Instances {
			instances = append(instances, ins.HostName + ":" + strconv.Itoa(ins.Port.Port))
		}

		newAppsMap.Store(servName, instances)
	}

	return newAppsMap
}

func refreshRegistryMap(s *Server, newRegistry *sync.Map) {
	exclusiveKeys := utils.FindExclusiveKey(s.registryMap, newRegistry)
	utils.DelKeys(s.registryMap, exclusiveKeys)
	utils.MergeSyncMap(newRegistry, s.registryMap)
}
