package throttle

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type RateLimiter struct {
	// 每秒生成的token数
	tokenPerSecond		int

	// 生成一个token需要的micro second
	tokenGenMicro		int64
	// 上次生成token的时间
	lastGenMicro		int64

	// 当前桶内token数量
	tokenCount			int
	mutex				*sync.Mutex
}

// 创建限速器
// qps: 每秒最大请求数
func NewRateLimiter(qps int) *RateLimiter {
	if qps < 1 {
		qps = 1
	}

	rl := new(RateLimiter)
	rl.tokenPerSecond = qps
	rl.mutex = new(sync.Mutex)
	rl.tokenGenMicro = int64(int64(1000 * 1000) / int64(qps))

	return rl
}

// 获取token, 如果没有则block
func (rl *RateLimiter) Acquire() {
	rl.mutex.Lock()
	rl.consumeToken(true)
	rl.mutex.Unlock()
}

// 获取token, 成功返回true, 没有则返回false
func (rl *RateLimiter) TryAcquire() bool {
	rl.mutex.Lock()
	got := rl.consumeToken(false)
	rl.mutex.Unlock()

	return got
}

func (rl *RateLimiter) fillBucket() {
	nowMicro := time.Now().UnixNano() / 1000
	// 如果是第一次获取, 直接填满1s的token
	if 0 == rl.lastGenMicro {
		rl.tokenCount = rl.tokenPerSecond
		rl.lastGenMicro = nowMicro
		return
	}

	// 计算上次生成token时的时间差
	microSecondDiff := nowMicro - rl.lastGenMicro
	// 计算应当生成的新token数
	newTokens := microSecondDiff / rl.tokenGenMicro
	rl.tokenCount += int(newTokens)
	// token总数不能超过qps值
	if rl.tokenCount > rl.tokenPerSecond {
		rl.tokenCount = rl.tokenPerSecond
	}

	rl.lastGenMicro = nowMicro

	// fmt.Printf("timeDiff = %v\n", timeDiff)
}

func (rl *RateLimiter) consumeToken(canSleep bool) bool {
	for rl.tokenCount == 0 {
		rl.fillBucket()
		if rl.tokenCount == 0 {
			if canSleep {
				time.Sleep(time.Microsecond * time.Duration(rl.tokenGenMicro))
			} else {
				return false
			}
		}
	}

	rl.tokenCount--
	return true
}

func (rl *RateLimiter) String() string {
	return "qps = " + strconv.Itoa(rl.tokenPerSecond) +
		",tokenGenMicro = " + strconv.FormatInt(rl.tokenGenMicro, 10) +
		",lastGenMicro = " + fmt.Sprintf("%v", rl.lastGenMicro) +
		",tokenCount = " + strconv.Itoa(rl.tokenCount)
}

