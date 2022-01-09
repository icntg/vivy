package config

import (
	"fmt"
	"net/url"
)

type Extend struct {
	Config
	ServiceHostPort string
	MysqlDSN        string
}

var globalConfigExtend *Extend = nil

func (ths *Extend) from(cfg *Config) *Extend {
	ths.Config = *cfg
	ths.ServiceHostPort = fmt.Sprintf(
		"%s:%d",
		cfg.Service.HTTP.Host,
		cfg.Service.HTTP.Port)
	ths.MysqlDSN = fmt.Sprintf("%s:%s@%s:%d/%s%s",
		url.QueryEscape(cfg.DataSource.Mysql.Username),
		url.QueryEscape(cfg.DataSource.Mysql.Password),
		cfg.DataSource.Mysql.Host,
		cfg.DataSource.Mysql.Port,
		cfg.DataSource.Mysql.Database,
		cfg.DataSource.Mysql.Options,
	)
	return ths
}

func GetGlobalConfigEx() *Extend {
	if nil == globalConfigExtend {
		globalConfigExtend.from(GetGlobalConfig())
	}
	return globalConfigExtend
}
