package initialize

import (
	"app/core/global"
	"app/core/utility/common"
	"app/core/utility/errno"
	"app/core/web/model/system"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func InitDatabase() {
	var (
		gCfg = global.GetConfig()
		//gDb = global.GetGORM()
	)
	conn, err := sql.Open("mysql", gCfg.DataSource.MySQL.GetDSNWithOutDatabase())
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	if nil != err {
		common.ErrPrintf("sql cannot open database with [%s]: %v\n", gCfg.DataSource.MySQL.GetMaskedDSNWithOutDatabase(), err)
		os.Exit(errno.ErrorConnectDatabase)
	}
	common.OutPrintf("conn = %v\n", conn)
	_sql := fmt.Sprintf("CREATE DATABASE `%s` /*!40100 COLLATE 'utf8mb4_general_ci' */", gCfg.DataSource.MySQL.Database)
	r, err := conn.Exec(_sql)
	if nil != err {
		common.ErrPrintf("sql cannot create database [%s] with [%s]: %v\n", gCfg.DataSource.MySQL.Database, _sql, err)
	} else {
		common.OutPrintf("sql created database [%s] successfully: %v\n", gCfg.DataSource.MySQL.Database, r)
	}
	_sql = fmt.Sprintf("USE `%s`", gCfg.DataSource.MySQL.Database)
	r, err = conn.Exec(_sql)
	if nil != err {
		common.ErrPrintf("sql cannot use database [%s] with [%s]: %v\n", gCfg.DataSource.MySQL.Database, _sql, err)
		os.Exit(errno.ErrorConnectDatabase)
	} else {
		common.OutPrintf("sql use database [%s] successfully: %v\n", gCfg.DataSource.MySQL.Database, r)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: conn,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if nil != err {
		common.ErrPrintf("gorm cannot open database with [sql connection]: %v\n", err)
		os.Exit(errno.ErrorConnectDatabase)
	}
	err = gormDB.AutoMigrate(system.User{})
	return
}
