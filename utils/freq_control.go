package utils

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	requests      int
	lastCheckTime time.Time
	mutex         sync.Mutex
}

var ipRateLimiters = make(map[string]*rateLimiter)
var globalMutex sync.Mutex

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		now := time.Now()

		globalMutex.Lock()
		limiter, exists := ipRateLimiters[ip]
		if !exists {
			limiter = &rateLimiter{
				requests:      0,
				lastCheckTime: now,
			}
			ipRateLimiters[ip] = limiter
		}
		globalMutex.Unlock()

		limiter.mutex.Lock()
		defer limiter.mutex.Unlock()

		elapsed := now.Sub(limiter.lastCheckTime)
		limiter.lastCheckTime = now

		limiter.requests -= int(elapsed.Seconds() * 10)
		if limiter.requests < 0 {
			limiter.requests = 0
		}

		if limiter.requests < 10 {
			limiter.requests++
			next.ServeHTTP(w, r)
		} else {
			Response(w, r, http.StatusTooManyRequests, "too many requests", nil)
		}
	})
}
