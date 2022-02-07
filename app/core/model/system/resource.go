package system

import "app/core/utility/entity"

// Resource 资源基类
type Resource struct {
	entity.Entity
}

// RcPage 可能用不到。
// 在隐藏html，使用模板渲染时，才需要用到
type RcPage struct {
	Resource
}

type RcAPI struct {
	Resource
	Path        string `json:"path" gorm:"comment:api路径"`             // api路径
	Description string `json:"description" gorm:"comment:api中文描述"`    // api中文描述
	ApiGroup    string `json:"apiGroup" gorm:"comment:api组"`          // api组
	Method      string `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
}

func (t *RcAPI) TableName() string {
	return "sys_resource_api"
}

type RcMenu struct {
	Resource
	MenuLevel  uint                              `json:"-"`
	ParentId   string                            `json:"parentId" gorm:"comment:父菜单ID"`     // 父菜单ID
	Path       string                            `json:"path" gorm:"comment:路由path"`        // 路由path
	Name       string                            `json:"name" gorm:"comment:路由name"`        // 路由name
	Hidden     bool                              `json:"hidden" gorm:"comment:是否在列表隐藏"`     // 是否在列表隐藏
	Component  string                            `json:"component" gorm:"comment:对应前端文件路径"` // 对应前端文件路径
	Sort       int                               `json:"sort" gorm:"comment:排序标记"`          // 排序标记
	RcMenuMeta `json:"meta" gorm:"comment:附加属性"` // 附加属性
	Roles      []Role                            // 多对多
	Children   []RcMenu                          `json:"children" gorm:"-"`
	Parameters []RcMenuParameter                 `json:"parameters"`
}

func (t *RcMenu) TableName() string {
	return "sys_resource_menu"
}

type RcMenuMeta struct {
	KeepAlive   bool   `json:"keepAlive" gorm:"comment:'是否缓存'"`           // 是否缓存
	DefaultMenu bool   `json:"defaultMenu" gorm:"comment:'是否是基础路由（开发中）'"` // 是否是基础路由（开发中）
	Title       string `json:"title" gorm:"comment:'菜单名'"`                // 菜单名
	Icon        string `json:"icon" gorm:"comment:'菜单图标'"`                // 菜单图标
	CloseTab    bool   `json:"closeTab" gorm:"comment:'自动关闭tab'"`         // 自动关闭tab
}

type RcMenuParameter struct {
	entity.Entity
	SysBaseMenuID uint
	Type          string `json:"type" gorm:"comment:'地址栏携带参数为params还是query'"` // 地址栏携带参数为params还是query
	Key           string `json:"key" gorm:"comment:'地址栏携带参数的key'"`            // 地址栏携带参数的key
	Value         string `json:"value" gorm:"comment:'地址栏携带参数的值'"`            // 地址栏携带参数的值
}

func (t *RcMenuParameter) TableName() string {
	return "sys_resource_menu_parameter"
}
