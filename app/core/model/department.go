package model

import "app/utility/entity"

type Department struct {
	entity.Entity
	Name     string `json:"name" xorm:"not null comment('部门/班组名称') VARCHAR(100)"`
	ParentId string `json:"parent_id" xorm:"comment('上级部门ID') CHAR(20)"`
}
