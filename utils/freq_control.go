package utils

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter      *rate.Limiter
	mu           sync.Mutex
	blocked      bool
	blockSeconds int64
}

func NewRateLimiter(r rate.Limit, b int, blockSeconds int64) *RateLimiter {
	return &RateLimiter{
		limiter:      rate.NewLimiter(r, b),
		blockSeconds: blockSeconds,
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.blocked {
		return false
	}

	if !r.limiter.Allow() {
		r.blocked = true
		go func() {
			time.Sleep(time.Duration(r.blockSeconds) * time.Second)
			r.mu.Lock()
			r.blocked = false
			r.mu.Unlock()
		}()
		return false
	}
	return true
}
