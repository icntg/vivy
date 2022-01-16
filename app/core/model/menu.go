package model

import (
	"app/utility/entity"
)

type Menu struct {
	entity.Entity
	Name        string `json:"name" xorm:"not null default '' comment('名称') VARCHAR(100)"`
	Path        string `json:"path" xorm:"not null default '' comment('路径') index VARCHAR(50)"`
	Component   string `json:"component" xorm:"not null default '' comment('组件') VARCHAR(100)"`
	Redirect    string `json:"redirect" xorm:"not null default '' comment('重定向') VARCHAR(200)"`
	Url         string `json:"url" xorm:"not null default '' comment('url') VARCHAR(200)"`
	MetaTitle   string `json:"meta_title" xorm:"not null default '' comment('meta标题') VARCHAR(50)"`
	MetaIcon    string `json:"meta_icon" xorm:"not null default '' comment('meta icon') VARCHAR(50)"`
	MetaNocache int    `json:"meta_nocache" xorm:"not null default 0 comment('是否缓存（1:是 0:否）') TINYINT(4)"`
	AlwaysShow  int    `json:"always_show" xorm:"not null default 0 comment('是否总是显示（1:是0：否）') TINYINT(4)"`
	MetaAffix   int    `json:"meta_affix" xorm:"not null default 0 comment('是否加固（1:是0：否）') TINYINT(4)"`
	Type        int    `json:"type" xorm:"not null default 2 comment('类型(1:固定,2:权限配置,3特殊)') TINYINT(4)"`
	Hidden      int    `json:"hidden" xorm:"not null default 0 comment('是否隐藏（0否1是）') TINYINT(4)"`
	ParentId    string `json:"parent_id" xorm:"not null default 0 comment('父ID') index(idx_list) INT(11)"`
	Sort        int    `json:"sort" xorm:"not null default 0 comment('排序') index(idx_list) INT(11)"`
	Level       int    `json:"level" xorm:"not null default 0 comment('层级') TINYINT(4)"`
}

type MenuRoute struct {
	entity.Entity
	Url string `json:"url" xorm:"not null default '' comment('url') VARCHAR(200)"`
}
