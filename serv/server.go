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

	fastServ		*fasthttp.Server

	// 保存每个instanceId对应的Http Client
	proxyClients	*sync.Map
}

const MAX_CONNECTION = 5000

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

	router, err := NewRouter(routePath)
	if nil != err {
		return nil, err
	}

	serv := &Server{
		host: host,
		port: port,

		router: router,
		proxyClients: new(sync.Map),
	}

	fastServ := &fasthttp.Server{
		Concurrency: maxConn,
		Handler: serv.HandleRequest,
	}

	serv.fastServ = fastServ

	return serv, nil

}

func (s *Server) Start() error {
	return s.fastServ.ListenAndServe(s.host + ":" + strconv.Itoa(s.port))
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




