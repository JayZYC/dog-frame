package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/dog-frame/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// CSRFConfig CSRF 防护配置
type CSRFConfig struct {
	// CSRF Token 的密钥
	Secret string
	// Cookie 名称
	CookieName string
	// Header 名称
	HeaderName string
	// 表单字段名称
	FormFieldName string
	// Cookie 过期时间（秒）
	CookieMaxAge int
	// 是否使用安全 Cookie
	CookieSecure bool
	// 是否使用 HttpOnly Cookie
	CookieHTTPOnly bool
	// Cookie 路径
	CookiePath string
	// Cookie 域名
	CookieDomain string
	// 排除的路径（不需要 CSRF 保护）
	ExcludedPaths []string
	// 排除的方法（不需要 CSRF 保护）
	ExcludedMethods []string
}

// DefaultCSRFConfig 返回默认的 CSRF 配置
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		Secret:         "dog-frame-csrf-secret", // 在生产环境中应该使用更安全的密钥
		CookieName:     "csrf_token",
		HeaderName:     "X-CSRF-Token",
		FormFieldName:  "_csrf",
		CookieMaxAge:   86400, // 24小时
		CookieSecure:   false, // 在生产环境中应该设置为 true
		CookieHTTPOnly: true,
		CookiePath:     "/",
		CookieDomain:   "",
		ExcludedPaths:  []string{"/api/health", "/api/metrics"},
		ExcludedMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
			http.MethodTrace,
		},
	}
}

// CSRF 返回一个 CSRF 防护中间件
func CSRF() gin.HandlerFunc {
	return CSRFWithConfig(DefaultCSRFConfig())
}

// CSRFWithConfig 返回一个使用自定义配置的 CSRF 防护中间件
func CSRFWithConfig(config CSRFConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否是排除的路径
		path := c.Request.URL.Path
		for _, excludedPath := range config.ExcludedPaths {
			if strings.HasPrefix(path, excludedPath) {
				c.Next()
				return
			}
		}

		// 检查是否是排除的方法
		method := c.Request.Method
		for _, excludedMethod := range config.ExcludedMethods {
			if method == excludedMethod {
				c.Next()
				return
			}
		}

		// 获取或生成 CSRF Token
		token, err := getOrGenerateCSRFToken(c, config)
		if err != nil {
			logger.Error(c, "Failed to get or generate CSRF token", "err", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 验证 CSRF Token
		if !validateCSRFToken(c, token, config) {
			logger.Warn(c, "CSRF token validation failed", "path", path, "method", method)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// getOrGenerateCSRFToken 获取或生成 CSRF Token
func getOrGenerateCSRFToken(c *gin.Context, config CSRFConfig) (string, error) {
	// 从 Cookie 中获取 Token
	cookie, err := c.Cookie(config.CookieName)
	if err == nil && cookie != "" {
		return cookie, nil
	}

	// 生成新的 Token
	token, err := generateCSRFToken(config.Secret)
	if err != nil {
		return "", err
	}

	// 设置 Cookie
	c.SetCookie(
		config.CookieName,
		token,
		config.CookieMaxAge,
		config.CookiePath,
		config.CookieDomain,
		config.CookieSecure,
		config.CookieHTTPOnly,
	)

	return token, nil
}

// generateCSRFToken 生成 CSRF Token
func generateCSRFToken(secret string) (string, error) {
	// 使用当前时间戳作为基础
	timestamp := time.Now().UnixNano()

	// 创建 HMAC
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(string(rune(timestamp))))
	if err != nil {
		return "", err
	}

	// 生成 Token
	token := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return token, nil
}

// validateCSRFToken 验证 CSRF Token
func validateCSRFToken(c *gin.Context, expectedToken string, config CSRFConfig) bool {
	// 从 Header 中获取 Token
	headerToken := c.GetHeader(config.HeaderName)
	if headerToken != "" && hmac.Equal([]byte(headerToken), []byte(expectedToken)) {
		return true
	}

	// 从表单中获取 Token
	formToken := c.PostForm(config.FormFieldName)
	if formToken != "" && hmac.Equal([]byte(formToken), []byte(expectedToken)) {
		return true
	}

	// 如果都没有找到有效的 Token，则验证失败
	return false
}

// GetCSRFToken 获取当前请求的 CSRF Token
func GetCSRFToken(c *gin.Context) (string, error) {
	token, err := c.Cookie(DefaultCSRFConfig().CookieName)
	if err != nil || token == "" {
		return "", errors.New("CSRF token not found")
	}
	return token, nil
}
