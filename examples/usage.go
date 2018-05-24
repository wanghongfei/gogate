package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/discovery"
	serv "github.com/wanghongfei/gogate/server"
)

func main() {
	// 初始化
	conf.LoadConfig("gogate.json")
	conf.InitLog("log.xml")
	discovery.InitEurekaClient()

	// 创建Server
	server, err := serv.NewGatewayServer(conf.App.Host, conf.App.Port, conf.App.RouteConfig, conf.App.MaxConnection)
	if nil != err {
		fmt.Println(err)
		return
	}

	// optional: 注册自定义过虑器, 在转发请求之前调用
	server.RegisterPreFilter(PreLogFilter)
	server.RegisterPreFilter(PreLogFilter)
	// optional: 注册自定义过虑器, 在转发请求之后调用
	server.RegisterPostFilter(PostLogFilter)
	server.RegisterPostFilter(PostLogFilter)

	// 启动Server
	err = server.Start()
	if nil != err {
		fmt.Println(err)
		return
	}
}

// 此方法会在gogate转发请求之前调用
// server: gogate服务器对象
// ctx: fasthttp请求上下文
// newRequest: gogate在转发请求时使用的请求对象指针, 可以做一些修改, 比如改请求参数,添加请求头之类
// return: 返回true则会继续执行后面的过虑器(如果有的话), 返回false则不会执行
func PreLogFilter(server *serv.Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	fmt.Println("request path: " + ctx.URI().String())

	return true
}

// 此方法会在gogate转发请求之后调用
// req: 转发给上游服务的HTTP请求
// resp: 上游服务的响应
// return: 返回true则会继续执行后面的过虑器(如果有的话), 返回false则不会执行
func PostLogFilter(req *fasthttp.Request, resp *fasthttp.Response) bool {
	fmt.Println("response: " + resp.String())

	return true
}
