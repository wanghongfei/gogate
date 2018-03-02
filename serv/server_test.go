package serv

import (
	"fmt"
	"testing"
	"time"
)

func TestServer_ReloadRoute(t *testing.T) {
	serv, err := NewGatewayServer("127.0.0.1", 8080, "../route.yml")
	if nil != err {
		t.Error(err)
		return
	}
	fmt.Println(serv.ExtractRoute())

	time.Sleep(time.Second * 5)
	serv.ReloadRoute()
	fmt.Println(serv.ExtractRoute())
}
