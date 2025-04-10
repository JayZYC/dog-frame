package router

import (
	"github.com/dog-frame/api/controller"
	"github.com/gin-gonic/gin"
)

func registerMetricsRoutes(r *gin.RouterGroup) {
	// 健康检查和就绪检查接口
	r.GET("/health", controller.HealthCheck)
	r.GET("/ready", controller.ReadinessCheck)

	// Prometheus 指标接口
	r.GET("/metrics", controller.PrometheusHandler())
}
