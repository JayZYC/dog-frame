package controller

import (
	"github.com/dog-frame/api/request"
	"github.com/dog-frame/common/app"
	"github.com/dog-frame/common/errcode"
	"github.com/dog-frame/logic/appservice"
	"github.com/gin-gonic/gin"
)

// CreateDemo 创建示例数据
//
//	@Summary		创建示
//	@Summary		创建示例数据
//	@Description	创建一个新的示例数据
//	@Tags			示例
//	@Accept			json
//	@Produce		json
//	@Param			demo	body		request.DemoCreate	true	"示例数据"
//	@Success		200		{object}	app.Response		"成功"
//	@Failure		400		{object}	app.Response		"请求参数错误"
//	@Failure		500		{object}	app.Response		"服务器内部错误"
//	@Router			/demo [post]
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
