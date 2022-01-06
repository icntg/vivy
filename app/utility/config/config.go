package config

import (
	"app/utility/copy"
	"app/utility/errno"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

/*
配置文件管理模块，主要功能：
	1. 配置文件读取
	2. 配置模板生成
*/

var defaultConfig = Config{
	DataSourceConfig{
		MysqlConfig{
			"127.0.0.1",
			3306,
			"<username>",
			"<password>",
			"vivy",
			"?parseTime=true&charset=utf8&loc=Local",
			50,
		},
		MongoDBConfig{
			"127.0.0.1",
			27017,
			"<username>",
			"<password>",
		},
		RedisConfig{
			"127.0.0.1",
			6379,
			"<username>",
			"<password>",
		},
	},
	ServiceConfig{
		HTTPConfig{
			"127.0.0.1",
			9088,
		},
	},
	LoggerConfig{
		"./",
	},
	false,
}

var globalConfig *Config = nil

type MysqlConfig struct {
	Host               string `yaml:"host"`
	Port               uint16 `yaml:"port"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	Database           string `yaml:"database"`
	Options            string `yaml:"options"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`
}

type MongoDBConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type DataSourceConfig struct {
	Mysql   MysqlConfig   `yaml:"mysql"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
	Redis   RedisConfig   `yaml:"redis"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

type ServiceConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

type LoggerConfig struct {
	Path string `yaml:"path"`
}

type Config struct {
	DataSource DataSourceConfig `yaml:"data_source"`
	Service    ServiceConfig    `yaml:"service"`
	Logger     LoggerConfig     `yaml:"logger"`
	Debug      bool             `yaml:"debug"`
}

func GetGlobalConfig() *Config {
	if nil != globalConfig {
		return globalConfig
	}
	log.SetOutput(os.Stderr)
	log.Printf("The various [globalConfig] is not initialized.\n")
	log.SetOutput(os.Stdout)
	log.Printf("Try to initialize with [config.yaml] ...\n")
	gc0, err := Read("config.yaml")
	if nil == err {
		globalConfig = gc0
		return globalConfig
	}
	log.SetOutput(os.Stderr)
	log.Printf("Failed to initialize with [config.yaml].\n")
	log.SetOutput(os.Stdout)
	log.Printf("Try to use default config ...\n")
	var gc1 Config
	err = copy.DeepCopyByGob(gc1, defaultConfig)
	if nil != err {
		globalConfig = &gc1
		return globalConfig
	}
	os.Exit(errno.ErrorReadConfig)
	return nil
}

func SetGlobalConfig(conf *Config) {
	globalConfig = conf
}

func Read(filename string) (*Config, error) {
	conf := &Config{}
	if fp, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		err = yaml.NewDecoder(fp).Decode(conf)
		if err != nil {
			return nil, err
		}
		return conf, nil
	}
}

func GenerateTemplate(filename string) error {
	data, err := yaml.Marshal(defaultConfig)
	if nil != err {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}
