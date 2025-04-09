package controller

import (
	"github.com/dog-frame/api/request"
	"github.com/dog-frame/common/app"
	"github.com/dog-frame/common/errcode"
	"github.com/dog-frame/logic/appservice"
	"github.com/gin-gonic/gin"
)

func CreateDemo(c *gin.Context) {
	// 解析参数
	req := new(request.DemoCreate)

	if err := c.ShouldBindJSON(&req); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	svc := appservice.NewDemoAppSvc(c)

	reply, err := svc.CreateDemo(req)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	app.NewResponse(c).Success(reply)

}
