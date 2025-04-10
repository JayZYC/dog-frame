package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Model struct {
	ID        int64                 `gorm:"column:id;primaryKey;autoIncrement;comment:'自增主键ID'"`                       // 自增主键id
	IsDel     soft_delete.DeletedAt `gorm:"column:is_del;softDelete:flag;default:0;index;comment:'删除标识(0-未删除 1-已删除)'"` // 删除标识字段
	CreatedAt time.Time             `gorm:"column:created_at;autoCreateTime;index;comment:'创建时间'"`                     // 创建时间
	UpdatedAt time.Time             `gorm:"column:updated_at;autoUpdateTime;comment:'更新时间'"`                           // 更新时间
}
