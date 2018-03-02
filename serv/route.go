package serv

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type ServiceInfo struct {
	Id		string
	Path	string
}

func LoadRoute(path string) (*sync.Map, error) {
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

		routeMap.Store(info.Path, info)
	}

	return &routeMap, nil
}

func validateServiceInfo(info *ServiceInfo) error {
	if nil == info {
		return errors.New("info is empty")
	}

	if "" == info.Id {
		return errors.New("id is empty")
	}

	if "" == info.Path {
		return errors.New("path is empty")
	}

	return nil
}