package server

import (
	"github.com/wanghongfei/gogate/discovery"
	"sync"

)

// 封装sync.map, 提供类型安全的方法调用
type InsInfoArrSyncMap struct {
	dataMap			*sync.Map
}

func NewInsInfoArrSyncMap() *InsInfoArrSyncMap {
	return &InsInfoArrSyncMap{
		dataMap: new(sync.Map),
	}
}

func (ism *InsInfoArrSyncMap) Get(key string) ([]*discovery.InstanceInfo, bool) {
	val, exist := ism.dataMap.Load(key)
	if !exist {
		return nil, false
	}

	info, ok := val.([]*discovery.InstanceInfo)
	if !ok {
		return nil, false
	}

	return info, true
}

func (ism *InsInfoArrSyncMap) Put(key string, val []*discovery.InstanceInfo) {
	ism.dataMap.Store(key, val)
}

func (ism *InsInfoArrSyncMap) Each(eachFunc func(key string, val []*discovery.InstanceInfo) bool) {
	ism.dataMap.Range(func(key, value interface{}) bool {
		return eachFunc(key.(string), value.([]*discovery.InstanceInfo))
	})
}

func (ism *InsInfoArrSyncMap) GetMap() *sync.Map {
	return ism.dataMap
}
