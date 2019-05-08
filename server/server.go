package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	log "github.com/alecthomas/log4go"
	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/discovery"
	"github.com/wanghongfei/gogate/redis"
	"github.com/wanghongfei/gogate/server/statistics"
	"github.com/wanghongfei/gogate/throttle"
)

type Server struct {
	host 					string
	port 					int

	// 负载均衡组件
	lb						LoadBalancer

	// 保存listener引用, 用于关闭server
	listen 					net.Listener
	// 是否启用优雅关闭
	graceShutdown 			bool
	// 优雅关闭最大等待时间
	maxWait 				time.Duration
	wg      				*sync.WaitGroup

	// URI路由组件
	Router 					*Router

	// 过滤器
	preFilters  			[]*PreFilter
	postFilters 			[]*PostFilter

	// fasthttp对象
	fastServ 				*fasthttp.Server

	isStarted 				bool

	// 保存服务地址
	// key: 服务名:版本号, 版本号为eureka注册信息中的metadata[version]值
	// val: []*InstanceInfo
	registryMap 			*InsInfoArrSyncMap

	// 服务id(string) -> 此服务的限速器对象(*MemoryRateLimiter)
	rateLimiterMap 			*RateLimiterSyncMap

	trafficStat 			*stat.TraficStat
}

const (
	// 默认最大连接数
	MAX_CONNECTION = 5000
	// 注册信息更新间隔, 秒
	REGISTRY_REFRESH_INTERVAL = 30
)

/*
* 创建网关服务对象
*
* PARAMS:
*	- host: 主机名(ip)
*	- port: 端口
*	- routePath: 路由配置文件路径
*	- maxConn: 最大连接数, 0表示使用默认值
*
 */
func NewGatewayServer(host string, port int, routePath string, maxConn int, useGracefullyShutdown bool, maxWait time.Duration) (*Server, error) {
	if "" == host {
		return nil, errors.New("invalid host")
	}

	if port <= 0 || port > 65535 {
		return nil, errors.New("invalid port")
	}

	if maxConn <= 0 {
		maxConn = MAX_CONNECTION
	}

	// 创建router
	router, err := NewRouter(routePath)
	if nil != err {
		return nil, err
	}

	// 创建Server对象
	serv := &Server{
		host: host,
		port: port,

		lb: &RoundRobinLoadBalancer{},

		Router:       router,

		preFilters:  make([]*PreFilter, 0, 3),
		postFilters: make([]*PostFilter, 0, 3),

		graceShutdown: useGracefullyShutdown,
		maxWait:       maxWait,
	}

	// 创建FastServer对象
	fastServ := &fasthttp.Server{
		Concurrency:  maxConn,
		Handler:      serv.HandleRequest,
		LogAllErrors: true,
	}

	serv.fastServ = fastServ

	// 创建每个服务的限速器
	serv.rebuildRateLimiter()

	// 注册过虑器
	serv.AppendPreFilter(NewPreFilter("service-match-pre-filter", ServiceMatchPreFilter))
	serv.InsertPreFilterBehind("service-match-pre-filter", NewPreFilter("rate-limit-pre-filter", RateLimitPreFilter))
	serv.InsertPreFilterBehind("rate-limit-pre-filter", NewPreFilter("url-rewrite-pre-filter", UrlRewritePreFilter))

	return serv, nil

}

// 启动服务器
func (serv *Server) Start() error {
	if conf.App.Traffic.EnableTrafficRecord {
		serv.trafficStat = stat.NewTrafficStat(1000, 1, stat.NewCsvFileTraficInfoStore(conf.App.Traffic.TrafficLogDir))
		serv.trafficStat.StartRecordTrafic()
	}

	serv.isStarted = true

	// 监听端口
	listen, err := net.Listen("tcp", serv.host+":"+strconv.Itoa(serv.port))
	if nil != err {
		return nil
	}

	// 是否启用优雅关闭功能
	if serv.graceShutdown {
		serv.wg = new(sync.WaitGroup)
	}

	// 保存Listener指针
	serv.listen = listen

	bothEnabled := conf.App.EurekaConfig.Enable && conf.App.ConsulConfig.Enable
	if bothEnabled {
		return errors.New("eureka and consul are both enabled")
	}

	// 注册
	go func() {
		time.Sleep(500 * time.Millisecond)

		if conf.App.EurekaConfig.Enable {
			log.Info("enable eureka")
			// 初始化eureka
			discovery.InitEurekaClient()
			// 注册
			discovery.StartRegister()

		} else if conf.App.ConsulConfig.Enable {
			log.Info("enable consul")
			// 初始化consul
			discovery.InitConsulClient()

		} else {
			panic("no registry center specified")
		}

		// 更新本地注册表
		serv.startRefreshRegistryInfo()

	}()

	// 启动http server
	return serv.fastServ.Serve(listen)
}

// 关闭server
func (serv *Server) Shutdown() error {
	serv.isStarted = false
	discovery.UnRegister()

	return serv.listen.Close()
}

// 需要在Shutdown()之后调用, 此方法会一直block直到所有连接关闭或者超时
func (serv *Server) WaitForGracefullyClose() error {
	select {
	case <-serv.waitAllRoutineDone():
		return nil

	case <-time.After(serv.maxWait):
		return fmt.Errorf("force shutdown after %v", serv.maxWait)
	}

}

// 等待所有请求处理routine完成;
// 此方法返回无缓冲channel, 只有当所有routine结束时channel会关闭
func (serv *Server) waitAllRoutineDone() chan struct{} {
	flagChan := make(chan struct{})

	go func() {
		if nil != serv.wg {
			serv.wg.Wait()
		}

		close(flagChan)
	}()

	return flagChan
}

// 更新路由配置文件
func (serv *Server) ReloadRoute() error {
	log.Info("start reloading route info")
	err := serv.Router.ReloadRoute()
	serv.rebuildRateLimiter()
	log.Info("route info reloaded")

	return err
}

// 将全部路由信息以字符串形式返回
func (serv *Server) ExtractRoute() string {
	return serv.Router.ExtractRoute()
}

// 启动定时更新注册表的routine
func (serv *Server) startRefreshRegistryInfo() {
	log.Info("refresh registry every %d sec", REGISTRY_REFRESH_INTERVAL)

	isBootstrap := true
	go func() {
		ticker := time.NewTicker(REGISTRY_REFRESH_INTERVAL * time.Second)

		for {
			log.Info("refresh registry started")
			err := serv.refreshRegistry()
			if nil != err {
				// 如果是第一次查询失败, 退出程序
				if isBootstrap {
					log.Error("failed to connect to eureka, err = %v, exit", err)
					os.Exit(1)
				}

				log.Error(err)
			}
			log.Info("done refreshing registry")

			isBootstrap = false

			<-ticker.C
		}
	}()
}

func (serv *Server) recordTraffic(ctx *fasthttp.RequestCtx, success bool) {
	if nil != serv.trafficStat {
		servName := GetStringFromUserValue(ctx, SERVICE_NAME)

		log.Debug("log traffic for %s", servName)

		info := &stat.TraficInfo{
			ServiceId: servName,
		}
		if success {
			info.SuccessCount = 1
		} else {
			info.FailedCount = 1
		}

		serv.trafficStat.RecordTrafic(info)
	}

}

// 给路由表中的每个服务重新创建限速器;
// 在更新过route.yml配置文件时调用
func (serv *Server) rebuildRateLimiter() {
	serv.rateLimiterMap = NewRateLimiterSyncMap()

	// 创建每个服务的限速器
	for _, info := range serv.Router.ServInfos {
		if 0 == info.Qps {
			continue
		}

		rl := serv.createRateLimiter(info)
		if nil != rl {
			serv.rateLimiterMap.Put(info.Id, rl)
			log.Debug("done building rateLimiter for %s", info.Id)
		}
	}
}

// 创建限速器对象
// 如果配置文件中设置了使用redis, 则创建RedisRateLimiter, 否则创建MemoryRateLimiter
func (serv *Server) createRateLimiter(info *ServiceInfo) throttle.RateLimiter {
	enableRedis := conf.App.RedisConfig.Enabled
	if !enableRedis {
		return throttle.NewMemoryRateLimiter(info.Qps)
	}

	client := redis.NewRedisClient(conf.App.RedisConfig.Addr, 5)
	err := client.Connect()
	if nil != err {
		log.Warn("failed to create ratelimiter, err = %v", err)
		return nil
	}

	rl, err := throttle.NewRedisRateLimiter(client, conf.App.RedisConfig.RateLimiterLua, info.Qps, info.Id)
	if nil != err {
		log.Warn("failed to create ratelimiter, err = %v", err)
		return nil
	}

	return rl
}
