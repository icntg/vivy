package config

import (
	"bytes"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

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

func ReadFromYamlFile(filename string) Config {
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if nil != err {
		log.SetOutput(os.Stderr)
		log.Fatalf("viper cannot read config [%s]: %v\n", filename, err)
	}
	ret := defaultConfig()
	err = v.Unmarshal(&ret)
	if nil != err {
		log.SetOutput(os.Stderr)
		log.Fatalf("viper cannot unmarshal data [%v]: %v\n", v, err)
	}
	return ret
}

func SaveToYamlFile(filename string) {
	bs, err := yaml.Marshal(defaultConfig())
	if nil != err {
		log.SetOutput(os.Stderr)
		log.Fatalf("yaml cannot marshal DefaultConfig: %v\n", err)
	}
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	err = v.ReadConfig(bytes.NewBuffer(bs))
	if nil != err {
		log.SetOutput(os.Stderr)
		log.Fatalf("viper cannot read in data [%v]: %v\n", string(bs), err)
	}
	err = v.SafeWriteConfig()
	if nil != err {
		log.SetOutput(os.Stderr)
		log.Fatalf("viper cannot write config [%v]: %v\n", filename, err)
	}
}

func defaultConfig() Config {
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
