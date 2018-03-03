package filter

import "github.com/valyala/fasthttp"

type PreFilterFunc func(ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool

type PostFilterFunc func(req *fasthttp.Request, resp *fasthttp.Response) bool
