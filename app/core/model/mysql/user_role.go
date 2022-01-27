package mysql

import "app/core/utility/entity"

type UserRole struct {
	entity.Entity
	UserId string `json:"user_id" gorm:"type:CHAR(20);not null;comment:'用户ID'"`
	RoleId string `json:"role_id" gorm:"type:CHAR(20);not null;comment:'角色ID'"`
}
