package throttle

import (
	"github.com/wanghongfei/gogate/perr"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/wanghongfei/gogate/redis"
)

type RedisRateLimiter struct {
	qps				string
	client			*redis.RedisClient
	luaCode			string

	serviceId		string
	luaArgs			[]string
}

func NewRedisRateLimiter(client *redis.RedisClient, luaPath string, qps int, serviceId string) (*RedisRateLimiter, error) {
	if nil == client {
		return nil, perr.WrapSystemErrorf(nil, "redis client cannot be nil")
	}

	if qps < 1 {
		qps = 1
	}

	if !client.IsConnected() {
		err := client.Connect()
		if nil != err {
			return nil, perr.WrapSystemErrorf(err, "failed to connect to redis")
		}
	}

	luaF, err := os.Open(luaPath)
	if nil != err {
		return nil, err
	}
	defer luaF.Close()

	luaBuf, _ := ioutil.ReadAll(luaF)
	luaCode := string(luaBuf)

	qpsStr := strconv.Itoa(qps)

	return &RedisRateLimiter{
		client: client,
		luaCode: luaCode,
		qps: qpsStr,
		serviceId: serviceId,
		luaArgs: []string{serviceId, qpsStr},
	}, nil
}

func (rrl *RedisRateLimiter) Acquire() {
	for {
		ok := rrl.TryAcquire()
		if ok {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func (rrl *RedisRateLimiter) TryAcquire() bool {
	resp, _ := rrl.client.ExeLuaInt(rrl.luaCode, nil, rrl.luaArgs)
	// resp, _ := rrl.client.ExeLuaInt(rrl.luaCode, nil, []string{rrl.serviceId, rrl.qps})
	return resp == 1
}

