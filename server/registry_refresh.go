package server

import (
	"github.com/wanghongfei/go-eureka-client/eureka"
	"github.com/wanghongfei/gogate/conf"
	. "github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/discovery"
	"github.com/wanghongfei/gogate/server/syncmap"
	"github.com/wanghongfei/gogate/utils"
	"strconv"
	"sync"
)

const META_VERSION = "version"

// 向eureka查询注册列表, 刷新本地列表
func (serv *Server) refreshRegistry() error {
	var instances []*discovery.InstanceInfo
	var err error

	if conf.App.EurekaConfig.Enable {
		instances, err = discovery.QueryEureka()

	} else if conf.App.ConsulConfig.Enable {
		instances, err = discovery.QueryConsul()
	}
	if nil != err {
		return utils.Errorf("failed to communicate with discovery service => %w", err)
	}

	Log.Infof("total app count: %d", len(instances))

	if nil == instances {
		Log.Info("no service found")
		return nil
	}

	newRegistryMap := groupByService(instances)
	// log.Info("refreshing registry")

	refreshRegistryMap(serv, newRegistryMap)
	// log.Info("refreshing clients")

	return nil
}

func createInstanceInfos(hosts []string) []*discovery.InstanceInfo {
	hostNum := len(hosts)

	infos := make([]*discovery.InstanceInfo, 0, hostNum)
	for _, host := range hosts {
		infos = append(infos, &discovery.InstanceInfo{
			Addr: host,
		})
	}

	return infos
}

// 将所有实例按服务名进行分组
func groupByService(instances []*discovery.InstanceInfo) *sync.Map {
	servMap := new(sync.Map)
	for _, ins := range instances {
		infosGeneric, exist := servMap.Load(ins.ServiceName)
		if !exist {
			infosGeneric = make([]*discovery.InstanceInfo, 0, 5)
			infosGeneric = append(infosGeneric.([]*discovery.InstanceInfo), ins)

		} else {
			infosGeneric = append(infosGeneric.([]*discovery.InstanceInfo), ins)
		}
		servMap.Store(ins.ServiceName, infosGeneric)
	}
	return servMap
}

// 将新服务列表保存为map
func convertToMap(apps []eureka.Application) *sync.Map {
	newAppsMap := new(sync.Map)
	for _, app := range apps {
		// 服务名
		servName := app.Name

		// 遍历每一个实例
		var instances []*discovery.InstanceInfo
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
				&discovery.InstanceInfo{
					ServiceName: servName,
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
		s.registryMap = syncmap.NewInsInfoArrSyncMap()
	}

	// 找出本地列表存在, 但新列表中不存在的服务
	exclusiveKeys, _ := utils.FindExclusiveKey(s.registryMap.GetMap(), newRegistry)
	// 删除本地多余的服务
	utils.DelKeys(s.registryMap.GetMap(), exclusiveKeys)
	// 将新列表中的服务合并到本地列表中
	utils.MergeSyncMap(newRegistry, s.registryMap.GetMap())
}
