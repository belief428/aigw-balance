package model

import "time"

type Params struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement;comment:主键" json:"id" form:"id"` // 主键
	Key       string    `gorm:"column:key;type:varchar(60);uniqueIndex:idx_params_key;default:'';comment:唯一标识Key" json:"key"`
	Value     string    `gorm:"column:value;type:varchar(255);default:'';comment:具体数值" json:"value"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;comment:更新时间" json:"updated_at" form:"updated_at"`
}

func (*Params) TableName() string {
	return "params"
}
