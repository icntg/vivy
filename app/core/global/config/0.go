package config

type Config struct {
	Dev Dev `mapstructure:"dev" json:"dev" yaml:"dev"`
	Zap Zap `mapstructure:"zap" json:"zap" yaml:"zap"`

	Casbin  Casbin  `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	Service Service `mapstructure:"service" json:"service" yaml:"service"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// auto
	// database
	Mysql   MySQL   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	MongoDB MongoDB `mapstructure:"mongodb" json:"mongodb" yaml:"mongodb"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
}

func DefaultConfig() Config {
	return Config{
		Dev{
			Debug: false,
		},
		Zap{
			"info",
			"--TODO--",
			"",
			"./logs",
			true,
			"--TODO--",
			"--TODO--",
			true,
		},
		Casbin{
			"--TODO--",
		},
		Service{
			"127.0.0.1",
			9088,
			"./web",
			true,
		},
		Captcha{
			5,
			120,
			60,
		},
		MySQL{
			"localhost",
			3306,
			"parseTime=true&charset=utf8&loc=Local",
			"vivy",
			"<username>",
			"<password>",
			20,
			100,
			"--TODO--",
			true,
			false,
			false,
		},
		MongoDB{
			"localhost",
			27017,
			"<username>",
			"<password>",
		},
		Redis{
			0,
			"localhost:6379",
			"<password>",
		},
		CORS{
			"--TODO--",
			nil,
		},
	}
}
