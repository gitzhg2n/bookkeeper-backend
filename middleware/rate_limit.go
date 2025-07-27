package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type rateLimiter struct {
	visits    map[string][]int64
	mu        sync.RWMutex
	lastClean time.Time
}

func NewRateLimiter() *rateLimiter {
	return &rateLimiter{
		visits:    make(map[string][]int64),
		lastClean: time.Now(),
	}
}

func (rl *rateLimiter) Limit(windowMs int64, max int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get real IP, handle X-Forwarded-For and X-Real-IP headers
			ip := rl.getRealIP(r)
			now := time.Now().UnixNano() / 1e6
			
			rl.mu.Lock()
			// Clean up old entries periodically (every 5 minutes)
			if now-rl.lastClean.UnixNano()/1e6 > 300000 {
				rl.cleanup(now, windowMs)
				rl.lastClean = time.Now()
			}
			
			rl.visits[ip] = filterRecent(rl.visits[ip], now, windowMs)
			if len(rl.visits[ip]) >= max {
				rl.mu.Unlock()
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", max))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", fmt.Sprintf("%.0f", float64(windowMs)/1000))
				http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
				return
			}
			rl.visits[ip] = append(rl.visits[ip], now)
			remaining := max - len(rl.visits[ip])
			rl.mu.Unlock()
			
			// Add rate limit headers
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", max))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			
			next.ServeHTTP(w, r)
		})
	}
}

// getRealIP extracts the real IP address from the request
func (rl *rateLimiter) getRealIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take the first IP from the comma-separated list
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}
	
	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" && net.ParseIP(xri) != nil {
		return xri
	}
	
	// Fall back to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// cleanup removes expired entries from all IPs
func (rl *rateLimiter) cleanup(now int64, windowMs int64) {
	for ip, visits := range rl.visits {
		filtered := filterRecent(visits, now, windowMs)
		if len(filtered) == 0 {
			delete(rl.visits, ip)
		} else {
			rl.visits[ip] = filtered
		}
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