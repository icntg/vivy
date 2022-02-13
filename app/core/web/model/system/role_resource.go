package system

import "app/core/utility/entity"

type RoleResource struct {
	entity.Entity
	RoleId     string `json:"role_id" gorm:"type:CHAR(20);not null;comment:角色ID"`
	ResourceId string `json:"resource_id" gorm:"type:CHAR(20);not null;comment:资源ID"`
}

func (r RoleResource) TableName() string {
	return "sys_role_resource"
}
