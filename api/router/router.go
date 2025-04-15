package router

import (
	"github.com/dog-frame/common/enum"
	"github.com/dog-frame/common/middleware"
	"github.com/dog-frame/config"
	"github.com/gin-gonic/gin"

	"github.com/dog-frame/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register(e *gin.Engine) {
	// 初始化路由
	r := e.Group("")

	// 注册指标路由
	registerMetricsRoutes(r)

	// 生产模式关闭Swagger
	if config.App.Env != enum.ModeProd {
		r.GET(docs.SwaggerInfo.BasePath+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

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
