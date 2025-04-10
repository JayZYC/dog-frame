package middleware

import (
	"bytes"
	"github.com/dog-frame/common/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// XSSConfig XSS 防护配置
type XSSConfig struct {
	// 是否对请求体进行过滤
	FilterRequestBody bool
	// 是否对响应体进行过滤
	FilterResponseBody bool
	// 是否对 URL 参数进行过滤
	FilterURLParams bool
	// 是否对 Header 进行过滤
	FilterHeaders bool
	// 需要过滤的 Header 名称列表
	FilterHeaderNames []string
}

// DefaultXSSConfig 返回默认的 XSS 配置
func DefaultXSSConfig() XSSConfig {
	return XSSConfig{
		FilterRequestBody:  true,
		FilterResponseBody: false, // 默认不过滤响应体，因为可能会影响正常的 HTML 内容
		FilterURLParams:    true,
		FilterHeaders:      true,
		FilterHeaderNames:  []string{"Cookie", "Authorization"},
	}
}

// XSS 返回一个 XSS 防护中间件
func XSS() gin.HandlerFunc {
	return XSSWithConfig(DefaultXSSConfig())
}

// XSSWithConfig 返回一个使用自定义配置的 XSS 防护中间件
func XSSWithConfig(config XSSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 过滤 URL 参数
		if config.FilterURLParams {
			for key, values := range c.Request.URL.Query() {
				for i, value := range values {
					c.Request.URL.Query()[key][i] = sanitizeString(value)
				}
			}
		}

		// 过滤 Header
		if config.FilterHeaders && len(config.FilterHeaderNames) > 0 {
			for _, headerName := range config.FilterHeaderNames {
				if value := c.Request.Header.Get(headerName); value != "" {
					c.Request.Header.Set(headerName, sanitizeString(value))
				}
			}
		}

		// 过滤请求体
		if config.FilterRequestBody {
			contentType := c.GetHeader("Content-Type")
			// 只处理 JSON 和表单数据，不处理文件上传等二进制内容
			if strings.Contains(contentType, "application/json") ||
				strings.Contains(contentType, "application/x-www-form-urlencoded") ||
				strings.Contains(contentType, "text/plain") {

				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err != nil {
					logger.Error(c, "XSS middleware read request body error", "err", err)
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}

				// 关闭原始请求体
				_ = c.Request.Body.Close()

				// 过滤请求体内容
				sanitizedBody := sanitizeString(string(bodyBytes))

				// 重新设置请求体
				c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(sanitizedBody)))
				c.Request.ContentLength = int64(len(sanitizedBody))
			}
		}

		// 如果需要过滤响应体，则需要包装 ResponseWriter
		if config.FilterResponseBody {
			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw

			c.Next()

			// 获取原始响应
			responseBody := blw.body.String()

			// 过滤响应内容
			sanitizedResponse := sanitizeString(responseBody)

			// 重写响应
			c.Writer.Header().Set("Content-Length", string(rune(len(sanitizedResponse))))
			c.Writer.WriteHeader(c.Writer.Status())
			_, _ = c.Writer.Write([]byte(sanitizedResponse))
			return
		}

		c.Next()
	}
}

// sanitizeString 对字符串进行 XSS 过滤
func sanitizeString(input string) string {
	// 替换常见的 XSS 攻击向量
	patterns := []struct {
		regex       *regexp.Regexp
		replacement string
	}{
		// 脚本标签
		{regexp.MustCompile(`(?i)<script[^>]*>[\s\S]*?</script[^>]*>`), ""},
		// 内联事件
		{regexp.MustCompile(`(?i)(on\w+)=["']?(.*?)["']?`), ""},
		// JavaScript URLs
		{regexp.MustCompile(`(?i)javascript:[^\s]`), ""},
		// CSS 表达式
		{regexp.MustCompile(`(?i)expression[\s]*\(`), ""},
		// CSS 属性
		{regexp.MustCompile(`(?i)behavior[\s]*:`), ""},
		// META 刷新
		{regexp.MustCompile(`(?i)<meta[^>]*refresh[^>]*>`), ""},
		// 内联框架
		{regexp.MustCompile(`(?i)<iframe[^>]*>[\s\S]*?</iframe[^>]*>`), ""},
		// 对象标签
		{regexp.MustCompile(`(?i)<object[^>]*>[\s\S]*?</object[^>]*>`), ""},
		// 嵌入标签
		{regexp.MustCompile(`(?i)<embed[^>]*>[\s\S]*?</embed[^>]*>`), ""},
		// 表单标签
		{regexp.MustCompile(`(?i)<form[^>]*>[\s\S]*?</form[^>]*>`), ""},
		// 常见的 XSS 攻击字符串
		{regexp.MustCompile(`(?i)<[^>]*\s+src\s*=\s*["']?[^"'>]*["']?[^>]*>`), ""},
		// 替换尖括号
		{regexp.MustCompile(`(?i)<`), "&lt;"},
		{regexp.MustCompile(`(?i)>`), "&gt;"},
	}

	result := input
	for _, pattern := range patterns {
		result = pattern.regex.ReplaceAllString(result, pattern.replacement)
	}

	return result
}
