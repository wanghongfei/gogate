package utils

import (
	"math/rand"
	"time"
)

var globalRand *rand.Rand

func init() {
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 根据权重值生成随机落点;
// weight: 权重数组, 例如传1,2,3, 则会按1:2:3的概率生成
// return: 此次落点的索引, 索引值对应于传入参数weight
func RandomByWeight(weight []int) int {
	// 计算出生成随机数时的范围最大值
	max := 0
	for _, w := range weight {
		// 乘以10进行放大
		max += w * 10
	}

	// 生成随机数
	stone := 0
	for stone == 0 {
		stone = globalRand.Intn(max)
	}

	// 判断随机数落在了哪个区间
	sum := 0
	for ix, w := range weight {
		start := sum
		sum += w * 10

		if stone > start && stone <= sum {
			return ix
		}
	}

	return -1
}
