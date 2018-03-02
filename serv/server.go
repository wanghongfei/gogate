package serv

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/utils"
)

type Server struct {
	host		string
	port		int

	routePath	string
	routeMap	*sync.Map
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

	routeMap, err := LoadRoute(routePath)
	if nil != err {
		return nil, err
	}

	return &Server{
		host: host,
		port: port,

		routeMap: routeMap,
		routePath: routePath,
	}, nil
}

func (s *Server) Start() error {
	return fasthttp.ListenAndServe(s.host + ":" + strconv.Itoa(s.port), HandleRequest)
}

func (s *Server) Shutdown() {
	// todo gracefully shutdown
}

func (s *Server) ReloadRoute() error {
	newRoute, err := LoadRoute(s.routePath)
	if nil != err {
		return err
	}

	s.refreshRoute(newRoute)

	return nil
}

func (s *Server) ExtractRoute() string {
	var strBuf bytes.Buffer
	s.routeMap.Range(func(key, value interface{}) bool {
		strKey := key.(string)
		info := value.(*ServiceInfo)

		str := fmt.Sprintf("%s -> id:%s, path:%s\n", strKey, info.Id, info.Path)
		strBuf.WriteString(str)

		return true
	})

	return strBuf.String()
}

func (s *Server) refreshRoute(newRoute *sync.Map) {
	exclusiveKeys := utils.FindExclusiveKey(s.routeMap, newRoute)
	utils.DelKeys(s.routeMap, exclusiveKeys)
	utils.MergeSyncMap(newRoute, s.routeMap)
}




