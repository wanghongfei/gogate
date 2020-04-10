package discovery

import (
	"github.com/wanghongfei/gogate/perr"
	"strconv"
	"time"

	"github.com/wanghongfei/go-eureka-client/eureka"
	"github.com/wanghongfei/gogate/conf"
	. "github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/utils"
)

// var euClient *eureka.Client
var gogateApp *eureka.InstanceInfo
var instanceId = ""

var ticker *time.Ticker
var tickerCloseChan chan struct{}

type EurekaClient struct {
	// 继承方法
	*periodicalRefreshClient

	client *eureka.Client

	// 保存服务地址
	// key: 服务名:版本号, 版本号为eureka注册信息中的metadata[version]值
	// val: []*InstanceInfo
	registryMap 			*InsInfoArrSyncMap
}

func NewEurekaClient(confFile string) (Client, error) {
	c, err := eureka.NewClientFromFile(confFile)
	if nil != err {
		return nil, perr.SystemErrorf("failed to init eureka client => %w", err)
	}

	euClient := &EurekaClient{client:c}
	euClient.periodicalRefreshClient = newPeriodicalRefresh(euClient)

	return euClient, nil
}

func (c *EurekaClient) Get(serviceId string) []*InstanceInfo {
	instance, exist := c.registryMap.Get(serviceId)
	if !exist {
		return nil
	}

	return instance
}

func (c *EurekaClient) GetInternalRegistryStore() *InsInfoArrSyncMap {
	return c.registryMap
}

func (c *EurekaClient) SetInternalRegistryStore(registry *InsInfoArrSyncMap) {
	c.registryMap = registry
}


// 查询所有服务
func (c *EurekaClient) QueryServices() ([]*InstanceInfo, error) {
	apps, err := c.client.GetApplications()
	if nil != err {
		return nil, perr.SystemErrorf("faield to query eureka => %w", err)
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

func (c *EurekaClient) Register() error {
	ip, err := utils.GetFirstNoneLoopIp()
	if nil != err {
		return perr.SystemErrorf("failed to get first none loop ip => %w", err)
	}


	instanceId = ip + ":" + strconv.Itoa(conf.App.ServerConfig.Port)

	// 注册
	Log.Infof("register to eureka as %s", instanceId)
	gogateApp = eureka.NewInstanceInfo(
		instanceId,
		conf.App.ServerConfig.AppName,
		ip,
		conf.App.ServerConfig.Port,
		conf.App.EurekaConfig.EvictionDuration,
		false,
	)
	gogateApp.Metadata = &eureka.MetaData{
		Class: "",
		Map: map[string]string {"version": conf.App.Version},
	}

	err = c.client.RegisterInstance(conf.App.ServerConfig.AppName, gogateApp)
	if nil != err {
		return perr.SystemErrorf("failed to register to eureka => %w", err)
	}

	// 心跳
	go func() {
		ticker = time.NewTicker(time.Second * time.Duration(conf.App.EurekaConfig.HeartbeatInterval))
		tickerCloseChan = make(chan struct{})

		for {
			select {
			case <-ticker.C:
				c.heartbeat()

			case <-tickerCloseChan:
				Log.Info("heartbeat stopped")
				return

			}
		}
	}()

	return nil
}

func (c *EurekaClient) UnRegister() error {
	c.stopHeartbeat()

	Log.Infof("unregistering %s", instanceId)
	err := c.client.UnregisterInstance("gogate", instanceId)

	if nil != err {
		return perr.SystemErrorf("failed to unregister => %w", err)
	}

	Log.Info("done unregistration")
	return nil
}

func (c *EurekaClient) stopHeartbeat() {
	ticker.Stop()
	close(tickerCloseChan)
}

func (c *EurekaClient) heartbeat() {
	err := c.client.SendHeartbeat(gogateApp.App, instanceId)
	if nil != err {
		Log.Warnf("failed to send heartbeat, %v", err)
		return
	}

	Log.Info("heartbeat sent")
}

