package appservice

import (
	"context"
	"github.com/dog-frame/api/request"
	"github.com/dog-frame/api/response"
	"github.com/dog-frame/common/errcode"
	"github.com/dog-frame/common/logger"
	"github.com/dog-frame/common/util"
	"github.com/dog-frame/dal/cache"
	"github.com/dog-frame/logic/do"
	"github.com/dog-frame/logic/domainservice"
)

type DemoAppSvc struct {
	ctx           context.Context
	demoDomainSvc *domainservice.DemoDomainSvc
}

func NewDemoAppSvc(ctx context.Context) *DemoAppSvc {
	return &DemoAppSvc{
		ctx:           ctx,
		demoDomainSvc: domainservice.NewDemoDomainSvc(ctx),
	}
}

func (das *DemoAppSvc) CreateDemo(orderRequest *request.DemoCreate) (*response.Demo, error) {
	DemoDo := new(do.Demo)
	err := util.CopyProperties(DemoDo, orderRequest)
	if err != nil {
		errcode.Wrap("请求转换成DemoDo失败", err)
		return nil, err
	}
	DemoDo, err = das.demoDomainSvc.CreateDemo(DemoDo)
	if err != nil {
		return nil, err
	}

	// 设置缓存和读取, 测试项目中缓存的使用, 没有其他任何意义
	err = cache.SetDemo(das.ctx, DemoDo)
	if err != nil {
		return nil, err
	}
	cacheData, _ := cache.GetDemo(das.ctx, DemoDo.Name)
	logger.Info(das.ctx, "redis data", "data", cacheData)

	replyDemo := new(response.Demo)
	err = util.CopyProperties(replyDemo, DemoDo)
	if err != nil {
		errcode.Wrap("DemoDo转换成replyDemo失败", err)
		return nil, err
	}

	return replyDemo, err
}
