package model

import (
	"app/utility/entity"
	"time"
)

type User struct {
	entity.Entity
	Name          string    `json:"name" xorm:"not null comment('姓名') VARCHAR(50)"`
	Nickname      string    `json:"nickname" xorm:"not null default '' comment('用户登录名') unique VARCHAR(50)"`
	Password      string    `json:"password" xorm:"not null comment('密码') index VARCHAR(50)"`
	Salt          string    `json:"salt" xorm:"not null comment('盐') VARCHAR(32)"`
	Token         string    `json:"token" xorm:"not null comment('Google令牌') VARCHAR(64)"`
	Phone         string    `json:"phone" xorm:"not null default '' comment('手机号') VARCHAR(11)"`
	Avatar        string    `json:"avatar" xorm:"not null default '' comment('头像') VARCHAR(300)"`
	LastLoginTime time.Time `json:"last_login_time" xorm:"not null default '0000-00-00 00:00:00' comment('上次登录时间') DATETIME"`
	LastLoginIp   string    `json:"last_login_ip" xorm:"not null default '' comment('最近登录IP') VARCHAR(50)"`
}
