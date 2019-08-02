package syncmap

import "sync"

// 封装sync.map, 提供类型安全的方法调用
type StrStrSyncMap struct {
	dataMap			*sync.Map
}

func NewStrStrSyncMap() *StrStrSyncMap {
	return &StrStrSyncMap{
		dataMap: new(sync.Map),
	}
}

func (ssm *StrStrSyncMap) Get(key string) (string, bool) {
	val, exist := ssm.dataMap.Load(key)
	if !exist {
		return "", false
	}

	info, ok := val.(string)
	if !ok {
		return "", false
	}

	return info, true
}

func (ssm *StrStrSyncMap) Put(key string, val string) {
	ssm.dataMap.Store(key, val)
}

func (ssm *StrStrSyncMap) Each(eachFunc func(key string, val string) bool) {
	ssm.dataMap.Range(func(key, value interface{}) bool {
		return eachFunc(key.(string), value.(string))
	})
}

func (ssm *StrStrSyncMap) GetMap() *sync.Map {
	return ssm.dataMap
}
