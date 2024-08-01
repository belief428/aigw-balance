package model

import "time"

type Archive struct {
	ID          uint64           `gorm:"column:id;primaryKey;autoIncrement;comment:主键" json:"id" form:"id"` // 主键
	GatewayCode string           `gorm:"column:gateway_code;type:varchar(20);default:'';comment:网关编号" json:"gateway_code"`
	Code        string           `gorm:"column:code;type:varchar(20);default:'';comment:设备编号" json:"code"`
	Attribute   ArchiveAttribute `gorm:"serializer:json;column:attribute;type:varchar(255);default:'';comment:属性信息" json:"attribute"`
	CreatedAt   time.Time        `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at" form:"created_at"`
	UpdatedAt   time.Time        `gorm:"column:updated_at;type:datetime;not null;comment:更新时间" json:"updated_at" form:"updated_at"`
}

type ArchiveAttribute struct {
	Regulate int     `json:"regulate"`
	Weight   float32 `json:"weight"`
}

func (*Archive) TableName() string {
	return "archive"
}
