package serv

import "github.com/valyala/fasthttp"

func (s *Server) HandleRequest(ctx *fasthttp.RequestCtx) {
}

func buildRequest(ctx *fasthttp.RequestCtx, newHost string) *fasthttp.Request {
	newReq := &fasthttp.Request{}
	ctx.Request.CopyTo(newReq)

	newReq.SetHost(newHost)

	return newReq
}