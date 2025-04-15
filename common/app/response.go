package app

import (
	"github.com/dog-frame/common/errcode"
	"github.com/dog-frame/common/logger"
	"github.com/gin-gonic/gin"
)

type Response struct {
	ctx        *gin.Context
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	RequestId  string      `json:"request_id"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func NewResponse(c *gin.Context) *Response {
	return &Response{ctx: c}
}

// SetPagination 设置Response的分页信息
func (r *Response) SetPagination(pagination *Pagination) *Response {
	r.Pagination = pagination
	return r
}

func (r *Response) Success(data interface{}) {
	r.Code = errcode.Success.Code()
	r.Msg = errcode.Success.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("traceid"); exists {
		val, _ := r.ctx.Get("traceid")
		requestId = val.(string)
	}
	r.RequestId = requestId
	r.Data = data

	r.ctx.JSON(errcode.Success.HttpStatusCode(), r)
}

func (r *Response) SuccessOk() {
	r.Success("")
}

func (r *Response) Error(err *errcode.AppError) {
	r.Code = err.Code()
	r.Msg = err.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("traceid"); exists {
		val, _ := r.ctx.Get("traceid")
		requestId = val.(string)
	}
	r.RequestId = requestId
	// 兜底记一条响应错误, 项目自定义的AppError中有错误链条, 方便出错后排查问题
	logger.Error(r.ctx, "api_response_error", "err", err)
	r.ctx.JSON(err.HttpStatusCode(), r)
}
