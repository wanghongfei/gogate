package server

import (
	"github.com/alecthomas/log4go"
	"github.com/valyala/fasthttp"
)

func UrlRewritePreFilter(s *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	info := ctx.UserValue(ROUTE_INFO).(*ServiceInfo)


	if info.StripPrefix {
		// path中去掉prefix
		original := string(newRequest.URI().Path())
		posToStrip := len(info.Prefix)

		newPath := original[posToStrip:]
		if newPath == "" {
			newPath = "/"
		}
		newRequest.URI().SetPath(newPath)

		log4go.Debug("rewrite path from %s to %s", original, newPath)
	}

	return true
}
