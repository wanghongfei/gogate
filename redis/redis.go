package redis

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
)

// Redis Client, 只能连接一个redis实例, 有连接池
type RedisClient struct {
	addr			string
	poolSize		int
	connPool		*pool.Pool

	isConnected		bool
}

func NewRedisClient(addr string, poolSize int) *RedisClient {
	if poolSize < 1 {
		poolSize = 1
	}

	return &RedisClient{
		addr: addr,
		poolSize: poolSize,
	}
}

func (crd *RedisClient) GetString(key string) (string, error) {
	resp := crd.connPool.Cmd("get", key)
	if nil != resp.Err {
		return "", fmt.Errorf("failed to GetString => %w", resp.Err)
	}

	return resp.Str()
}

func (crd *RedisClient) ExeLuaInt(lua string, keys []string, args []string) (int, error) {
	resp := crd.connPool.Cmd("eval", lua, len(keys), keys, args)
	if nil != resp.Err {
		return 0, resp.Err
	}

	return resp.Int()
}

func (crd *RedisClient) Close() {
	crd.connPool.Empty()
	crd.isConnected = false
}

func (crd *RedisClient) IsConnected() bool {
	return crd.isConnected
}

func (crd *RedisClient) Connect() error {
	conn, err := pool.New("tcp", crd.addr, crd.poolSize)
	if err != nil {
		return fmt.Errorf("failed to connect to redis => %w", err)
	}

	crd.connPool = conn
	crd.isConnected = true
	return nil
}
