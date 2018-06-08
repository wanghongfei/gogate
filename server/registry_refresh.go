package server

import (
	"strconv"
	"strings"
	"sync"
	"time"

	asynclog "github.com/alecthomas/log4go"
	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/go-eureka-client/eureka"
	"github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/discovery"
	"github.com/wanghongfei/gogate/utils"
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
	serv.refreshClients()

	return nil
}

// 刷新HttpClient.
// 遍历注册列表, 通过serviceId查找对应的Client, 如果没有则创建;
// 如果Client存在, 则对比host(每个实例)是否与注册表中的相同，如果不相同则抛弃老Client, 重新创建
func (serv *Server) refreshClients() error {
	if nil == serv.proxyClients {
		serv.proxyClients = NewInsMetaLbClientSyncMap()
	}

	changeCount := 0
	newCount := 0

	// 遍历注册列表
	serv.registryMap.Each(func(name string, infos []*InstanceInfo) bool {
		// 服务名
		name = strings.ToLower(name)

		// 按版本号分组
		groupMap := groupByVersion(name, infos)

		// 遍历分组后的服务表
		for fullname, hosts := range groupMap {
			// 根据serviceId查对应的Http Client
			client, exist := serv.proxyClients.Get(fullname)

			// 如果注册表中的service不存在Client
			// 则为此服务创建Client
			if !exist {
				asynclog.Debug("create new client for service: %s", name)

				// 此service不存在, 创建新的
				newClient := &fasthttp.LBClient{
					Clients: createClients(hosts),
					Timeout: time.Millisecond * time.Duration(conf.App.ServerConfig.Timeout),
				}

				serv.proxyClients.Put(fullname, newClient)
				newCount++

			} else {
				// service对应的Client
				// 对比所有实例(host)是否有变化
				changed := isHostsChanged(client, hosts)
				if changed {
					// 发生了变化
					// 创建新的LBClient替换掉老的
					asynclog.Debug("service %s changed", name)
					newClient := &fasthttp.LBClient{
						Clients: createClients(hosts),
						Timeout: time.Millisecond * time.Duration(conf.App.ServerConfig.Timeout),
					}

					serv.proxyClients.Put(fullname, newClient)
					changeCount++
				}
			}

		}

		return true
	})


	asynclog.Info("%d services updated, %d services created", changeCount, newCount)
	return nil
}

// 将服务信息分组, 组的key为 serviceId:version, 如果没有version则key为serviceId
func groupByVersion(serviceName string, infos []*InstanceInfo) map[string][]string {
	groupMap := make(map[string][]string)

	for _, info := range infos {
		var key = serviceName
		if info.Meta != nil {
			if ver, exist := info.Meta[META_VERSION]; exist {
				key = serviceName + ":" + ver
			}
		}

		hosts := groupMap[key]
		hosts = append(hosts, info.Addr)
		groupMap[key] = hosts
	}

	return groupMap
}

// 对比LBClient中的host与注册列表中的host有没有变化
// 返回true表示有变化
func isHostsChanged(lbClient *fasthttp.LBClient, newHosts []string) bool {
	if len(lbClient.Clients) != len(newHosts) {
		return true
	}

	// 遍历LBClient里的每一个Client对象
	for _, client := range lbClient.Clients {
		c := client.(*fasthttp.HostClient)

		// 判断此Client的地址在不在newHosts中
		match := false
		for _, h := range newHosts {
			if h == c.Addr {
				match = true
				break
			}
		}

		// 有一个不存在的, 就认为发生了变化
		if !match {
			return true
		}
	}

	return false
}

// 为每一个host创建一个HostClient
func createClients(hosts[] string) []fasthttp.BalancingClient {
	var cs []fasthttp.BalancingClient
	for _, host := range hosts {
		client := &fasthttp.HostClient{
			Addr: host,
		}

		cs = append(cs, client)
	}

	return cs
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
