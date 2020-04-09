package discovery

import (
	. "github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/perr"
	"github.com/wanghongfei/gogate/utils"
	"sync"
	"time"
)
const REGISTRY_REFRESH_INTERVAL = 30

// 向eureka查询注册列表, 刷新本地列表
func startPeriodicalRefresh(c Client) error {
	Log.Infof("refresh registry every %d sec", REGISTRY_REFRESH_INTERVAL)

	refreshRegistryChan := make(chan error)

	isBootstrap := true
	go func() {
		ticker := time.NewTicker(REGISTRY_REFRESH_INTERVAL * time.Second)

		for {
			Log.Info("registry refresh started")
			err := doRefresh(c)
			if nil != err {
				// 如果是第一次查询失败, 退出程序
				if isBootstrap {
					refreshRegistryChan <- perr.SystemErrorf("failed to refresh registry => %w", err)
					return

				} else {
					Log.Error(err)
				}

			}
			Log.Info("done refreshing registry")

			if isBootstrap {
				isBootstrap = false
				close(refreshRegistryChan)
			}

			<-ticker.C
		}
	}()

	return <- refreshRegistryChan
}

func doRefresh(c Client) error {
	instances, err := c.QueryServices()

	if nil != err {
		return perr.SystemErrorf("failed to refresh registry => %w", err)
	}

	if nil == instances {
		Log.Info("no instance found")
		return nil
	}

	Log.Infof("total app count: %d", len(instances))

	newRegistryMap := groupByService(instances)

	refreshRegistryMap(c, newRegistryMap)

	return nil

}


// 将所有实例按服务名进行分组
func groupByService(instances []*InstanceInfo) *sync.Map {
	servMap := new(sync.Map)
	for _, ins := range instances {
		infosGeneric, exist := servMap.Load(ins.ServiceName)
		if !exist {
			infosGeneric = make([]*InstanceInfo, 0, 5)
			infosGeneric = append(infosGeneric.([]*InstanceInfo), ins)

		} else {
			infosGeneric = append(infosGeneric.([]*InstanceInfo), ins)
		}
		servMap.Store(ins.ServiceName, infosGeneric)
	}
	return servMap
}


// 更新本地注册列表
// s: gogate server对象
// newRegistry: 刚从eureka查出的最新服务列表
func refreshRegistryMap(c Client, newRegistry *sync.Map) {
	if nil == c.GetInternalRegistryStore() {
		c.SetInternalRegistryStore(NewInsInfoArrSyncMap())
	}

	// 找出本地列表存在, 但新列表中不存在的服务
	exclusiveKeys, _ := utils.FindExclusiveKey(c.GetInternalRegistryStore().GetMap(), newRegistry)
	// 删除本地多余的服务
	utils.DelKeys(c.GetInternalRegistryStore().GetMap(), exclusiveKeys)
	// 将新列表中的服务合并到本地列表中
	utils.MergeSyncMap(newRegistry, c.GetInternalRegistryStore().GetMap())
}
