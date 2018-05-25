package server

import (
	"code.google.com/p/log4go"
	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/throttle"
)

// 控制QPS的前置过虑器
func RateLimitPreFilter(s *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	// 取出router结果
	ctxVal := ctx.UserValue(ROUTE_INFO)
	if nil == ctxVal {
		return true
	}

	// 取出对应service的限速器
	info := ctxVal.(*ServiceInfo)
	if 0 == info.Qps {
		// 如果没有说明不需要限速
		log4go.Debug("no limiter for service %s", info.Id)
		return true
	}

	// 取出限速器
	rlVal, _ := s.rateLimiterMap.Load(info.Id)
	rl := rlVal.(*throttle.RateLimiter)

	pass := rl.TryAcquire()
	if !pass {
		// token不足
		NewResponse(ctx.UserValue(REQUEST_PATH).(string), "reach QPS limitation").Send(ctx)
		log4go.Info("drop request for %s due to rate limitation", info.Id)
	}

	return pass
}
