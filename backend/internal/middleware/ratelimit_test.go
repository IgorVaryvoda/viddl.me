package middleware

import (
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestIPRateLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(rate.Every(time.Minute/3), 3)

	ip := "192.168.1.1"
	l := limiter.GetLimiter(ip)
	if l == nil {
		t.Error("GetLimiter() returned nil")
	}

	l2 := limiter.GetLimiter(ip)
	if l != l2 {
		t.Error("GetLimiter() should return same limiter for same IP")
	}

	l3 := limiter.GetLimiter("192.168.1.2")
	if l == l3 {
		t.Error("GetLimiter() should return different limiter for different IP")
	}

	testIP := "10.0.0.1"
	testLimiter := limiter.GetLimiter(testIP)

	for i := 0; i < 3; i++ {
		if !testLimiter.Allow() {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	if testLimiter.Allow() {
		t.Error("4th request should be rate limited")
	}
}
