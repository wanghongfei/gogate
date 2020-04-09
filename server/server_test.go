package server

import (
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	// 模拟两个下游服务
	server1 := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString("server1 at 8080")
			ctx.SetConnectionClose()
		},
	}
	server2 := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString("server2 at 8081")
			ctx.SetConnectionClose()
		},
	}

	// 启动服务
	go func() {
		server1.ListenAndServe("0.0.0.0:8080")
	}()

	go func() {
		server2.ListenAndServe("0.0.0.0:8081")
	}()


	// 启动gogate
	InitGogate("gogate-test.yml")
	gogate, err := NewGatewayServer("localhost", 7000, "route-test.yml", 10)
	if nil != err {
		t.Fatal(err)
	}
	go func() {
		err := gogate.Start()
		if nil != err {
			t.Fatal(err)
		}
	}()

	time.Sleep(time.Second)

	// 发请求
	_, buf, err := fasthttp.Get(make([]byte, 0, 20), "http://localhost:7000/service1/info")
	if nil != err {
		t.Fatal(err)
	}
	if "server1 at 8080" != string(buf) {
		t.Error("service 1 failed")
	}

	// 发请求
	_, buf, err = fasthttp.Get(make([]byte, 0, 20), "http://localhost:7000/service2/info")
	if nil != err {
		t.Fatal(err)
	}
	if "server2 at 8081" != string(buf) {
		t.Error("service 2 failed")
	}


	server1.Shutdown()
	server2.Shutdown()
	gogate.Shutdown()
}


