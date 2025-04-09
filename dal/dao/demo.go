package dao

import (
	"context"
	"github.com/dog-frame/common/util"
	"github.com/dog-frame/dal/model"
	"github.com/dog-frame/logic/do"
)

type DemoDao struct {
	ctx context.Context
}

func NewDemoDao(ctx context.Context) *DemoDao {
	return &DemoDao{ctx: ctx}
}

func (demo *DemoDao) GetAllDemos() (demos []*model.Demo, err error) {

	err = DB().WithContext(demo.ctx).Find(&demos).Error
	if err != nil {
		return nil, err
	}

	return demos, err
}

func (demo *DemoDao) CreateDemo(d *do.Demo) (*model.Demo, error) {
	m := new(model.Demo)
	err := util.CopyProperties(m, d)
	if err != nil {
		return nil, err
	}
	err = DB().WithContext(demo.ctx).
		Create(m).Error
	return m, err
}
