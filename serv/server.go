package serv

import (
	"errors"
	"strconv"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/serv/filter"
)

type Server struct {
	host			string
	port			int

	router 			*Router

	preFilters		[]filter.PreFilterFunc
	postFilters		[]filter.PostFilterFunc

	// 保存每个instanceId对应的Http Client
	proxyClients	*sync.Map
}

/*
* 创建网关服务对象
*
* PARAMS:
*	- host: 主机名(ip)
*	- port: 端口
*	- routePath: 路由配置文件路径
*
*/
func NewGatewayServer(host string, port int, routePath string) (*Server, error) {
	if "" == host {
		return nil, errors.New("invalid host")
	}

	if port <= 0 || port > 65535 {
		return nil, errors.New("invalid port")
	}

	router, err := NewRouter(routePath)
	if nil != err {
		return nil, err
	}

	return &Server{
		host: host,
		port: port,

		router: router,
		proxyClients: new(sync.Map),
	}, nil
}

func (s *Server) Start() error {
	return fasthttp.ListenAndServe(s.host + ":" + strconv.Itoa(s.port), s.HandleRequest)
}

func (s *Server) Shutdown() {
	// todo gracefully shutdown
}

func (s *Server) ReloadRoute() error {
	return s.router.ReloadRoute()
}

func (s *Server) ExtractRoute() string {
	return s.router.ExtractRoute()
}

func (s *Server) RegisterPreFilter(preFunc filter.PreFilterFunc) {
	s.preFilters = append(s.preFilters, preFunc)
}

func (s *Server) RegisterPostFilter(postFunc filter.PostFilterFunc) {
	s.postFilters = append(s.postFilters, postFunc)
}




