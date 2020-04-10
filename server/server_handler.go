package server

import (
	"github.com/valyala/fasthttp"
	. "github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/perr"
	"github.com/wanghongfei/gogate/utils"
	"strconv"
)

const (
	SERVICE_NAME = "key_service_name"
	REQUEST_PATH = "key_request_path"
	ROUTE_INFO   = "key_route_info"

	RELOAD_PATH = "/_mgr/reload"
)

// HTTP请求处理方法.
func (serv *Server) HandleRequest(ctx *fasthttp.RequestCtx) {
	defer recoverPanic(ctx, serv)

	// 计时器
	sw := utils.NewStopwatch()


	// 取出请求path
	path := string(ctx.Path())
	ctx.SetUserValue(REQUEST_PATH, path)

	// log.Info("request received: %s %s", string(ctx.Method()), path)

	// 处理reload请求
	if path == RELOAD_PATH {
		err := serv.ReloadRoute()
		if nil != err {
			Log.Error(err)
			NewResponse(path, err.Error()).Send(ctx)
			return
		}

		// ctx.WriteString(serv.ExtractRoute())
		ctx.WriteString("ok")
		return
	}

	newReq := new(fasthttp.Request)
	ctx.Request.CopyTo(newReq)

	// 调用Pre过虑器
	ok := invokePreFilters(serv, ctx, newReq)
	if !ok {
		return
	}

	// 发请求
	resp, logRecordName, err := serv.sendRequest(ctx, newReq)
	// 错误处理
	if nil != err {
		// 解析错误类型
		bizErr, sysErr, _ := perr.ParseError(err)
		var responseMessage string
		if nil != bizErr {
			// 业务错误
			responseMessage = bizErr.Msg
			Log.Error(bizErr.ErrorWithEnv())

		} else if nil != sysErr {
			// 系统错误
			responseMessage = "system error"
			Log.Error(sysErr.ErrorWithEnv())

		} else {
			responseMessage = err.Error()
			Log.Error(err)
		}

		NewResponse(path, responseMessage).Send(ctx)

		serv.recordTraffic(logRecordName, false)
		return
	}
	serv.recordTraffic(logRecordName, true)


	// 调用Post过虑器
	ok = invokePostFilters(serv, newReq, resp)
	if !ok {
		return
	}

	timeCost := sw.Record()
	resp.Header.Add("Time", strconv.FormatInt(timeCost, 10))
	resp.Header.Set("Server", "gogate")

	Log.Infof("request %s finished, time = %vms, response = %s", path, timeCost, resp.Body())

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
	serv.recordTraffic(GetStringFromUserValue(ctx, SERVICE_NAME), false)

}

func recoverPanic(ctx *fasthttp.RequestCtx, serv *Server) {
	if r := recover(); r != nil {
		Log.Errorf("panic: %v", r)
		processPanic(ctx, serv)
	}
}

