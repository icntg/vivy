package global

import (
	"app/core/global/config"
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
