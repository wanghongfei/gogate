package server

import (
	"strconv"

	"github.com/alecthomas/log4go"
	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/utils"
)

const (
	SERVICE_NAME = "key_service_name"
	REQUEST_PATH = "key_request_path"
	ROUTE_INFO = "key_route_info"

	RELOAD_PATH = "/_mgr/reload"
)

// HTTP请求处理方法
func (s *Server) HandleRequest(ctx *fasthttp.RequestCtx) {
	defer func() {
		if r := recover(); r != nil {
			log4go.Error(r)
			processPanic(ctx, s)
		}
	}()

	// 取出请求path
	path := string(ctx.Path())
	ctx.SetUserValue(REQUEST_PATH, path)

	log4go.Info("request received: %s %s", string(ctx.Method()), path)

	// 处理reload请求
	if path == RELOAD_PATH {
		err := s.ReloadRoute()
		if nil != err {
			log4go.Error(err)
			NewResponse(path, err.Error()).Send(ctx)
			return
		}

		ctx.WriteString(s.ExtractRoute())
		return
	}

	newReq := new(fasthttp.Request)
	ctx.Request.CopyTo(newReq)

	// 调用Pre过虑器
	ok := invokePreFilters(s, ctx, newReq)
	if !ok {
		return
	}

	// 发请求
	sw := utils.NewStopwatch()
	resp, err := s.sendRequest(ctx, newReq)
	if nil != err {
		log4go.Error(err)
		NewResponse(path, err.Error()).Send(ctx)

		s.recordTraffic(ctx, false)
		return
	}
	s.recordTraffic(ctx, true)

	resp.Header.Add("Time", strconv.FormatInt(sw.Record(), 10))

	// 调用Post过虑器
	ok = invokePostFilters(s, newReq, resp)
	if !ok {
		return
	}

	// 返回响应
	sendResponse(ctx, resp)
}

func sendResponse(ctx *fasthttp.RequestCtx, resp *fasthttp.Response) {
	// copy header
	ctx.Response.Header = resp.Header
	ctx.Response.Header.Add("proxy", "gogate")

	ctx.Write(resp.Body())
}


func invokePreFilters(s *Server, ctx *fasthttp.RequestCtx, newReq *fasthttp.Request) bool {
	for _, f := range s.preFilters {
		next := f.FilterFunc(s, ctx, newReq)
		if !next {
			return false
		}
	}

	return true
}

func invokePostFilters(s *Server, newReq *fasthttp.Request, resp *fasthttp.Response) bool {
	for _, f := range s.postFilters {
		next := f.FilterFunc(newReq, resp)
		if !next {
			return false
		}
	}

	return true
}

func processPanic(ctx *fasthttp.RequestCtx, serv *Server) {
	path := string(ctx.Path())
	NewResponse(path, "system error").SendWithStatus(ctx, 500)

	// 记录流量
	serv.recordTraffic(ctx, false)

}
