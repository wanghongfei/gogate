package serv

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

	// path -> serviceId
	routeMap	*sync.Map
}

type ServiceInfo struct {
	Id		string
	Prefix	string
	Host	string
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
	routeMap, err := loadRoute(path)
	if nil != err {
		return nil, err
	}

	return &Router{
		routeMap: routeMap,
		cfgPath: path,
	}, nil
}

/*
* 重新加载路由器
*/
func (r *Router) ReloadRoute() error {
	newRoute, err := loadRoute(r.cfgPath)
	if nil != err {
		return err
	}

	r.refreshRoute(newRoute)

	return nil
}

/*
* 将路由信息转换成string返回
*/
func (r *Router) ExtractRoute() string {
	var strBuf bytes.Buffer
	r.routeMap.Range(func(key, value interface{}) bool {
		strKey := key.(string)
		info := value.(*ServiceInfo)

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

		appId, exist := r.routeMap.Load(matchTerm)
		if exist {
			return appId.(*ServiceInfo)
		}
	}

	return nil
}

func (r *Router) refreshRoute(newRoute *sync.Map) {
	exclusiveKeys, _ := utils.FindExclusiveKey(r.routeMap, newRoute)
	utils.DelKeys(r.routeMap, exclusiveKeys)
	utils.MergeSyncMap(newRoute, r.routeMap)
}

func loadRoute(path string) (*sync.Map, error) {
	// 打开配置文件
	routeFile, err := os.Open(path)
	if nil != err {
		return nil, err
	}
	defer routeFile.Close()

	// 读取
	buf, err := ioutil.ReadAll(routeFile)
	if nil != err {
		return nil, err
	}

	// 解析yml
	ymlMap := make(map[string]*ServiceInfo)
	err = yaml.UnmarshalStrict(buf, &ymlMap)
	if nil != err {
		return nil, err
	}


	// 构造 path->serviceId 映射
	var routeMap sync.Map
	for name, info := range ymlMap {
		// 验证
		err = validateServiceInfo(info)
		if nil != err {
			return nil, errors.New("invalid config for " + name + ":" + err.Error())
		}

		routeMap.Store(info.Prefix, info)
	}

	return &routeMap, nil
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