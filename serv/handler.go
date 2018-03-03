package serv

import (
	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/serv/filter"
)

func (s *Server) HandleRequest(ctx *fasthttp.RequestCtx) {
	newReq := new(fasthttp.Request)
	ctx.Request.CopyTo(newReq)

	// 调用Pre过虑器
	invokePreFilters(s.preFilters, ctx, newReq)

	// todo 发请求
	resp := new(fasthttp.Response)

	// 调用Post过虑器
	invokePostFilters(s.postFilters, newReq, resp)
}

func (s *Server) sendRequest(req *fasthttp.Request) {
	// 获取Client
}

func invokePreFilters(filters []filter.PreFilterFunc, ctx *fasthttp.RequestCtx, newReq *fasthttp.Request) {
	for _, f := range filters {
		next := f(ctx, newReq)
		if !next {
			return
		}
	}

}

func invokePostFilters(filters []filter.PostFilterFunc, newReq *fasthttp.Request, resp *fasthttp.Response) {
	for _, f := range filters {
		next := f(newReq, resp)
		if !next {
			return
		}
	}

}
