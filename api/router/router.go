package router

import (
	"github.com/dog-frame/common/middleware"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) {
	// 初始化路由
	r := e.Group("")

	// 注册指标路由
	registerMetricsRoutes(r)

	r.Use(
		middleware.StartTrace,
		middleware.LogAccess,
		middleware.GinPanicRecovery,
		middleware.PrometheusMiddleware(), // 添加 Prometheus 监控中间件
	)

	// 注册测试路由
	registerTestRoutes(r)

	// 注册示例路由
	registerDemo(r)

}
