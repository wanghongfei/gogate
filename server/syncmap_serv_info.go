package server


import "sync"

// 封装sync.map, 提供类型安全的方法调用
type ServInfoSyncMap struct {
	dataMap			*sync.Map
}

func NewServInfoSyncMap() *ServInfoSyncMap {
	return &ServInfoSyncMap{
		dataMap: new(sync.Map),
	}
}

func (ssm *ServInfoSyncMap) Get(key string) (*ServiceInfo, bool) {
	val, exist := ssm.dataMap.Load(key)
	if !exist {
		return nil, false
	}

	info, ok := val.(*ServiceInfo)
	if !ok {
		return nil, false
	}

	return info, true
}

func (ssm *ServInfoSyncMap) Put(key string, val *ServiceInfo) {
	ssm.dataMap.Store(key, val)
}

func (ssm *ServInfoSyncMap) Each(eachFunc func(key string, val *ServiceInfo) bool) {
	ssm.dataMap.Range(func(key, value interface{}) bool {
		return eachFunc(key.(string), value.(*ServiceInfo))
	})
}

func (ssm *ServInfoSyncMap) GetMap() *sync.Map {
	return ssm.dataMap
}
