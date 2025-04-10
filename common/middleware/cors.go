package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// CORSConfig 定义 CORS 配置
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig 返回默认的 CORS 配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400, // 24小时
	}
}

// CORS 返回一个处理跨域请求的中间件
func CORS() gin.HandlerFunc {
	return CORSWithConfig(DefaultCORSConfig())
}

// CORSWithConfig 返回一个使用自定义配置的 CORS 中间件
func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查请求的 Origin 是否在允许列表中
		allowOrigin := "*"
		if len(config.AllowOrigins) > 0 && config.AllowOrigins[0] != "*" {
			allowOrigin = ""
			for _, o := range config.AllowOrigins {
				if o == origin {
					allowOrigin = origin
					break
				}
				if strings.HasSuffix(o, "*") {
					prefix := strings.TrimSuffix(o, "*")
					if strings.HasPrefix(origin, prefix) {
						allowOrigin = origin
						break
					}
				}
			}
			if allowOrigin == "" {
				c.Next()
				return
			}
		}

		// 设置 CORS 响应头
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
		c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))

		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
