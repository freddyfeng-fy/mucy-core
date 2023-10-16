package baseModel

import (
	"gorm.io/gorm"
	"time"
)

// 自增ID主键
type ID struct {
	ID int64 `json:"id" gorm:"primaryKey;SERIAL"`
}

// 创建、更新时间
type Timestamps struct {
	CreateAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createAt"`
	CreateBy int64     `json:"createBy"`
	UpdateAt time.Time `json:"updateAt"`
	UpdateBy int64     `json:"updateBy"`
}

// 软删除
type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type LabelValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
