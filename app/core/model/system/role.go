package system

import "app/core/utility/entity"

type Role struct {
	entity.Entity
	ParentId string `json:"parent_id" gorm:"type:CHAR(20);comment:'上级角色ID'"`
	Name     string `json:"name" gorm:"type:VARCHAR(100);not null;comment:'角色名称'"`
	Level    int32  `json:"level" gorm:"type:INT;not null;comment:'角色级别'"`
}
