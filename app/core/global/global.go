package global

// 全局变量，包括：
//1系统参数
//2配置文件信息
//3日志对象
//4数据访问对象
//5http服务对象

import (
	"app/core/global/args"
	globalConfig "app/core/global/config"
	"app/core/global/database"
	"app/core/global/logger"
	"app/core/global/service"
	"app/core/system/config"
	"gorm.io/gorm"
)

const (
	ProductName = "VIVY"
)

var (
	constTest     = false
	GetSysArgs    func() *args.SystemArgs
	GetConfig     func() *config.Config
	GetLoggers    func() *logger.Loggers
	GetGORM       func() *gorm.DB
	GetGinService func() *service.Service
)

func init() {
	if !constTest {
		GetSysArgs = args.Instance
		GetConfig = globalConfig.Instance
		GetLoggers = logger.Instance
		GetGORM = database.Instance
		GetGinService = service.Instance
	}

}
