package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type ConcurrentDownloadLimiter struct {
	mu          sync.Mutex
	activeIPs   map[string]int
	maxPerIP    int
}

func NewConcurrentDownloadLimiter(maxPerIP int) *ConcurrentDownloadLimiter {
	return &ConcurrentDownloadLimiter{
		activeIPs: make(map[string]int),
		maxPerIP:  maxPerIP,
	}
}

func (l *ConcurrentDownloadLimiter) Acquire(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.activeIPs[ip] >= l.maxPerIP {
		return false
	}
	l.activeIPs[ip]++
	return true
}

func (l *ConcurrentDownloadLimiter) Release(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.activeIPs[ip] > 0 {
		l.activeIPs[ip]--
		if l.activeIPs[ip] == 0 {
			delete(l.activeIPs, ip)
		}
	}
}

func ConcurrentLimit(limiter *ConcurrentDownloadLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Acquire(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many concurrent downloads. Please wait for current download to finish.",
			})
			c.Abort()
			return
		}

		defer limiter.Release(ip)
		c.Next()
	}
}
