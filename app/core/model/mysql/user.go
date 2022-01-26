package mysql

import (
	"app/utility/entity"
	"time"
)

type User struct {
	entity.Entity
	Code          string    `json:"code" xorm:"not null comment('工号') VARCHAR(50) unique"`
	Name          string    `json:"name" xorm:"not null comment('姓名') VARCHAR(50) index"`
	Nickname      string    `json:"nickname" xorm:"not null default '' comment('用户登录名') VARCHAR(50) unique"`
	Password      string    `json:"password" xorm:"not null comment('密码') VARCHAR(50)"`
	Salt          string    `json:"salt" xorm:"not null comment('盐') VARCHAR(50)"`
	Token         string    `json:"token" xorm:"not null comment('Google令牌') VARCHAR(100)"`
	Phone         string    `json:"phone" xorm:"not null default '' comment('手机号') VARCHAR(20) index"`
	Avatar        string    `json:"avatar" xorm:"not null default '' comment('头像') VARCHAR(300)"`
	LastLoginTime time.Time `json:"last_login_time" xorm:"not null default '0000-00-00 00:00:00' comment('最近登录时间') DATETIME"`
	LastLoginIp   string    `json:"last_login_ip" xorm:"not null default '' comment('最近登录IP') VARCHAR(50)"`
}
