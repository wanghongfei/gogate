package redis

import (
	"github.com/mediocregopher/radix.v2/cluster"
)

type ClusterRedisClient struct {
	addr			string
	cluster			*cluster.Cluster
}

func NewClusterRedisClient(addr string) *ClusterRedisClient {
	return &ClusterRedisClient{
		addr: addr,
	}
}

func (crd *ClusterRedisClient) GetString(key string) (string, error) {
	resp := crd.cluster.Cmd("get", key)
	if nil != resp.Err {
		return "", resp.Err
	}

	return resp.Str()
}

func (crd *ClusterRedisClient) Close() {
	crd.cluster.Close()
}

// for test only
func (crd *ClusterRedisClient) Connect() error {
	cluster, err := cluster.New(crd.addr)
	if err != nil {
		return err
	}

	crd.cluster = cluster
	return nil
}
