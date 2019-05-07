package server

import (
	"errors"
	log "github.com/wanghongfei/gogate/asynclog"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/wanghongfei/gogate/utils"
)

// 转发请求到指定微服务
func (serv *Server) sendRequest(ctx *fasthttp.RequestCtx, req *fasthttp.Request) (*fasthttp.Response, error) {
	// 获取服务信息
	info := ctx.UserValue(ROUTE_INFO).(*ServiceInfo)

	// 需要从注册列表中查询地址
	if info.Id != "" {
		// 获取Client
		appId := strings.ToUpper(info.Id)

		// 灰度, 选择版本
		version := chooseVersion(info.Canary)

		// 取出指定服务的所有实例
		serviceInstances, exist := serv.registryMap.Get(appId)
		if !exist {
			return nil, errors.New("no instance " + appId + " for service " + appId + ", (service is offline)")
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
		// 直接使用后面的地址
		// todo 有优化空间, 不需要每次都new
		hostList := strings.Split(info.Host, ",")
		instances := createInstanceInfos(hostList)

		target := serv.lb.Choose(instances)
		req.URI().SetHost(target.Addr)
	}

	// 发请求
	resp := new(fasthttp.Response)
	err := fasthttp.Do(req, resp)
	if nil != err {
		return nil, err
	}

	return resp, nil
}

// 过滤出meta里version字段为指定值的实例
func filterWithVersion(instances []*InstanceInfo, targetVersion string) []*InstanceInfo {
	result := make([]*InstanceInfo, 5, 5)

	for _, ins := range instances {
		if ins.Meta[META_VERSION] == targetVersion {
			result = append(result, ins)
		}
	}

	return result
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
		log.Warn("random interval returned -1")
		return ""
	}

	return canaryInfos[index].Meta
}
