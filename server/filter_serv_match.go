package server

import (
	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/asynclog"
)

func ServiceMatchPreFilter(s *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	uri := GetStringFromUserValue(ctx, REQUEST_PATH)

	servInfo := s.Router.Match(uri)
	if nil == servInfo {
		// 没匹配到
		ctx.Response.SetStatusCode(404)
		ctx.WriteString("no match")
		return false
	}
	ctx.SetUserValue(ROUTE_INFO, servInfo)
	ctx.SetUserValue(SERVICE_NAME, servInfo.Id)

	asynclog.Debug("%s matched to %s", uri, servInfo.Id)

	return true
}
