package mysql

import (
	"app/core/global/config"
	"app/core/global/logger"
	"app/core/utility/common"
	"app/core/utility/errno"
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"sync"
)

var (
	_gormInstance *gorm.DB = nil
	_gormOnce     sync.Once
)

func Instance() *gorm.DB {
	_gormOnce.Do(func() {
		var (
			err  error
			gCfg = config.Instance()
		)
		dsn := gCfg.DataSource.MySQL.GetDSN()
		dsnMask := gCfg.DataSource.MySQL.GetMaskedDSN()
		_gormInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			PrepareStmt:                              true,
			Logger:                                   *logger.Instance().GORMLogger,
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if nil != err {
			common.ErrPrintf("gorm cannot open with [%s]: %v\n", dsnMask, err)
			common.ErrPrintf("gorm cannot connect database. maybe you should use --init to initialize database?\n")
			os.Exit(errno.ErrorGORMOpen.Code())
		}
		sqlDB, err := _gormInstance.DB()
		if nil != err {
			common.ErrPrintf("_gormInstance cannot call MaxIdle: %v\n", err)
			common.ErrPrintf("gorm cannot connect database. maybe you should use --init to initialize database?\n")
			os.Exit(errno.ErrorGORMDBHandler.Code())
		}
		sqlDB.SetMaxOpenConns(gCfg.DataSource.MySQL.MaxOpen)
		sqlDB.SetMaxIdleConns(gCfg.DataSource.MySQL.MaxIdle)
		stats, err := json.Marshal(sqlDB.Stats())
		if nil != err {
			common.ErrPrintf("sqlDB cannot call Stats: %v\n", err)
			common.ErrPrintf("gorm cannot connect database. maybe you should use --init to initialize database?\n")
			os.Exit(errno.ErrorGORMDBStats.Code())
		}
		common.OutPrintf("gorm open with [%s]: %s\n", dsnMask, string(stats))
	})
	return _gormInstance
}
