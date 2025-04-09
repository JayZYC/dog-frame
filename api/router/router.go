package router

import (
	"github.com/dog-frame/common/middleware"
	"github.com/gin-gonic/gin"
)

func Register(e *gin.Engine) {
	// 初始化路由
	e.Use(middleware.StartTrace, middleware.LogAccess, middleware.GinPanicRecovery)

	r := e.Group("")

	registerTestRoutes(r)
}
