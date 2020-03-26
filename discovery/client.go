package discovery

// 封装服务实例信息
type InstanceInfo struct {
	ServiceName		string

	// 格式为 host:port
	Addr 			string
	// 此实例附加信息
	Meta 			map[string]string
}

// 服务发现客户端接口
type Client interface {
	// 查询所有服务实例
	QueryServices() ([]*InstanceInfo, error)

	// 注册自己
	Register() error

	// 取消注册自己
	UnRegister() error
}

