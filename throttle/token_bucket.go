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
	tokenGenTime		int64
	// 上次生成token的时间
	lastGenTime			*time.Time

	// 当前桶内token数量
	tokenCount			int
	mutex				*sync.Mutex
}

func NewRateLimiter(qps int) *RateLimiter {
	if qps < 1 {
		qps = 1
	}

	rl := new(RateLimiter)
	rl.tokenPerSecond = qps
	rl.mutex = new(sync.Mutex)
	rl.tokenGenTime = int64(int64(1000 * 1000) / int64(qps))

	return rl
}

func (rl *RateLimiter) Acquire() {
	rl.mutex.Lock()
	rl.consumeToken()
	rl.mutex.Unlock()
}

func (rl *RateLimiter) fillBucket() {
	now := time.Now()
	// 如果是第一次获取, 直接填满1s的token
	if nil == rl.lastGenTime {
		rl.tokenCount = rl.tokenPerSecond
		rl.lastGenTime = &now
		return
	}

	// 计算上次生成token时的时间差
	timeDiff := now.Sub(*rl.lastGenTime)
	// 转成micro second
	microSecondDiff := timeDiff.Nanoseconds() / 1000
	// 计算应当生成的新token数
	newTokens := microSecondDiff / rl.tokenGenTime
	rl.tokenCount += int(newTokens)
	// token总数不能超过qps值
	if rl.tokenCount > rl.tokenPerSecond {
		rl.tokenCount = rl.tokenPerSecond
	}

	rl.lastGenTime = &now

	// fmt.Printf("timeDiff = %v\n", timeDiff)
}

func (rl *RateLimiter) consumeToken() {
	for rl.tokenCount == 0 {
		rl.fillBucket()
		if rl.tokenCount == 0 {
			time.Sleep(time.Microsecond * time.Duration(rl.tokenGenTime))
		}
	}

	rl.tokenCount--
}

func (rl *RateLimiter) String() string {
	return "qps = " + strconv.Itoa(rl.tokenPerSecond) +
		",tokenGenTime = " + strconv.FormatInt(rl.tokenGenTime, 10) +
		",lastGenTime = " + fmt.Sprintf("%v", rl.lastGenTime) +
		",tokenCount = " + strconv.Itoa(rl.tokenCount)
}

