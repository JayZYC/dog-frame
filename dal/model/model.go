package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Model struct {
	ID        int64                 `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 自增主键id
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag"`                                      // 删除标识字段(0-未删除 1-已删除)
	CreatedAt time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}
