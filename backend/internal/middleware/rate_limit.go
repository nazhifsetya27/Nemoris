package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"nemoris/internal/utils"
)

const (
	rateLimitMaxRequests = 10
	rateLimitWindowSec   = 60
)

var (
	rateLimitStore   = make(map[string][]time.Time)
	rateLimitStoreMu sync.Mutex
)

// RateLimit limits requests per IP: max 10 per 60 seconds.
// Returns 429 with "rate limit exceeded" when exceeded.
// Uses in-memory map with mutex for concurrency safety.
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		now := time.Now()
		cutoff := now.Add(-rateLimitWindowSec * time.Second)

		rateLimitStoreMu.Lock()
		// 1. Remove timestamps older than 60 seconds
		ts := rateLimitStore[ip]
		pruned := ts[:0]
		for _, t := range ts {
			if t.After(cutoff) {
				pruned = append(pruned, t)
			}
		}

		// 2. Count remaining
		count := len(pruned)

		// 2b. Delete key if pruned empty (avoid storing empty slice)
		if len(pruned) == 0 {
			delete(rateLimitStore, ip)
		}

		// 3. Reject if count >= 10
		if count >= rateLimitMaxRequests {
			rateLimitStore[ip] = pruned
			rateLimitStoreMu.Unlock()
			utils.LogSystem("rate limit blocked: " + ip)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limit exceeded"))
			return
		}

		// 4. Append current timestamp if accepted
		pruned = append(pruned, now)
		rateLimitStore[ip] = pruned
		rateLimitStoreMu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// extractIP returns client IP from r.RemoteAddr, stripping port.
// Falls back to raw RemoteAddr if SplitHostPort fails.
func extractIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
