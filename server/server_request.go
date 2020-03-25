package server

import (
	. "github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/discovery"
	"github.com/wanghongfei/gogate/server/route"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/utils"
)

// 转发请求到指定微服务
// return:
// Response: 响应对象;
// string: 下游服务名
// error: 错误
func (serv *Server) sendRequest(ctx *fasthttp.RequestCtx, req *fasthttp.Request) (*fasthttp.Response, string, error) {
	// 获取服务信息
	info := ctx.UserValue(ROUTE_INFO).(*route.ServiceInfo)

	var logRecordName string
	// 需要从注册列表中查询地址
	if info.Id != "" {
		logRecordName = info.Id

		// 获取Client
		appId := strings.ToUpper(info.Id)

		// 灰度, 选择版本
		version := chooseVersion(info.Canary)

		// 取出指定服务的所有实例
		serviceInstances, exist := serv.registryMap.Get(appId)
		if !exist || 0 == len(serviceInstances) {
			// return nil, "", errors.New("no instance " + appId + " for service " + appId + ", (service is offline)")
			return nil, "", utils.Errorf("no instance %s for service (service is offline)", appId)
		}

		// 按version过滤
		if "" != version {
			serviceInstances = filterWithVersion(serviceInstances, version)
		}

		// 负载均衡
		targetInstance := serv.lb.Choose(serviceInstances)
		// 修改请求的host为目标主机地址
		req.URI().SetHost(targetInstance.Addr)

	} else {
		logRecordName = info.Name

		// 直接使用后面的地址
		hostList := strings.Split(info.Host, ",")

		targetAddr := serv.lb.ChooseByAddresses(hostList)
		req.URI().SetHost(targetAddr)
	}

	// 发请求
	resp := new(fasthttp.Response)
	err := fasthttp.Do(req, resp)
	if nil != err {
		return nil, "", utils.Errorf("failed to send request to downstream service => %w", err)
	}

	return resp, logRecordName, nil
}

// 过滤出meta里version字段为指定值的实例
func filterWithVersion(instances []*discovery.InstanceInfo, targetVersion string) []*discovery.InstanceInfo {
	result := make([]*discovery.InstanceInfo, 5, 5)

	for _, ins := range instances {
		if ins.Meta[META_VERSION] == targetVersion {
			result = append(result, ins)
		}
	}

	return result
}

func chooseVersion(canaryInfos []*route.CanaryInfo) string {
	if nil == canaryInfos || len(canaryInfos) == 0 {
		return ""
	}

	var weights []int
	for _, info := range canaryInfos {
		weights = append(weights, info.Weight)
	}

	index := utils.RandomByWeight(weights)
	if -1 == index {
		Log.Warn("random interval returned -1")
		return ""
	}

	return canaryInfos[index].Meta
}
