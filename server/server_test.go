package server

import (
	"testing"
	"time"
)

func TestServer_ReloadRoute(t *testing.T) {
	serv, err := NewGatewayServer("127.0.0.1", 8080, "../route.yml", 50, false, time.Second)
	if nil != err {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 5)
	serv.ReloadRoute()
}

