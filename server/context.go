package server

import (
	"github.com/valyala/fasthttp"
)

// 从请求上下文中取出*ServiceInfo
func GetServiceInfoFromUserValue(ctx *fasthttp.RequestCtx, key string) (*ServiceInfo, bool) {
	val := ctx.UserValue(key)
	if nil == val {
		return nil, false
	}

	info, ok := val.(*ServiceInfo)
	if !ok {
		return nil, false
	}

	return info, true
}

// 从请求上下文中取出string
func GetStringFromUserValue(ctx *fasthttp.RequestCtx, key string) string {
	val := ctx.UserValue(key)
	if nil == val {
		return ""
	}

	str, ok := val.(string)
	if !ok {
		return ""
	}

	return str
}
