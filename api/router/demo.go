package router

import (
	"github.com/dog-frame/api/controller"
	"github.com/gin-gonic/gin"
)

func registerDemo(r *gin.RouterGroup) {
	g := r.Group("/demo")

	g.POST("/demo", controller.CreateDemo)
}
