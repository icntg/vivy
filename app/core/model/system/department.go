package system

import "app/core/utility/entity"

type Department struct {
	entity.Entity
	Name     string `json:"name" gorm:"type:VARCHAR(100);not null;comment:'部门/班组名称'"`
	ParentId string `json:"parent_id" gorm:"type:CHAR(20);comment:'上级部门ID'"`
}
