package serv

import (
	"strings"

	"code.google.com/p/log4go"
	"github.com/valyala/fasthttp"
)

func ServiceMatchPreFilter(s *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	uri := ctx.UserValue(REQUEST_PATH).(string)

	servInfo := s.Router.Match(uri)
	if nil == servInfo {
		// 没匹配到
		ctx.Response.SetStatusCode(404)
		ctx.WriteString("no match")
		return false
	}

	addr := ""
	if "" != servInfo.Host {
		addr = "HOST:" + servInfo.Host
		ctx.SetUserValue(SERVICE_NAME, strings.ToUpper(addr))

	} else {
		addr = "ID:" + servInfo.Id
	}

	log4go.Debug("%s matched to %s", uri, addr)

	return true
}
