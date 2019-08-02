package syncmap

import (
	"sync"

	"github.com/valyala/fasthttp"
)

// 封装sync.map, 提供类型安全的方法调用
type InsLbClientSyncMap struct {
	mcMap			*sync.Map
}

func NewInsMetaLbClientSyncMap() *InsLbClientSyncMap {
	return &InsLbClientSyncMap{
		mcMap: new(sync.Map),
	}
}

func (ism *InsLbClientSyncMap) Get(key string) (*fasthttp.LBClient, bool) {
	val, exist := ism.mcMap.Load(key)
	if !exist {
		return nil, false
	}

	syncMap, ok := val.(*fasthttp.LBClient)
	if !ok {
		return nil, false
	}

	return syncMap, true
}

func (ism *InsLbClientSyncMap) Put(key string, val *fasthttp.LBClient) {
	ism.mcMap.Store(key, val)
}
