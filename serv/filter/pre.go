package filter

import "github.com/valyala/fasthttp"

func ServiceMatchPreFilter(ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool {
	return true
}
