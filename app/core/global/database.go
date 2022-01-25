package global

import (
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	gormInstance *gorm.DB = nil
	gormOnce     sync.Once
)

func GormInstance() *gorm.DB {
	gormOnce.Do(func() {
		var (
			err error
		)
		dsn := ConfigInstance().Mysql.GetDSN()
		dsnMask := ConfigInstance().Mysql.GetDSNWithMask()
		gormInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if nil != err {
			LoggersInstance().OutPutLogger.Fatalf("gorm cannot open with [%s]: %v\n", dsnMask, err)
		}
		sqlDB, err := gormInstance.DB()
		if nil != err {
			LoggersInstance().OutPutLogger.Fatalf("gormInstance cannot call DB: %v\n", err)
		}
		sqlDB.SetMaxOpenConns(ConfigInstance().Mysql.MaxOpen)
		sqlDB.SetMaxIdleConns(ConfigInstance().Mysql.MaxIdle)
		stats, err := json.Marshal(sqlDB.Stats())
		if nil != err {
			LoggersInstance().OutPutLogger.Fatalf("sqlDB cannot call Stats: %v\n", err)
		}
		LoggersInstance().OutPutLogger.Infof("gorm open with [%s]: %s\n", dsnMask, string(stats))
	})
	return gormInstance
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
