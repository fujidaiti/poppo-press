package httpserver

import (
	"net/http"
	"sync"
	"time"
)

type ipCounter struct {
	count int
	reset time.Time
}

type rateLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	byIP   map[string]*ipCounter
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{limit: limit, window: window, byIP: make(map[string]*ipCounter)}
}

func (rl *rateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		now := time.Now()
		rl.mu.Lock()
		c := rl.byIP[ip]
		if c == nil || now.After(c.reset) {
			c = &ipCounter{count: 0, reset: now.Add(rl.window)}
			rl.byIP[ip] = c
		}
		c.count++
		over := c.count > rl.limit
		rl.mu.Unlock()
		if over {
			writeError(w, http.StatusTooManyRequests, "rate_limited", "too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}
