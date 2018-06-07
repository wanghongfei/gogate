package server

import (
	asynclog "github.com/alecthomas/log4go"
	"github.com/valyala/fasthttp"
)

// 控制QPS的前置过虑器
func RateLimitPreFilter(s *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	// 取出router结果
	info, ok := GetServiceInfoFromUserValue(ctx, ROUTE_INFO)
	if !ok {
		return true
	}

	// 取出对应service的限速器
	if 0 == info.Qps {
		// 如果没有说明不需要限速
		asynclog.Debug("no limiter for service %s", info.Id)
		return true
	}

	// 取出限速器
	rl, ok := s.rateLimiterMap.Get(info.Id)
	if !ok {
		asynclog.Error("lack rate limiter for %s", info.Id)
		return true
	}

	pass := rl.TryAcquire()
	if !pass {
		// token不足
		NewResponse(ctx.UserValue(REQUEST_PATH).(string), "reach QPS limitation").Send(ctx)
		asynclog.Info("drop request for %s due to rate limitation", info.Id)
	}

	return pass
}
