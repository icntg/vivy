package model

import (
	"app/core/global"
	"app/core/utility/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"net/url"
)

var (
	mEngine *xorm.Engine
	gCfg    = global.ConfigInstance()
	oLog    = global.LoggersInstance().OutPutLogger
)

func init() {
	if mEngine == nil {
		var err error
		mEngine, err = xorm.NewEngine("mysql", gCfg.Mysql.GetDSN())
		if err != nil {
			oLog.Fatalf("Cannot connect to MySQL server with [%s:%s@%s:%d/%s%s]\n",
				url.QueryEscape(gCfg.Mysql.Username),
				common.CoverWithStars(url.QueryEscape(gCfg.Mysql.Password)),
				gCfg.Mysql.Host,
				gCfg.Mysql.Port,
				gCfg.Mysql.Database,
				gCfg.Mysql.Option,
			)
		}
		mEngine.SetMaxIdleConns(gCfg.Mysql.MaxIdle) //空闲连接
		mEngine.SetMaxOpenConns(gCfg.Mysql.MaxOpen) //最大连接数
		mEngine.ShowSQL(gCfg.Mysql.ShowSQL)
		mEngine.ShowExecTime(gCfg.Mysql.ShowExecTime)
	}
}
