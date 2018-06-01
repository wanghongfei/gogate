package throttle

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-samples/chaincode/pkg/errors"
	"github.com/wanghongfei/gogate/redis"
)

type RedisRateLimiter struct {
	qps				string
	client			*redis.RedisClient
	luaCode			string

	serviceId		string
}

func NewRedisRateLimiter(client *redis.RedisClient, luaPath string, qps int, serviceId string) (*RedisRateLimiter, error) {
	if nil == client {
		return nil, errors.New("redis client cannot be nil")
	}

	if qps < 1 {
		qps = 1
	}

	if !client.IsConnected() {
		err := client.Connect()
		if nil != err {
			return nil, err
		}
	}

	luaF, err := os.Open(luaPath)
	if nil != err {
		return nil, err
	}
	defer luaF.Close()

	luaBuf, _ := ioutil.ReadAll(luaF)
	luaCode := string(luaBuf)

	return &RedisRateLimiter{
		client: client,
		luaCode: luaCode,
		qps: strconv.Itoa(qps),
		serviceId: serviceId,
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
	resp, _ := rrl.client.ExeLuaInt(rrl.luaCode, nil, []string{rrl.serviceId, rrl.qps})
	// resp, _ := rrl.client.ExeLuaInt(rrl.luaCode, nil, []string{rrl.serviceId, rrl.qps})
	fmt.Println(resp)
	return resp == 1
}

