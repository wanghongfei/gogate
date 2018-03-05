package serv

import (
	"strings"

	"code.google.com/p/log4go"
	"github.com/valyala/fasthttp"
)

func ServiceMatchPreFilter(s *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	uri := ctx.UserValue(REQUEST_PATH).(string)

	appId := s.Router.Match(uri)
	if "" == appId {
		// 没匹配到
		ctx.Response.SetStatusCode(404)
		ctx.WriteString("no match")
		return false
	}

	ctx.SetUserValue(SERVICE_NAME, strings.ToUpper(appId))

	log4go.Debug("%s matched to %s", uri, appId)

	return true
}
