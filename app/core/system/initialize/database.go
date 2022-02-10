package initialize

import (
	"app/core/global"
	"app/core/utility/common"
	"app/core/utility/errno"
	"database/sql"
	_ "gorm.io/driver/mysql"
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
	r, err := conn.Exec("CREATE DATABASE `vivy` /*!40100 COLLATE 'utf8mb4_general_ci' */")
	common.OutPrintf("r = %v\n", r)
	common.ErrPrintf("err = %v\n", err)
}
