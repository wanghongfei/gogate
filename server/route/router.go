package route

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Router struct {
	// 配置文件路径
	cfgPath			string

	// path(string) -> *ServiceInfo
	pathMatcher		*PathMatcher

	ServInfos		[]*ServiceInfo
}

type ServiceInfo struct {
	Id				string
	Prefix			string
	Host			string
	Name			string
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
	matcher, servInfos, err := loadRoute(path)
	if nil != err {
		return nil, err
	}


	return &Router{
		pathMatcher: matcher,
		cfgPath: path,
		ServInfos: servInfos,
	}, nil
}

/*
* 重新加载路由器
*/
func (r *Router) ReloadRoute() error {
	matcher, servInfos, err := loadRoute(r.cfgPath)
	if nil != err {
		return err
	}

	r.ServInfos = servInfos
	r.pathMatcher = matcher

	return nil
}

/*
* 根据uri选择一个最匹配的appId
*
* RETURNS:
*	返回最匹配的ServiceInfo
*/
func (r *Router) Match(reqPath string) *ServiceInfo {

	return r.pathMatcher.Match(reqPath)
}

func loadRoute(path string) (*PathMatcher, []*ServiceInfo, error) {
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
	// 保存到字典树中
	tree := NewTrieTree()
	// 保存到map中
	routeMap := make(map[string]*ServiceInfo)
	for name, info := range ymlMap["services"] {
		// 验证
		err = validateServiceInfo(info)
		if nil != err {
			return nil, nil, errors.New("invalid config for " + name + ":" + err.Error())
		}

		tree.PutString(info.Prefix, info)
		routeMap[info.Prefix] = info

		servInfos = append(servInfos, info)
	}


	matcher := &PathMatcher{
		routeMap: routeMap,
		routeTrieTree: tree,
	}
	return matcher, servInfos, nil
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