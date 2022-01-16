package config

import (
	"bytes"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Zap Zap `mapstructure:"zap" json:"zap" yaml:"zap"`

	Casbin  Casbin  `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// auto
	// database
	Mysql   MySQL   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	MongoDB MongoDB `mapstructure:"mongodb" json:"mongodb" yaml:"mongodb"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
}

func (ths *Config) ReadFromYamlFile(filename string) error {
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if nil != err {
		return err
	}
	err = v.Unmarshal(ths)
	if nil != err {
		return err
	}

	return nil
}

func SaveToYamlFile(filename string) error {
	bs, err := yaml.Marshal(defaultConfig())
	if nil != err {
		return err
	}
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	err = v.ReadConfig(bytes.NewBuffer(bs))
	if nil != err {
		return err
	}
	err = v.SafeWriteConfig()
	if nil != err {
		return err
	}
	return nil
}

func defaultConfig() Config {
	return Config{
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
		System{
			"--TODO--",
			9088,
			"mysql",
			"--TODO--",
			true,
			0,
			0,
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
			50,
			200,
			"--TODO--",
			true,
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
