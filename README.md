![build](https://api.travis-ci.org/wanghongfei/gogate.svg?branch=master)

# GoGate

Go语言实现的Spring Cloud网关，目标是性能，即使用更少的资源达到更高的QPS。

GoGate使用以高性能著称的`FastHttp`库收发HTTP请求，且会为每个host单独创建一个`HostClient`以减少锁竞争。



目前已经实现的功能有:

- 基于Eureka的服务发现、注册
- 请求路由、路由配置热更新
- 负载均衡
- 灰度发布(基于Eureka meta信息里的version字段分配流量)
- 微服务粒度的QPS控制(有基于内存的令牌桶算法限流和Redis + Lua限流两种可选)
- 微服务粒度的流量统计(暂时实现为记录日志到/tmp目录下)
- 优雅关闭(开启此功能会略微损耗性能)

初步测试了一下性能，结论如下：

相同的硬件环境、Zuul充分预热且关闭Hystrix的前提下，Go版的网关QPS为Zuul的2.3倍，同时内存占用仅为Zuul的十分之一(600M vs 50M)。而且Go基本上第一波请求就能达到最大QPS, zuul要预热几次才会稳定。

如果按消耗相同资源的前提下算的话，go一定要比zuul节省多的多的多的机器。



## 什么情况下可以考虑使用非Java语言的网关

- 系统使用Spring Cloud全家桶
- 对Zuul 1性能不满意
- 对Cloud官方已经明确不会整合Zuul 2的行为不爽
- 认为Spring Cloud Gateway不够成熟(相比Zuul 2.0)
- 对网关的CPU/内存资源使用非常敏感

那么就可以考虑！



## 流程

![arc](http://ovbyjzegm.bkt.clouddn.com/gogate-arc.jpg)

服务路由`service-match-pre-filter`: 根据URL匹配后端微服务

流量控制`rate-limit-pre-filter`: 令牌桶算法控制qps

URL重写`url-rewrite-pre-filter`: 调整向后端服务发请求的URL

转发请求: 负载均衡、按比例分配流量



gogate没有提供默认的Post Filter，可根据需要自己实现相应函数。



## 构建

此项目使用go官方的依赖解决方法modules, 因此需要版本>=1.11。



在$GOPATH之外的任意目录下clone项目:

```shell
git clone https://github.com/wanghongfei/gogate
```

安装依赖:

```shell
go mod tidy
```

最后构建:

```shell
go build
```





## 使用

可以编译`main.go`直接生成可执行文件，也可以当一个库来使用。

可以在转发请求之前和之后添加自定义Filter来添加自定义逻辑。

详见`examples/usage.go`



## 关于优雅关闭

创建`server`对象时如果指定了启用优雅关闭功能, 则在调用

```go
server.WaitForGracefullyClose()
```

方法后，会发生：

- 停止eureka心跳
- 向eureka取消注册
- block当前协程，直到所有正在处理的请求退出或者超过最大等待时间




## 路由配置

路由匹配规则:

- 当`id`不为空时，会使用eureka的注册信息查询此服务的地址
- 当`host`不为空时, 会优先使用此字段指定的服务地址, 多个地址用逗号分隔
- 当请求路径匹配多个`prefix`时，配置文件中`prefix`最长的获胜

当路由配置文件发生变动时，访问

```
GET /_mgr/reload
```

即可应用新配置。



示例配置：

```yaml
services:
  user-service:
    # eureka中的服务名
    id: user-service
    # 以/user开头的请求, 会被转发到user-service服务中
    prefix: /user
    # 转发时是否去掉请求前缀, 即/user
    strip-prefix: true
    # 灰度配置
    canary:
      -
        # 对应eurekai注册信息中元数据(metadata map)中key=version的值
        meta: "1.0"
        # 流量比重
        weight: 3
      -
        meta: "2.0"
        weight: 4
      -
        # 对应没有metadata的服务
        meta: ""
        weight: 1

  trends-service:
    id: trends-service
    # 请求路径当匹配多个prefix时, 长的获胜
    prefix: /trends
    strip-prefix: false
    # 设置qps限制, 每秒最多请求数
    qps: 1

  order-service:
    id: order-service
    prefix: /order
    strip-prefix: false

  img-service:
    # 如果有host, 则不查注册中心直接使用此地址, 多个地址逗号分隔
    host: localhost:4444,localhost:5555
    prefix: /img
    strip-prefix: false

# 上面都没有匹配到时
  common-service:
    id: common-service
    prefix: /
    strip-prefix: false
```



## 自定义过滤器

前置fitler和后置filter都可以在任意位置添加自定义过滤器以实现定制化的功能。



- 前置过滤器

函数签名为:

```go
type PreFilterFunc func(server *Server, ctx *fasthttp.RequestCtx, newRequest *fasthttp.Request) bool
```

`server`: gogate server对象的指针

`ctx`: 请求上下文对象指针

`newRequest`: 要转发给下游微服务的请求对象的指针，可以对相关参数进行修改，如header, body, method等

返回`true`时gogate会继续触发下一个过滤器，返回`false`则表示请求到此为止， 不会执行后续过滤器，也不会转发请求。



- 后置过滤器

函数签名名:

```go
type PostFilterFunc func(req *fasthttp.Request, resp *fasthttp.Response) bool
```

`req`: 已经转发给微服务的请求对象指针

`resp`: 微服务返回的响应对象指针, 可进行修改

返回`true`时gogate会继续触发下一个过滤器，返回`false`则表示请求到此为止， 不会执行后续过滤器，也不会转发请求。



- 添加过滤器

`Server.AppendPreFilter`: 在末尾追加前置过滤器

`Server.AppendPostFilter`: 在末尾追加后置过滤器

`Server.InsertPreFilter`: 在指定过滤器的后面插入前置过滤器

`Server.InsertPostFilter`: 在指定过滤器的后面插入后置过滤器

`Server.InsertPreFilterAhead`: 插入前置过滤器到最头部

`Server.InsertPostFilterAhead`: 插入后置过滤器到最头部 

## Eureka配置

`eureka.json`文件



## gogate配置

`gogate.yml`文件:

```yaml
version: 1.0

server:
  # 向eureka注册自己时使用的服务名
  appName: gogate
  host: 127.0.0.1
  port: 8080
  # gateway最大连接数
  maxConnection: 1000
  # gateway请求后端服务超时时间, 毫秒
  timeout: 3000

eureka:
  # eureka配置文件名
  configFile: eureka.json
  # 路由配置文件名
  routeFile: route.yml
  # eureka剔除服务的最大时间限值, 秒
  evictionDuration: 30
  # 心跳间隔, 秒
  heartbeatInterval: 20


traffic:
  # 是否开启流量记录功能
  enableTrafficRecord: true
  # 流量日志文件所在目录
  trafficLogDir: /tmp

redis:
  # 是否使用redis做限速器
  enabled: false
  # 目前只支持单实例, 不支持cluster
  addr: 127.0.0.1:6379
  # 限速器lua代码文件
  rateLimiterLua: lua/rate_limiter.lua
```



## 限流器

gogate有两个限流器实现, `MemoryRateLimiter`和`RedisRateLimiter`，通过`gogate.yml`配置文件里的`redis.enabled`控制。前者使用令牌桶算法实现，适用于单实例部署的场景；后者基于 Redis + Lua 实现，适用于多实例部署。但有一个限制是目前Redis只支持连接单个实例，不支持cluster。



## 流量日志

gogate会记录过去`1s`内各个微服务的请求数据，包括成功请求数和失败请求数，然后写入`/tmp/{service-id}_yyyyMMdd.log`文件中:

```
1527580599228,2,1,user-service
1527580600230,4,1,user-service
1527580601228,1,1,user-service
```

即`毫秒时间戳,成功请求数,失败请求数,服务名`。

如果在过去的1s内没有请求, 则不会向日志中写入任何数据。

