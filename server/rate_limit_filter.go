package server

import (
	"code.google.com/p/log4go"
	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/throttle"
)

func RateLimitPreFilter(s *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	// 取出router结果
	ctxVal := ctx.UserValue(ROUTE_INFO)
	if nil == ctxVal {
		return true
	}

	// 取出对应service的限速器
	info := ctxVal.(*ServiceInfo)
	rlVal, ok := s.rateLimiterMap.Load(info.Id)
	if !ok {
		// 如果没有说明不需要限速
		log4go.Debug("no limiter for service %s", info.Id)
		return true
	}

	rl := rlVal.(*throttle.RateLimiter)
	pass := rl.TryAcquire()
	if !pass {
		// 没匹配到
		NewResponse(ctx.UserValue(REQUEST_PATH).(string), "reach QPS limitation").Send(ctx)
		log4go.Info("drop request for %s due to rate limitation", info.Id)
	}

	return pass
}
