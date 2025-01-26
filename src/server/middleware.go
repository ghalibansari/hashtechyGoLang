package server

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter manages request rate limiting per IP address to prevent abuse
type RateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewRateLimiter creates a new rate limiter with specified rate and burst capacity
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// RateLimit middleware checks and enforces request rate limits for each unique IP
func (l *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		l.mu.Lock()
		limiter, exists := l.ips[ip]
		if !exists {
			limiter = rate.NewLimiter(l.r, l.b)
			l.ips[ip] = limiter
		}
		l.mu.Unlock()

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
