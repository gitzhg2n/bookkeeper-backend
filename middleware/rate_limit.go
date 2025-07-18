package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	visits map[string][]int64
	mu     sync.Mutex
}

func NewRateLimiter() *rateLimiter {
	return &rateLimiter{visits: make(map[string][]int64)}
}

func (rl *rateLimiter) Limit(windowMs int64, max int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			now := time.Now().UnixNano() / 1e6
			rl.mu.Lock()
			rl.visits[ip] = filterRecent(rl.visits[ip], now, windowMs)
			if len(rl.visits[ip]) >= max {
				rl.mu.Unlock()
				http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
				return
			}
			rl.visits[ip] = append(rl.visits[ip], now)
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}

func filterRecent(visits []int64, now int64, windowMs int64) []int64 {
	var out []int64
	for _, ts := range visits {
		if now-ts < windowMs {
			out = append(out, ts)
		}
	}
	return out
}