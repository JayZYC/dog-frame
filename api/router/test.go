package router

import (
	"github.com/dog-frame/api/controller"
	"github.com/gin-gonic/gin"
)

func registerTestRoutes(r *gin.RouterGroup) {

	g := r.Group("/test")

	g.GET("/panic-log-test", controller.TestPanicLog)

}
