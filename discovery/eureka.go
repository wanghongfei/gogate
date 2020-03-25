package discovery

import (
	"strconv"
	"time"

	"github.com/wanghongfei/go-eureka-client/eureka"
	"github.com/wanghongfei/gogate/conf"
	. "github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/utils"
)

var euClient *eureka.Client
var gogateApp *eureka.InstanceInfo
var instanceId = ""

var ticker *time.Ticker
var tickerCloseChan chan struct{}

func InitEurekaClient() error {
	c, err := eureka.NewClientFromFile(conf.App.EurekaConfig.ConfigFile)
	if nil != err {
		utils.Errorf("failed to init eureka client => %w", err)
	}

	euClient = c

	return nil
}

func StartRegister() {
	ip, err := utils.GetFirstNoneLoopIp()
	if nil != err {
		panic(err)
	}

	//host, err := os.Hostname()
	//if nil != err {
	//	panic(err)
	//}

	instanceId = ip + ":" + strconv.Itoa(conf.App.ServerConfig.Port)

	// 注册
	Log.Info("register to eureka as %s", instanceId)
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

	err = euClient.RegisterInstance(conf.App.ServerConfig.AppName, gogateApp)
	if nil != err {
		Log.Warn("failed to register to eureka, %v", err)
	}

	// 心跳
	go func() {
		ticker = time.NewTicker(time.Second * time.Duration(conf.App.EurekaConfig.HeartbeatInterval))
		tickerCloseChan = make(chan struct{})

		for {
			select {
			case <-ticker.C:
				heartbeat()

			case <-tickerCloseChan:
				Log.Info("heartbeat stopped")
				return

			}
		}
	}()
}

func UnRegister() {
	stopHeartbeat()

	Log.Info("unregistering %s", instanceId)
	err := euClient.UnregisterInstance("gogate", instanceId)

	if nil != err {
		Log.Error(err)
		return
	}

	Log.Info("done unregistration")
}

func stopHeartbeat() {
	ticker.Stop()
	close(tickerCloseChan)
}

func heartbeat() {
	err := euClient.SendHeartbeat(gogateApp.App, instanceId)
	if nil != err {
		Log.Warn("failed to send heartbeat, %v", err)
		return
	}

	Log.Info("heartbeat sent")
}

