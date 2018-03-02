package filter

import "github.com/valyala/fasthttp"

func PreRoute(ctx *fasthttp.RequestCtx) (string, string, error) {
	return "", "", nil
}
