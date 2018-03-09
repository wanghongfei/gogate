![build](https://travis-ci.org/wanghongfei/spring-cloud-gogate.svg?branch=master)

# GoGate

Go语言实现的Spring Cloud网关，目标是性能，即使用更少的资源达到更高的QPS。

目前GoGate已经实现的功能有:

- 基于Eureka的服务发现、注册
- 请求路由(反向代理)
- 负载均衡

初步测试了一下性能，结论如下：

相同的硬件环境、Zuul充分预热且关闭Hystrix的前提下，Go版的网关QPS为Zuul的2.3倍，同时内存占用仅为Zuul的十分之一(600M vs 50M)。而且Go基本上第一波请求就能达到最大QPS, zuul要预热几次才会稳定。更详细的测试近期更新。



## 路由配置

```yaml
user-service:
  # eureka中的服务名
  id: user-service
  # 以/user开头的请求, 会被转发到user-service服务中
  prefix: /user

dog-service:
  id: dog-service
  # 请求路径当匹配多个prefix时, 长的获胜
  prefix: /user/dog

order-service:
  id: order-service
  prefix: /order

img-service:
  # 如果有host, 则不查注册中心直接使用此地址, 多个地址逗号分隔
  host: localhost:4444,localhost:5555
  prefix: /img

# 上面都没有匹配到时
common-service:
  id: common-service
  prefix: /
```

