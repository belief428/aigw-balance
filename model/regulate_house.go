package model

import "time"

type RegulateHouse struct {
	GatewayCode string           `gorm:"column:gateway_code;type:varchar(20);default:'';comment:网关编号" json:"gateway_code"`
	ArchiveCode string           `gorm:"column:archive_code;type:varchar(20);default:'';comment:设备编号" json:"archive_code"`
	ArchiveName string           `gorm:"column:archive_name;type:varchar(100);default:'';comment:设备名称区域" json:"archive_name"`
	Params      []RegulateParams `gorm:"serializer:json;column:params;type:varchar(255);default:'';comment:参数信息" json:"params"`
	PrevDeg     uint8            `gorm:"column:prev_deg;type:tinyint(1);default:0;comment:调控前开度" json:"prev_deg"`
	NextDeg     uint8            `gorm:"column:next_deg;type:tinyint(1);default:0;comment:调控后开度" json:"next_deg"`
	Status      int              `gorm:"column:status;type:tinyint(1);default:0;comment:状态" json:"status"`
	Remark      string           `gorm:"column:remark;type:varchar(255);default:'';comment:备注信息" json:"remark"`
	Date        time.Time        `gorm:"column:date;type:datetime;not null;comment:调控时间" json:"date" form:"date"`
}

func (*RegulateHouse) TableName() string {
	return "regulate_house"
}
