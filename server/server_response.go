package server

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	. "github.com/wanghongfei/gogate/conf"
)

type GogateResponse struct {
	RequestId	int64`json:"requestId"`
	Path		string`json:"path"`
	Error		string`json:"error"`
}

func NewResponse(path, msg string) *GogateResponse {
	return &GogateResponse{
		Path: path,
		Error: msg,
	}
}

func (resp *GogateResponse) ToJson() string {
	return string(resp.ToJsonBytes())
}

func (resp *GogateResponse) ToJsonBytes() []byte {
	buf, _ := json.Marshal(resp)

	return buf
}

func (resp *GogateResponse) Send(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json;charset=utf8")
	resp.RequestId = GetInt64FromUserValue(ctx, REQUEST_ID)
	timer := GetStopWatchFromUserValue(ctx)

	responseBody := resp.ToJson()
	Log.Infof("request %d finished, cost = %dms, statusCode = %d, response = %s", resp.RequestId, timer.Record(), ctx.Response.StatusCode(), responseBody)
	ctx.WriteString(responseBody)
}

func (resp *GogateResponse) SendWithStatus(ctx *fasthttp.RequestCtx, statusCode int) {
	ctx.SetStatusCode(statusCode)
	resp.Send(ctx)
}
