package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestPanicLog(c *gin.Context) {
	var a map[string]string
	a["k"] = "v"
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   a,
	})
	return
}
