package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter stores a rate limiter per IP address.
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       *sync.RWMutex
	r        rate.Limit
	b        int
}

// NewIPRateLimiter creates a new rate limiter.
// r: requests per second (e.g., 10)
// b: burst size (max tokens)
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		mu:       &sync.RWMutex{},
		r:        r,
		b:        b,
	}
}

// getLimiter returns the limiter for the given IP, creating it if necessary.
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.limiters[ip] = limiter
	}
	return limiter
}

// RateLimitMiddleware returns a Gin middleware that limits requests per IP.
func RateLimitMiddleware(r rate.Limit, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(r, burst)

	return func(c *gin.Context) {
		// Get client IP (respect X-Forwarded-For if behind proxy)
		clientIP := c.ClientIP()
		l := limiter.getLimiter(clientIP)
		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, please slow down",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
