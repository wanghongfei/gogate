package server

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/wanghongfei/gogate/utils"
	"gopkg.in/yaml.v2"
)

type Router struct {
	// 配置文件路径
	cfgPath		string

	// path(string) -> *ServiceInfo
	routeMap	*ServInfoSyncMap

	ServInfos	[]*ServiceInfo
}

type ServiceInfo struct {
	Id				string
	Prefix			string
	Host			string
	StripPrefix		bool`yaml:"strip-prefix"`
	Qps				int

	Canary			[]*CanaryInfo
}

type CanaryInfo struct {
	Meta		string
	Weight		int
}

func (info *ServiceInfo) String() string {
	return "prefix = " + info.Prefix + ", id = " + info.Id + ", host = " + info.Host
}

/*
* 创建路由器
*
* PARAMS:
*	- path: 路由配置文件路径
*
*/
func NewRouter(path string) (*Router, error) {
	routeMap, servInfos, err := loadRoute(path)
	if nil != err {
		return nil, err
	}

	return &Router{
		routeMap: routeMap,
		cfgPath: path,
		ServInfos: servInfos,
	}, nil
}

/*
* 重新加载路由器
*/
func (r *Router) ReloadRoute() error {
	newRoute, servInfos, err := loadRoute(r.cfgPath)
	if nil != err {
		return err
	}

	r.ServInfos = servInfos
	r.refreshRoute(newRoute.GetMap())

	return nil
}

/*
* 将路由信息转换成string返回
*/
func (r *Router) ExtractRoute() string {
	var strBuf bytes.Buffer
	r.routeMap.Each(func(strKey string, info *ServiceInfo) bool {

		str := fmt.Sprintf("%s -> id:%s, path:%s\n", strKey, info.Id, info.Prefix)
		strBuf.WriteString(str)

		return true
	})

	return strBuf.String()
}

/*
* 根据uri选择一个最匹配的appId
*
* RETURNS:
*	返回最匹配的ServiceInfo
*/
func (r *Router) Match(reqPath string) *ServiceInfo {
	if !strings.HasSuffix(reqPath, "/") {
		reqPath = reqPath + "/"
	}

	if "/" == reqPath {
		reqPath = "//"
	}

	// 以/为分隔符, 从后向前匹配
	// 每次循环都去掉最后一个/XXXX节点
	term := reqPath
	for {
		lastSlash := strings.LastIndex(term, "/")
		if -1 == lastSlash {
			break
		}

		matchTerm := term[0:lastSlash]
		term = matchTerm

		if "" == matchTerm {
			matchTerm = "/"
		}

		appId, exist := r.routeMap.Get(matchTerm)
		if exist {
			return appId
		}
	}

	return nil
}

func (r *Router) refreshRoute(newRoute *sync.Map) {
	exclusiveKeys, _ := utils.FindExclusiveKey(r.routeMap.GetMap(), newRoute)
	utils.DelKeys(r.routeMap.GetMap(), exclusiveKeys)
	utils.MergeSyncMap(newRoute, r.routeMap.GetMap())
}

func loadRoute(path string) (*ServInfoSyncMap, []*ServiceInfo, error) {
	// 打开配置文件
	routeFile, err := os.Open(path)
	if nil != err {
		return nil, nil, err
	}
	defer routeFile.Close()

	// 读取
	buf, err := ioutil.ReadAll(routeFile)
	if nil != err {
		return nil, nil, err
	}

	// 解析yml
	// ymlMap := make(map[string]*ServiceInfo)
	ymlMap := make(map[string]map[string]*ServiceInfo)
	err = yaml.UnmarshalStrict(buf, &ymlMap)
	if nil != err {
		return nil, nil, err
	}

	servInfos := make([]*ServiceInfo, 0, 10)

	// 构造 path->serviceId 映射
	// var routeMap sync.Map
	routeMap := NewServInfoSyncMap()
	for name, info := range ymlMap["services"] {
		// 验证
		err = validateServiceInfo(info)
		if nil != err {
			return nil, nil, errors.New("invalid config for " + name + ":" + err.Error())
		}

		routeMap.Put(info.Prefix, info)
		servInfos = append(servInfos, info)
	}

	return routeMap, servInfos, nil
}

func validateServiceInfo(info *ServiceInfo) error {
	if nil == info {
		return errors.New("info is empty")
	}

	if "" == info.Id && "" == info.Host {
		return errors.New("id and host are both empty")
	}

	if "" == info.Prefix {
		return errors.New("path is empty")
	}

	return nil
}