package global

import (
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

func gormInstance() *gorm.DB {
	_gormOnce.Do(func() {
		var (
			err error
		)
		dsn := configInstance().DataSource.MySQL.GetDSN()
		dsnMask := configInstance().DataSource.MySQL.GetDSNWithMask()
		_gormInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			PrepareStmt: true,
			Logger:      *loggersInstance().GORMLogger,
		})
		if nil != err {
			common.ErrPrintf("gorm cannot open with [%s]: %v\n", dsnMask, err)
			common.ErrPrintf("gorm cannot connect database. maybe you should use --init to initialize database?\n")
			os.Exit(errno.ErrorConnectDatabase)
		}
		sqlDB, err := _gormInstance.DB()
		if nil != err {
			common.ErrPrintf("_gormInstance cannot call DB: %v\n", err)
			os.Exit(errno.ErrorConnectDatabase)
		}
		sqlDB.SetMaxOpenConns(configInstance().DataSource.MySQL.MaxOpen)
		sqlDB.SetMaxIdleConns(configInstance().DataSource.MySQL.MaxIdle)
		stats, err := json.Marshal(sqlDB.Stats())
		if nil != err {
			common.ErrPrintf("sqlDB cannot call Stats: %v\n", err)
			os.Exit(errno.ErrorConnectDatabase)
		}
		common.OutPrintf("gorm open with [%s]: %s\n", dsnMask, string(stats))
	})
	return _gormInstance
}

func TestDatabaseConnection() {
	// test mysql

	// test redis
	// test mongodb
}

func TestMySQLConnection() {

}

func TestRedisConnection() {

}

func TestMongoDBConnection() {

}
