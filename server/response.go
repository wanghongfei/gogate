package server

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type GogateResponse struct {
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
	ctx.WriteString(resp.ToJson())
}

func (resp *GogateResponse) SendWithStatus(ctx *fasthttp.RequestCtx, statusCode int) {
	ctx.SetStatusCode(statusCode)
	resp.Send(ctx)
}
