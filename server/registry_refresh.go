package server

import (
	asynclog "github.com/alecthomas/log4go"
	"github.com/wanghongfei/go-eureka-client/eureka"
	"github.com/wanghongfei/gogate/discovery"
	"github.com/wanghongfei/gogate/utils"
	"strconv"
	"sync"
)

const META_VERSION = "version"

// 向eureka查询注册列表, 刷新本地列表
func (serv *Server) refreshRegistry() error {
	apps, err := discovery.QueryAll()
	if nil != err {
		return err
	}
	asynclog.Info("total app count: %d", len(apps))

	if nil == apps {
		asynclog.Error("no service found")
		return nil
	}

	newRegistryMap := convertToMap(apps)
	asynclog.Info("refreshing registry")

	refreshRegistryMap(serv, newRegistryMap)
	asynclog.Info("refreshing clients")

	return nil
}

func createInstanceInfos(hosts []string) []*InstanceInfo {
	hostNum := len(hosts)

	infos := make([]*InstanceInfo, 0, hostNum)
	for _, host := range hosts {
		infos = append(infos, &InstanceInfo{
			Addr: host,
		})
	}

	return infos
}

// 将新服务列表保存为map
func convertToMap(apps []eureka.Application) *sync.Map {
	newAppsMap := new(sync.Map)
	for _, app := range apps {
		// 服务名
		servName := app.Name

		// 遍历每一个实例
		var instances []*InstanceInfo
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
					Addr: addr,
					Meta: meta,
				},
			)
		}

		newAppsMap.Store(servName, instances)
	}

	return newAppsMap
}

// 更新本地注册列表
// s: gogate server对象
// newRegistry: 刚从eureka查出的最新服务列表
func refreshRegistryMap(s *Server, newRegistry *sync.Map) {
	if nil == s.registryMap {
		s.registryMap = NewInsInfoArrSyncMap()
	}

	// 找出本地列表存在, 但新列表中不存在的服务
	exclusiveKeys, _ := utils.FindExclusiveKey(s.registryMap.GetMap(), newRegistry)
	// 删除本地多余的服务
	utils.DelKeys(s.registryMap.GetMap(), exclusiveKeys)
	// 将新列表中的服务合并到本地列表中
	utils.MergeSyncMap(newRegistry, s.registryMap.GetMap())
}
