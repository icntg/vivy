package initialize

import (
	"app/core/global"
	"app/core/utility/common"
	"app/core/utility/crypto"
	"app/core/utility/errno"
	"app/core/web/model/system"
	"bufio"
	"database/sql"
	"encoding/hex"
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strings"
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

	// create database
	{
		_sql := fmt.Sprintf("CREATE DATABASE `%s` /*!40100 COLLATE 'utf8mb4_general_ci' */", gCfg.DataSource.MySQL.Database)
		r, err := conn.Exec(_sql)
		if nil != err {
			common.ErrPrintf("sql cannot create database [%s] with [%s]: %v\n", gCfg.DataSource.MySQL.Database, _sql, err)
		} else {
			common.OutPrintf("sql created database [%s] successfully: %v\n", gCfg.DataSource.MySQL.Database, r)
		}
	}

	// select to database
	{
		_sql := fmt.Sprintf("USE `%s`", gCfg.DataSource.MySQL.Database)
		r, err := conn.Exec(_sql)
		if nil != err {
			common.ErrPrintf("sql cannot use database [%s] with [%s]: %v\n", gCfg.DataSource.MySQL.Database, _sql, err)
			os.Exit(errno.ErrorConnectDatabase)
		} else {
			common.OutPrintf("sql use database [%s] successfully: %v\n", gCfg.DataSource.MySQL.Database, r)
		}
	}

	// wrap connection with GORM
	var gormDB *gorm.DB
	{
		gormDB, err = gorm.Open(mysql.New(mysql.Config{
			Conn: conn,
		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if nil != err {
			common.ErrPrintf("gorm cannot open database with [sql connection]: %v\n", err)
			os.Exit(errno.ErrorConnectDatabase)
		}
	}

	// create tables
	{
		err = gormDB.AutoMigrate(system.Department{})
		err = gormDB.AutoMigrate(system.Resource{})
		err = gormDB.AutoMigrate(system.Role{})
		err = gormDB.AutoMigrate(system.RoleResource{})
		err = gormDB.AutoMigrate(system.User{})
		err = gormDB.AutoMigrate(system.UserRole{})
	}
	return
}

func initAdmin(db *gorm.DB) {
	admin := system.User{}
	admin.Service.Id = common.ObjectIdB32x()

	input := bufio.NewScanner(os.Stdin)

	common.OutPrintf("===== CREATE ADMIN =====\n")
	common.OutPrintf("Please input the username of admin(default: 'admin'): ")
	input.Scan()
	adminName := strings.TrimSpace(input.Text())
	if len(adminName) == 0 {
		adminName = "admin"
	}
	passcode := crypto.Rand(10, true)
	password := hex.EncodeToString(passcode)
	common.OutPrintf("Initial password of [%s] is: %s\n", adminName, password)
	common.OutPrintf("Do you want to use Google Token? (Y/n): ")
	input.Scan()
	answer := strings.TrimSpace(input.Text())
	if answer != "n" && answer != "N" {
		token := crypto.Rand(10, true)

	}
}
