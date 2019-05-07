package discovery

// 封装服务实例信息
type InstanceInfo struct {
	ServiceName		string

	// 格式为 host:port
	Addr 			string
	// 此实例附加信息
	Meta 			map[string]string
}
