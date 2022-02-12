package system

import (
	"app/core/utility/entity"
	"time"
)

type User struct {
	entity.Entity
	Code          string    `json:"code" gorm:"type:VARCHAR(20);not null;uniqueIndex;comment:工号"`
	Name          string    `json:"name" gorm:"type:VARCHAR(20);not null;comment:姓名"`
	LoginName     string    `json:"login_name" gorm:"type:VARCHAR(20);uniqueIndex;comment:登录名/昵称"`
	EMail         string    `json:"email" gorm:"type:VARCHAR(100);uniqueIndex;comment:邮箱（可登录）"`
	Telephone     string    `json:"telephone" gorm:"type:VARCHAR(100);comment:座机"`
	MobilePhone   string    `json:"mobile_phone" gorm:"type:VARCHAR(100);comment:手机"`
	Avatar        string    `json:"avatar" gorm:"type:TEXT;comment:头像"`
	Password      string    `json:"-" gorm:"type:VARCHAR(50);comment:密码"`
	Salt          string    `json:"-" gorm:"type:VARCHAR(50);comment:密码盐"`
	Token         string    `json:"-" gorm:"type:VARCHAR(50);comment:Google令牌"`
	LastLoginTime time.Time `json:"last_login_time" gorm:"type:DATETIME;comment:最近登录时间"`
	LastLoginIp   string    `json:"last_login_ip" gorm:"type:VARCHAR(50);comment:最近登录IP"`
}

func (t *User) TableName() string {
	return "sys_user"
}
