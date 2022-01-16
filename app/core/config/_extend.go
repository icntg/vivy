//package config
//
//import (
//	"fmt"
//	"net/url"
//)
//
//type Extend struct {
//	*Config
//	ServiceHostPort string
//	MysqlDSN        string
//}
//
//var globalConfigExtend *Extend = nil
//
//func (ths *Extend) from(cfg *Config) *Extend {
//	ths.Config = cfg
//	ths.ServiceHostPort = fmt.Sprintf(
//		"%s:%d",
//		cfg.Service.HTTP.Host,
//		cfg.Service.HTTP.Port)
//	ths.MysqlDSN = fmt.Sprintf("%s:%s@%s:%d/%s%s",
//		url.QueryEscape(cfg.DataSource.MySQL.Username),
//		url.QueryEscape(cfg.DataSource.MySQL.Password),
//		cfg.DataSource.MySQL.Host,
//		cfg.DataSource.MySQL.Port,
//		cfg.DataSource.MySQL.Database,
//		cfg.DataSource.MySQL.Options,
//	)
//	return ths
//}
//
//func GetGlobalConfigEx() *Extend {
//	if nil == globalConfigExtend {
//		globalConfigExtend = &Extend{}
//		globalConfigExtend.from(GetGlobalConfig())
//	}
//	return globalConfigExtend
//}
