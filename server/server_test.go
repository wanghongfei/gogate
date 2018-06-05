package server

import (
	"fmt"
	"testing"
	"time"
)

func TestServer_ReloadRoute(t *testing.T) {
	serv, err := NewGatewayServer("127.0.0.1", 8080, "../route.yml", 50)
	if nil != err {
		t.Error(err)
		return
	}
	fmt.Println(serv.ExtractRoute())

	time.Sleep(time.Second * 5)
	serv.ReloadRoute()
	fmt.Println(serv.ExtractRoute())
}

func TestOther(t *testing.T) {
	a := make([]int, 0, 10)
	fmt.Printf("cap(a) = %d, len(a) = %d\n", cap(a), len(a))

	b := append(a, 1)
	fmt.Printf("cap(a) = %d, len(a) = %d, cap(b) = %d, len(b) = %d\n", cap(a), len(a), cap(b), len(b))

	_ = append(a, 2)
	fmt.Printf("cap(a) = %d, len(a) = %d, cap(b) = %d, len(b) = %d\n", cap(a), len(a), cap(b), len(b))

	println(b[0])
}
