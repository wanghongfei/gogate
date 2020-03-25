package server

import (
	"github.com/valyala/fasthttp"
	. "github.com/wanghongfei/gogate/conf"
)

func UrlRewritePreFilter(s *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	info, ok := GetServiceInfoFromUserValue(ctx, ROUTE_INFO)
	if !ok {
		return true
	}


	if info.StripPrefix {
		// path中去掉prefix
		original := string(newRequest.URI().Path())
		posToStrip := len(info.Prefix)

		newPath := original[posToStrip:]
		if newPath == "" {
			newPath = "/"
		}
		newRequest.URI().SetPath(newPath)

		Log.Debug("rewrite path from %s to %s", original, newPath)
	}

	return true
}
