package server

import (
	"github.com/valyala/fasthttp"
)

// 前置过滤器函数
type PreFilterFunc func(server *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool
// 前置过滤器对象
type PreFilter struct {
	FilterFunc			PreFilterFunc
	Name				string
}

func NewPreFilter(name string, filter PreFilterFunc) *PreFilter {
	return &PreFilter{
		FilterFunc: filter,
		Name: name,
	}
}

func (pf *PreFilter) String() string {
	return pf.Name
}


// 后置过滤器函数
type PostFilterFunc func(req *fasthttp.Request, resp *fasthttp.Response) bool
// 后置过滤器对象
type PostFilter struct {
	FilterFunc			PostFilterFunc
	Name				string
}

func NewPostFilter(name string, filter PostFilterFunc) *PostFilter {
	return &PostFilter{
		FilterFunc: filter,
		Name: name,
	}
}

func (pf *PostFilter) String() string {
	return pf.Name
}
