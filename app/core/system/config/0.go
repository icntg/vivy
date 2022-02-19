package config

import "app/core/system/config/ds"

type Config struct {
	Dev Dev `mapstructure:"dev" json:"dev" yaml:"dev"`
	Zap Zap `mapstructure:"zap" json:"zap" yaml:"zap"`

	Casbin  Casbin  `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	Service Service `mapstructure:"service" json:"service" yaml:"service"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// auto
	// database
	DataSource ds.DataSource `mapstructure:"data-source" json:"dataSource" yaml:"data-source"`
	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
}

func DefaultConfig() Config {
	return Config{
		Dev: Dev{
			Debug: false,
		},
		Zap: Zap{
			"info",
			"--TODO--",
			"",
			"./logs",
			true,
			"--TODO--",
			"--TODO--",
			true,
		},
		Casbin: Casbin{
			"--TODO--",
		},
		Service: Service{
			"127.0.0.1",
			9088,
			"./web",
			true,
			"<session-secret-in-HEX>",
			nil,
			3600,
		},
		Captcha: Captcha{
			5,
			120,
			60,
		},
		DataSource: ds.DataSource{
			MySQL: ds.MySQL{
				Host:     "localhost",
				Port:     3306,
				Option:   "parseTime=true&charset=utf8mb4&loc=Local",
				Database: "vivy",
				Username: "<username>",
				Password: "<password>",
				MaxIdle:  20,
				MaxOpen:  100,
				LogMode:  "--TODO--",
				LogZap:   true,
			},
			MongoDB: ds.MongoDB{
				Host:     "localhost",
				Port:     27017,
				Username: "<username>",
				Password: "<password>",
			},
			Redis: ds.Redis{
				MaxIdle:  10,
				Protocol: "tcp",
				Address:  "localhost:6379",
				Password: "<password>",
			},
		},
		Cors: CORS{
			"--TODO--",
			nil,
		},
	}
}
