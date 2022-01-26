package global

import (
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		dsn := configInstance().Mysql.GetDSN()
		dsnMask := configInstance().Mysql.GetDSNWithMask()
		_gormInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			PrepareStmt: true,
		})
		if nil != err {
			loggersInstance().OutPutLogger.Fatalf("gorm cannot open with [%s]: %v\n", dsnMask, err)
		}
		sqlDB, err := _gormInstance.DB()
		if nil != err {
			loggersInstance().OutPutLogger.Fatalf("_gormInstance cannot call DB: %v\n", err)
		}
		sqlDB.SetMaxOpenConns(configInstance().Mysql.MaxOpen)
		sqlDB.SetMaxIdleConns(configInstance().Mysql.MaxIdle)
		stats, err := json.Marshal(sqlDB.Stats())
		if nil != err {
			loggersInstance().OutPutLogger.Fatalf("sqlDB cannot call Stats: %v\n", err)
		}
		loggersInstance().OutPutLogger.Infof("gorm open with [%s]: %s\n", dsnMask, string(stats))
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
