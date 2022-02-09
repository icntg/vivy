package global

// 全局变量，包括：
//1系统参数
//2配置文件信息
//3日志对象
//4数据访问对象
//5http服务对象

import (
	"app/core/system/config"
	"gorm.io/gorm"
)

var (
	constTest     = false
	GetSysArgs    func() *SystemArgs
	GetConfig     func() *config.Config
	GetLoggers    func() *Loggers
	GetGorm       func() *gorm.DB
	GetGinService func() *Service
)

func init() {
	if !constTest {
		GetSysArgs = systemArgsInstance
		GetConfig = configInstance
		GetLoggers = loggersInstance
		GetGorm = gormInstance
		GetGinService = serviceInstance
	}

}
