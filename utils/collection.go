package utils

import (
	"sync"
)

/*
* 从map中删除指定的key

* PARAMS:
*	- baseMap: 要删除key的map
*	- keys: 要删除的key数组
*
* RETURNS:
*	集合了2个map所有元素的map指针
*/
func DelKeys(baseMap *sync.Map, keys []interface{}) {
	for _, key := range keys {
		baseMap.Delete(key)
	}
}

/*
* 两个map取并集
*
* PARAMS:
*	- fromMap: 源map
*	- toMap: 合并后的map
*
*/
func MergeSyncMap(fromMap, toMap *sync.Map) {
	fromMap.Range(func(key, value interface{}) bool {
		toMap.Store(key, value)
		return true
	})
}

/*
* 找出在baseMap中存在但yMap中不存在的元素
*
* PARAMS:
*	- baseMap: 独有元素所在的map
*	- yMap: 对比map
*
* RETURNS:
*	baseMap中独有元素的key的数组
*/
func FindExclusiveKey(baseMap, yMap *sync.Map) []interface{} {
	var keys []interface{}
	baseMap.Range(func(key, value interface{}) bool {
		_, exist := yMap.Load(key)
		if !exist {
			keys = append(keys, key)
		}

		return true
	})

	return keys
}
