package middleware

import (
	"github.com/dog-frame/common/errcode"
	"github.com/dog-frame/common/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

// RateLimiterConfig 速率限制配置
type RateLimiterConfig struct {
	// 每秒允许的请求数
	Rate float64
	// 令牌桶容量
	Burst int
	// 限流器过期时间（秒）
	ExpirationTime int64
	// 限流器清理间隔（秒）
	CleanupInterval int64
	// 限流键生成函数
	KeyFunc func(*gin.Context) string
	// 排除的路径
	ExcludedPaths []string
}

// DefaultRateLimiterConfig 返回默认的速率限制配置
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Rate:            10,   // 每秒10个请求
		Burst:           20,   // 最多允许20个请求同时处理
		ExpirationTime:  3600, // 1小时后过期
		CleanupInterval: 60,   // 每分钟清理一次过期的限流器
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP() // 默认使用客户端IP作为限流键
		},
		ExcludedPaths: []string{"/api/health", "/api/metrics"},
	}
}

// rateLimiter 速率限制器
type rateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// rateLimiterStore 存储所有的速率限制器
type rateLimiterStore struct {
	limiters map[string]*rateLimiter
	mu       sync.RWMutex
	config   RateLimiterConfig
}

// 全局速率限制器存储
var limiterStore *rateLimiterStore
var limiterOnce sync.Once

// getLimiterStore 获取或创建速率限制器存储
func getLimiterStore(config RateLimiterConfig) *rateLimiterStore {
	limiterOnce.Do(func() {
		limiterStore = &rateLimiterStore{
			limiters: make(map[string]*rateLimiter),
			config:   config,
		}

		// 启动清理过期限流器的协程
		go limiterStore.cleanup()
	})
	return limiterStore
}

// cleanup 清理过期的限流器
func (s *rateLimiterStore) cleanup() {
	ticker := time.NewTicker(time.Duration(s.config.CleanupInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for key, limiter := range s.limiters {
			if time.Since(limiter.lastSeen).Seconds() > float64(s.config.ExpirationTime) {
				delete(s.limiters, key)
			}
		}
		s.mu.Unlock()
	}
}

// getLimiter 获取指定键的限流器
func (s *rateLimiterStore) getLimiter(key string) *rate.Limiter {
	s.mu.RLock()
	limiter, exists := s.limiters[key]
	s.mu.RUnlock()

	if !exists {
		s.mu.Lock()
		defer s.mu.Unlock()

		// 再次检查，避免并发创建
		limiter, exists = s.limiters[key]
		if !exists {
			limiter = &rateLimiter{
				limiter:  rate.NewLimiter(rate.Limit(s.config.Rate), s.config.Burst),
				lastSeen: time.Now(),
			}
			s.limiters[key] = limiter
		}
	}

	// 更新最后访问时间
	limiter.lastSeen = time.Now()
	return limiter.limiter
}

// RateLimit 返回一个速率限制中间件
func RateLimit() gin.HandlerFunc {
	return RateLimitWithConfig(DefaultRateLimiterConfig())
}

// RateLimitWithConfig 返回一个使用自定义配置的速率限制中间件
func RateLimitWithConfig(config RateLimiterConfig) gin.HandlerFunc {
	store := getLimiterStore(config)

	return func(c *gin.Context) {
		// 检查是否是排除的路径
		path := c.Request.URL.Path
		for _, excludedPath := range config.ExcludedPaths {
			if path == excludedPath {
				c.Next()
				return
			}
		}

		// 获取限流键
		key := config.KeyFunc(c)

		// 获取限流器
		limiter := store.getLimiter(key)

		// 尝试获取令牌
		if !limiter.Allow() {
			logger.Warn(c, "Rate limit exceeded", "ip", c.ClientIP(), "path", path)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": errcode.ErrTooManyRequests.Code(),
				"msg":  errcode.ErrTooManyRequests.Msg(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
