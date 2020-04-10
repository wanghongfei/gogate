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
	REQUEST_ID   = "key_req_id"
	STOPWATCH    = "key_stopwatch"

	RELOAD_PATH = "/_mgr/reload"
)

// HTTP请求处理方法.
func (serv *Server) HandleRequest(ctx *fasthttp.RequestCtx) {
	defer recoverPanic(ctx, serv)

	// 计时器
	sw := utils.NewStopwatch()
	ctx.SetUserValue(STOPWATCH, sw)


	// 取出请求path
	path := string(ctx.Path())
	ctx.SetUserValue(REQUEST_PATH, path)

	// 生成唯一id
	reqId := utils.GenerateUuid()
	ctx.SetUserValue(REQUEST_ID, reqId)

	Log.Infof("request %d received, method = %s, path = %s, body = %s", reqId, string(ctx.Method()), path, string(ctx.Request.Body()))

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
			Log.Errorf("request %d, %s", reqId, bizErr.ErrorWithEnv())

		} else if nil != sysErr {
			// 系统错误
			responseMessage = "system error"
			Log.Errorf("request %d, %s", reqId, sysErr.ErrorWithEnv())
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)

		} else {
			responseMessage = err.Error()
			Log.Errorf("request %d, %s", reqId, err)
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

	// 返回响应
	sendResponse(ctx, resp, reqId, sw)

}

func sendResponse(ctx *fasthttp.RequestCtx, resp *fasthttp.Response, reqId int64, timer *utils.Stopwatch) {
	// copy header
	ctx.Response.Header = resp.Header
	ctx.Response.Header.Add("proxy", "gogate")

	timeCost := timer.Record()
	resp.Header.Add("Time", strconv.FormatInt(timeCost, 10))
	resp.Header.Set("Server", "gogate")

	Log.Infof("request %d finished, cost = %dms, statusCode = %d, response = %s", reqId, timeCost, ctx.Response.StatusCode(), string(resp.Body()))
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

