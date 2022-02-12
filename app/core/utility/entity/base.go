package entity

import (
	"time"
)

type Entity struct {
	DatabaseId   uint64        `json:"-" gorm:"column:_id;type:BIGINT UNSIGNED;primaryKey;autoIncrement:true;comment:数据库主键"`
	Service      ServiceEntity `gorm:"embedded"`
	DisplayOrder int           `json:"-" gorm:"type:INT;index:idx_display_order;comment:显示顺序"`
	Available    int           `json:"-" gorm:"default:1;type:INT;index:idx_available;comment:可用状态"`
	CreateTime   time.Time     `json:"-" gorm:"autoCreateTime;comment:创建时间"`
	UpdateTime   time.Time     `json:"-" gorm:"autoUpdateTime;comment:更新时间"`
	//DeleteTime   gorm.DeletedAt `json:"-" gorm:"index;comment:'删除时间标记'"`
	Comment string `json:"comment" gorm:"type:TEXT;comment:备注"`
}

type ServiceEntity struct {
	Id string `json:"id" gorm:"type:CHAR(20);uniqueIndex;not null;comment:业务主键"`
}
