package middleware

import (
	"bytes"
	"github.com/dog-frame/common/metrics"
	"github.com/gin-gonic/gin"
	"time"
)

// PrometheusMiddleware Prometheus 监控中间件
func PrometheusMiddleware() gin.HandlerFunc {
	// 获取 Prometheus 指标实例
	prom := metrics.GetMetrics()

	return func(c *gin.Context) {
		// 获取请求信息
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// 记录请求开始时间
		startTime := time.Now()

		// 获取请求大小
		requestSize := 0
		if c.Request.ContentLength > 0 {
			requestSize = int(c.Request.ContentLength)
		}

		// 创建自定义响应写入器以捕获响应大小
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 获取响应状态码
		status := c.Writer.Status()

		// 计算请求处理时间
		duration := time.Since(startTime)

		// 获取响应大小
		responseSize := c.Writer.Size()

		// 记录指标
		prom.RecordMetrics(method, path, status, duration, requestSize, responseSize)

		// 记录错误
		if status >= 400 {
			var errorType string
			switch {
			case status >= 500:
				errorType = "server_error"
			case status >= 400:
				errorType = "client_error"
			}
			prom.RecordError(method, path, errorType)
		}
	}
}
