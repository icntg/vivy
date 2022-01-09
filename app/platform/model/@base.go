package model

import (
	"app/utility/common"
	"app/utility/config"
	"app/utility/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"net/url"
)

var mEngine *xorm.Engine

func init() {
	if mEngine == nil {
		var err error
		cfg := config.GetGlobalConfigEx()
		mEngine, err = xorm.NewEngine(cfg.DataSource.Mysql.DriverName, cfg.MysqlDSN)
		if err != nil {
			oLogger := logger.GetOutputLogger()
			oLogger.Fatalf("Cannot connect to MySQL server with [%s:%s@%s:%d/%s%s]\n",
				url.QueryEscape(cfg.DataSource.Mysql.Username),
				common.CoverWithStars(url.QueryEscape(cfg.DataSource.Mysql.Password)),
				cfg.DataSource.Mysql.Host,
				cfg.DataSource.Mysql.Port,
				cfg.DataSource.Mysql.Database,
				cfg.DataSource.Mysql.Options,
			)
		}
		mEngine.SetMaxIdleConns(cfg.DataSource.Mysql.MaxIdle) //空闲连接
		mEngine.SetMaxOpenConns(cfg.DataSource.Mysql.MaxOpen) //最大连接数
		mEngine.ShowSQL(cfg.DataSource.Mysql.ShowSQL)
		mEngine.ShowExecTime(cfg.DataSource.Mysql.ShowExecTime)
	}
}
