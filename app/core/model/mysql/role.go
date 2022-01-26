package mysql

import "app/utility/entity"

type Role struct {
	entity.Entity
	Name  string `json:"name" xorm:"not null comment('角色名称') VARCHAR(100)"`
	Level int32  `json:"level" xorm:"not null comment('角色级别') INT(11)"`
}
