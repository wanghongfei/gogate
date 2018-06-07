package server

import (
	"errors"
	"strings"

	asynclog "github.com/alecthomas/log4go"
	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/utils"
)


func (serv *Server) sendRequest(ctx *fasthttp.RequestCtx, req *fasthttp.Request) (*fasthttp.Response, error) {
	// 获取服务信息
	info := ctx.UserValue(ROUTE_INFO).(*ServiceInfo)

	var c *fasthttp.LBClient
	// 需要从注册列表中查询地址
	if info.Id != "" {
		// 获取Client
		appId := info.Id

		// 灰度, 选择版本
		version := chooseVersion(info.Canary)

		// 构造HTTP client名
		clientName := appId
		if "" != version {
			clientName = clientName + ":" + version
		}

		client, exist := serv.proxyClients.Get(clientName)
		if !exist {
			return nil, errors.New("no client " + clientName + " for service " + appId + ", (service is offline)")
		}

		c = client

	} else {
		// 直接使用后面的地址
		hostList := strings.Split(info.Host, ",")
		c = &fasthttp.LBClient{
			Clients: createClients(hostList),
		}
	}



	// 发请求
	resp := new(fasthttp.Response)
	err := c.Do(req, resp)
	if nil != err {
		return nil, err
	}

	return resp, nil
}

func chooseVersion(canaryInfos []*CanaryInfo) string {
	if nil == canaryInfos || len(canaryInfos) == 0 {
		return ""
	}

	var weights []int
	for _, info := range canaryInfos {
		weights = append(weights, info.Weight)
	}

	index := utils.RandomByWeight(weights)
	if -1 == index {
		asynclog.Warn("random interval returned -1")
		return ""
	}

	return canaryInfos[index].Meta
}
