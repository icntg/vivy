package config

import "fmt"

type MySQL struct {
	Host         string `mapstructure:"host" json:"host" yaml:"host"`             // 服务器地址
	Port         uint16 `mapstructure:"port" json:"port" yaml:"port"`             // 端口
	Option       string `mapstructure:"option" json:"option" yaml:"option"`       // 高级配置
	Database     string `mapstructure:"database" json:"database" yaml:"database"` // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"` // 数据库用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"` // 数据库密码
	MaxIdle      int    `mapstructure:"max-idle" json:"maxIdle" yaml:"max-idle"`  // 空闲中的最大连接数
	MaxOpen      int    `mapstructure:"max-open" json:"maxOpen" yaml:"max-open"`  // 打开到数据库的最大连接数
	LogMode      string `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`  // 是否开启Gorm全局日志
	LogZap       bool   `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"`     // 是否通过zap写入日志文件
	ShowSQL      bool   `mapstructure:"show-sql" json:"showSQL" yaml:"show-sql"`
	ShowExecTime bool   `mapstructure:"show-exec-time" json:"showExecTime" yaml:"show-exec-time"`
}

func (m *MySQL) GetDSN() string {
	var ret string
	ret = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.Database,
	)
	if len(m.Option) > 0 {
		ret += "?" + m.Option
	}
	return ret
}
