package server

import (
	"code.google.com/p/log4go"
	"github.com/valyala/fasthttp"
)

func RateLimitPreFilter(s *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	// 取出router结果
	ctxVal := ctx.UserValue(ROUTE_INFO)
	if nil == ctxVal {
		return true
	}

	// 取出对应service的限速器
	info := ctxVal.(*ServiceInfo)
	rl, ok := s.rateLimiterMap[info.Id]
	if !ok {
		// 如果没有说明不需要限速
		log4go.Debug("no limiter for service %s", info.Id)
		return true
	}

	return rl.TryAcquire()
}
