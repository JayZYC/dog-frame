package domainservice

import (
	"context"
	"github.com/dog-frame/common/errcode"
	"github.com/dog-frame/common/util"
	"github.com/dog-frame/dal/dao"
	"github.com/dog-frame/logic/do"
)

type DemoDomainSvc struct {
	ctx     context.Context
	DemoDao *dao.DemoDao
}

func NewDemoDomainSvc(ctx context.Context) *DemoDomainSvc {
	return &DemoDomainSvc{
		ctx:     ctx,
		DemoDao: dao.NewDemoDao(ctx),
	}
}

func (dds *DemoDomainSvc) GetDemos() ([]*do.Demo, error) {
	demos, err := dds.DemoDao.GetAllDemos()
	if err != nil {
		err = errcode.Wrap("query entity error", err)
		return nil, err
	}

	Demos := make([]*do.Demo, 0, len(demos))
	for _, demo := range demos {
		Demo := new(do.Demo)
		err := util.CopyProperties(Demo, demo)
		if err != nil {
			return nil, err
		}
		Demos = append(Demos, Demo)
	}

	return Demos, nil
}

func (dds *DemoDomainSvc) CreateDemo(Demo *do.Demo) (*do.Demo, error) {
	DemoModel, err := dds.DemoDao.CreateDemo(Demo)
	if err != nil {
		err = errcode.Wrap("创建Demo失败", err)
		return nil, err
	}

	err = util.CopyProperties(Demo, DemoModel)
	// 返回领域对象
	return Demo, err
}
