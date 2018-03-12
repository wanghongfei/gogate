package utils

import (
	"fmt"
	"testing"
)

func TestRandomByWeight(t *testing.T) {
	var weights = []int {
		//80, 20, 10,
		1, 1, 3,
	}

	stoneMap := make(map[int]int)
	for ix := 0; ix < 100; ix++ {
		next := RandomByWeight(weights)
		stoneMap[next]++
	}

	fmt.Println(stoneMap)
}

func BenchmarkRandomByWeight(b *testing.B) {
	var weights = []int {
		//80, 20, 10,
		1, 1, 3,
	}

	for i := 0; i < b.N; i++ {
		RandomByWeight(weights)
	}
}