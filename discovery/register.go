package discovery

import (
	"os"
	"time"

	"github.com/alecthomas/log4go"
	"github.com/wanghongfei/go-eureka-client/eureka"
	"github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/utils"
)

var euClient *eureka.Client
var gogateApp *eureka.InstanceInfo

func init() {
	c, err := eureka.NewClientFromFile(conf.App.EurekaConfig)
	if nil != err {
		panic(err)
	}

	euClient = c
}

func StartRegister() {
	ip, err := utils.GetFirstNoneLoopIp()
	if nil != err {
		panic(err)
	}

	host, err := os.Hostname()
	if nil != err {
		panic(err)
	}

	// 注册
	log4go.Info("register to eureka")
	gogateApp = eureka.NewInstanceInfo(host, conf.App.AppName, ip, conf.App.Port, 30, false)
	err = euClient.RegisterInstance("gogate", gogateApp)
	if nil != err {
		log4go.Warn("failed to register to eureka, %v", err)
	}

	// 心跳
	go func() {
		ticker := time.NewTicker(time.Second * 20)
		<- ticker.C

		heartbeat()
	}()
}

func heartbeat() {
	err := euClient.SendHeartbeat(gogateApp.App, gogateApp.HostName)
	if nil != err {
		log4go.Warn("failed to send heartbeat, %v", err)
		return
	}

	log4go.Info("heartbeat sent")
}

