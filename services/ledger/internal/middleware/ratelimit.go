package middleware

import (
	"net/http"
	"sync"
	"time"
)

type bucket struct {
	tokens float64
	last   time.Time
}

// RateLimit is a per-user token bucket limiter.
func RateLimit(rps int) func(http.Handler) http.Handler {
	var mu sync.Mutex
	buckets := map[string]*bucket{}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := UserID(r.Context())
			if key == "" {
				key = r.RemoteAddr
			}
			mu.Lock()
			b, ok := buckets[key]
			now := time.Now()
			if !ok {
				b = &bucket{tokens: float64(rps), last: now}
				buckets[key] = b
			}
			b.tokens += now.Sub(b.last).Seconds() * float64(rps)
			if b.tokens > float64(rps) {
				b.tokens = float64(rps)
			}
			b.last = now
			allowed := b.tokens >= 1
			if allowed {
				b.tokens--
			}
			mu.Unlock()
			if !allowed {
				w.Header().Set("Retry-After", "1")
				http.Error(w, `{"code":"rate_limited","message":"too many requests"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
