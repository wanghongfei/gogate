package serv

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"code.google.com/p/log4go"
	"github.com/valyala/fasthttp"
)

type Server struct {
	host			string
	port			int

	// URI路由组件
	Router 			*Router

	preFilters		[]PreFilterFunc
	postFilters		[]PostFilterFunc

	// fasthttp对象
	fastServ		*fasthttp.Server

	// 保存每个instanceId对应的Http Client
	// key: instanceId
	// val: *LBClient
	proxyClients	*sync.Map

	// 保存服务地址
	// key: 服务名
	// val: host:port数组, []string类型
	registryMap		*sync.Map
}

const (
	// 默认最大连接数
	MAX_CONNECTION = 5000
	// 注册信息更新间隔, 秒
	REGISTRY_REFRESH_INTERVAL = 20
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
func NewGatewayServer(host string, port int, routePath string, maxConn int) (*Server, error) {
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

		Router: router,
		proxyClients: new(sync.Map),
	}

	// 创建FastServer对象
	fastServ := &fasthttp.Server{
		Concurrency: maxConn,
		Handler: serv.HandleRequest,
		LogAllErrors: true,
	}

	serv.fastServ = fastServ

	// 注册过虑器
	serv.RegisterPreFilter(ServiceMatchPreFilter)

	return serv, nil

}

// 启动服务器
func (s *Server) Start() error {
	s.startRefreshRegistryInfo()

	return s.fastServ.ListenAndServe(s.host + ":" + strconv.Itoa(s.port))
}

// 优雅关闭
func (s *Server) Shutdown() {
	// todo gracefully shutdown
}

// 更新路由配置文件
func (s *Server) ReloadRoute() error {
	log4go.Info("start reloading route info")
	err := s.Router.ReloadRoute()
	log4go.Info("route info reloaded")

	return err
}

// 将全部路由信息以字符串形式返回
func (s *Server) ExtractRoute() string {
	return s.Router.ExtractRoute()
}

// 注册过滤器, 追加到末尾
func (s *Server) RegisterPreFilter(preFunc PreFilterFunc) {
	s.preFilters = append(s.preFilters, preFunc)
}

// 注册过滤器, 追加到末尾
func (s *Server) RegisterPostFilter(postFunc PostFilterFunc) {
	s.postFilters = append(s.postFilters, postFunc)
}

func (s *Server) startRefreshRegistryInfo() {
	log4go.Info("refresh registry every %d sec", REGISTRY_REFRESH_INTERVAL)

	go func() {
		ticker := time.NewTicker(REGISTRY_REFRESH_INTERVAL * time.Second)

		for {
			log4go.Info("refresh registry started")
			err := refreshRegistry(s)
			if nil != err {
				log4go.Error(err)
			}
			log4go.Info("done refreshing registry")

			<- ticker.C
		}
	}()
}



