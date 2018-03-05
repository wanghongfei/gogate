package serv

import (
	"strings"

	"code.google.com/p/log4go"
	"github.com/valyala/fasthttp"
)

func ServiceMatchPreFilter(s *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	uri := string(ctx.URI().Path())

	appId := s.Router.Match(uri)
	if "" == appId {
		// 没匹配到
		ctx.WriteString("no match")
		return false
	}

	ctx.SetUserValue(SERVICE_NAME, strings.ToUpper(appId))

	log4go.Debug("%s matched to %s", uri, appId)

	return true
}
