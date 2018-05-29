package throttle

type RateLimiter interface {
	Acquire()
	TryAcquire() bool
}
