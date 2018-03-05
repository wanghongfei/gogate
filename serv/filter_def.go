package serv

import (
	"github.com/valyala/fasthttp"
)

type PreFilterFunc func(server *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool

type PostFilterFunc func(req *fasthttp.Request, resp *fasthttp.Response) bool
