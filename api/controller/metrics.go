package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// PrometheusHandler 处理 Prometheus 指标请求
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// HealthCheck 健康检查接口
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// ReadinessCheck 就绪检查接口
func ReadinessCheck(c *gin.Context) {
	// 这里可以添加更多的就绪检查逻辑，例如检查数据库连接等
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
