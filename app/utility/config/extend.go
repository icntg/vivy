package config

import "fmt"

type ExConfig struct {
	Config
	MysqlDSN string
}

func (ths *ExConfig) From(cfg *Config) {
	ths.MysqlDSN = fmt.Sprintf("%s:%s@%s:%d/%s%s",
		cfg.DataSource.Mysql.Username,
		cfg.DataSource.Mysql.Password,
		cfg.DataSource.Mysql.Host,
		cfg.DataSource.Mysql.Port,
		cfg.DataSource.Mysql.Database,
		cfg.DataSource.Mysql.Options,
	)
}
