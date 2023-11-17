package baseModel

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ID struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
}

func (id *ID) BeforeCreate(tx *gorm.DB) (err error) {
	id.ID = uuid.New()
	return nil
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
