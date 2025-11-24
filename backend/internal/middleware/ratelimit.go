package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
	lastSeen map[string]time.Time
	rate     rate.Limit
	burst    int
}

func NewIPRateLimiter(r rate.Limit, burst int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		lastSeen: make(map[string]time.Time),
		rate:     r,
		burst:    burst,
	}
	go limiter.cleanupOldLimiters()
	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.limiters[ip]
	i.mu.RUnlock()

	if exists {
		i.mu.Lock()
		i.lastSeen[ip] = time.Now()
		i.mu.Unlock()
		return limiter
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	if limiter, exists := i.limiters[ip]; exists {
		i.lastSeen[ip] = time.Now()
		return limiter
	}

	limiter = rate.NewLimiter(i.rate, i.burst)
	i.limiters[ip] = limiter
	i.lastSeen[ip] = time.Now()
	return limiter
}

func (i *IPRateLimiter) cleanupOldLimiters() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		i.mu.Lock()
		now := time.Now()
		for ip, lastSeen := range i.lastSeen {
			if now.Sub(lastSeen) > 30*time.Minute {
				delete(i.limiters, ip)
				delete(i.lastSeen, ip)
				log.Printf("INFO: Cleaned up rate limiter for IP: %s", ip)
			}
		}
		i.mu.Unlock()
	}
}

func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
