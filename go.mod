module github.com/wanghongfei/gogate

replace golang.org/x/net => github.com/golang/net v0.0.0-20190404232315-eb5bcb51f2a3

replace golang.org/x/text => github.com/golang/text v0.3.0

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190411191339-88737f569e3a

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190412213103-97732733099d

go 1.12

require (
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/hashicorp/consul/api v1.0.1
	github.com/mediocregopher/radix.v2 v0.0.0-20181115013041-b67df6e626f9
	github.com/valyala/fasthttp v1.3.0
	github.com/wanghongfei/go-eureka-client v1.1.0
	gopkg.in/yaml.v2 v2.2.2
)
