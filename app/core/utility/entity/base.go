package entity

import (
	"gorm.io/gorm"
	"time"
)

type Entity struct {
	DatabaseId   uint64    `json:"db_id" xorm:"not null pk autoincr comment('数据库主键') INT(11)"`
	Id           string    `json:"id" xorm:"not null comment('业务主键') CHAR(20) unique"`
	DisplayOrder int       `json:"display_order" xorm:"comment('显示顺序') INT(11) index"`
	Available    int       `json:"available" xorm:"not null default 1 comment('状态（0 停止1启动）') TINYINT(4) index"`
	CreateTime   time.Time `json:"create_time" xorm:"not null comment('创建时间') TIMESTAMP"`
	UpdateTime   time.Time `json:"update_time" xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	DeleteTime   gorm.DeletedAt
	Comment      string `json:"comment" xorm:"comment('备注') TEXT"`
}

func EntityToMap(entity interface{}, fields ...string) map[string]interface{} {
	return nil
}

func EntitiesToListMap(entities []interface{}, fields ...string) []map[string]interface{} {
	return nil
}
