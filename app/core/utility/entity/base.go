package entity

import (
	"gorm.io/gorm"
	"time"
)

type Entity struct {
	DatabaseId   uint64         `json:"db_id" gorm:"type:BIGINT UNSIGNED;primaryKey;autoIncrement:true;comment:'数据库主键'"`
	Id           string         `json:"id" gorm:"type:CHAR(20);uniqueIndex;not null;comment:'业务主键'"`
	DisplayOrder int            `json:"display_order" gorm:"type:INT;index:idx_display_order;comment:'显示顺序'"`
	Available    int            `json:"available" gorm:"default:1;type:INT;index:idx_available;comment:'可用状态'"`
	CreateTime   time.Time      `json:"create_time" gorm:"autoCreateTime;comment:'创建时间'"`
	UpdateTime   time.Time      `json:"update_time" gorm:"autoUpdateTime;comment:'更新时间'"`
	DeleteTime   gorm.DeletedAt `json:"delete_time" gorm:"index;comment:'删除时间标记'"`
	Comment      string         `json:"comment" gorm:"type:TEXT;comment:'备注'"`
}
