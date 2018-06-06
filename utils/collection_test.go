package utils

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	a := make([]int, 0, 10)
	fmt.Printf("cap(a) = %d, len(a) = %d\n", cap(a), len(a))

	b := append(a, 1)
	fmt.Printf("cap(a) = %d, len(a) = %d, cap(b) = %d, len(b) = %d\n", cap(a), len(a), cap(b), len(b))

	_ = append(a, 2)
	fmt.Printf("cap(a) = %d, len(a) = %d, cap(b) = %d, len(b) = %d\n", cap(a), len(a), cap(b), len(b))

	println(b[0])
}

func TestCopy(t *testing.T) {
	arr := make([]int, 0, 6)
	for ix := 0; ix < cap(arr); ix++ {
		arr = append(arr, ix + 1)
	}
	fmt.Println(arr)

	targetIx := 1
	copy(arr[targetIx + 2:], arr[targetIx + 1:])
	arr[targetIx + 1] = 999
	fmt.Println(arr)

}
