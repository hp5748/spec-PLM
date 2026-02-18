package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"plm/pkg/response"
)

// RateLimiter 限流器
type RateLimiter struct {
	records map[string]*record
	mu      sync.RWMutex
}

type record struct {
	count     int
	expiresAt time.Time
}

// NewRateLimiter 创建限流器
func NewRateLimiter() *RateLimiter {
	limiter := &RateLimiter{
		records: make(map[string]*record),
	}
	// 启动清理协程
	go limiter.cleanup()
	return limiter
}

// cleanup 定期清理过期记录
func (r *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		r.mu.Lock()
		now := time.Now()
		for key, rec := range r.records {
			if rec.expiresAt.Before(now) {
				delete(r.records, key)
			}
		}
		r.mu.Unlock()
	}
}

// Allow 检查是否允许请求
func (r *RateLimiter) Allow(key string, limit int, window time.Duration) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	rec, exists := r.records[key]

	if !exists || rec.expiresAt.Before(now) {
		// 创建新记录
		r.records[key] = &record{
			count:     1,
			expiresAt: now.Add(window),
		}
		return true
	}

	// 检查是否超限
	if rec.count >= limit {
		return false
	}

	rec.count++
	return true
}

// RateLimit 限流中间件
// limit: 时间窗口内允许的最大请求数
// window: 时间窗口
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter()

	return func(c *gin.Context) {
		// 使用IP作为限流Key
		key := c.ClientIP()

		if !limiter.Allow(key, limit, window) {
			response.Error(c, 42900, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUser 按用户限流
func RateLimitByUser(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter()

	return func(c *gin.Context) {
		var key string

		// 优先使用用户ID
		if userID, exists := c.Get("userID"); exists {
			key = "user:" + string(rune(userID.(uint)))
		} else {
			// 未登录使用IP
			key = "ip:" + c.ClientIP()
		}

		if !limiter.Allow(key, limit, window) {
			response.Error(c, 42900, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
