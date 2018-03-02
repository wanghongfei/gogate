package serv

import (
	"errors"
	"strconv"

	"github.com/valyala/fasthttp"
)

type Server struct {
	host		string
	port		int

	router 		*Router
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




