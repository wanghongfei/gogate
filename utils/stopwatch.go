package utils

import "time"

type Stopwatch struct {
	start	time.Time
}

// 创建一个计时器
func NewStopwatch() *Stopwatch {
	return &Stopwatch{time.Now()}
}

// 返回从上次调用Record()到现的经过的毫秒数
func (st *Stopwatch) Record() int64 {
	now := time.Now()
	diff := now.Sub(st.start).Nanoseconds() / 1000 / 1000

	st.start = now

	return diff
}
