package mysql

import "app/utility/entity"

type UserRole struct {
	entity.Entity
	UserId string `json:"user_id" xorm:"not null comment('用户ID') CHAR(20)"`
	RoleId string `json:"role_id" xorm:"not null comment('角色ID') CHAR(20)"`
}
