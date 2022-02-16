package config

import (
	"app/core/global/args"
	"app/core/system/config"
	"app/core/utility/common"
	"app/core/utility/errno"
	"bytes"
	"encoding/hex"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
	"os"
	"path"
	"sync"
)

var (
	_configInstance *config.Config = nil
	_configOnce     sync.Once
)

func Instance() *config.Config {
	_configOnce.Do(func() {
		systemArgs := args.Instance()
		if nil == systemArgs {
			errno.ErrorSystemArgsIsNil.Exit()
			// 但是上面应该不可能遇到。除非手贱。
		}
		if systemArgs.ConfigTemplate != nil && len(*systemArgs.ConfigTemplate) > 0 {
			// 输出配置模板
			SaveToYamlFile(*systemArgs.ConfigTemplate)
			// exit
		} else {
			// 读取配置文件
			cfg := ReadFromYamlFile(*systemArgs.ConfigFilename)
			_configInstance = &cfg
		}
	})
	return _configInstance
}

func ReadFromYamlFile(filename string) config.Config {
	systemArgs := args.Instance()

	if !common.FileExists(filename) {
		common.ErrPrintf("config file [%s] does not exist.\n", filename)
		common.OutPrintf(systemArgs.FlagUsage)
		os.Exit(errno.ErrorConfigNotExist.Code())
	}
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if nil != err {
		common.ErrPrintf("viper cannot read config [%s]: %v\n", filename, err)
		common.OutPrintf(systemArgs.FlagUsage)
		os.Exit(errno.ErrorConfigRead.Code())
	}
	ret := config.DefaultConfig() // 默认配置
	err = v.Unmarshal(&ret)       // 外部配置文件覆盖默认配置
	if nil != err {
		common.ErrPrintf("viper cannot unmarshal data [%v]: %v\n", v, err)
		common.OutPrintf(systemArgs.FlagUsage)
		os.Exit(errno.ErrorConfigUnmarshal.Code())
	}

	ret.Service.SessionSecretBytes, err = hex.DecodeString(ret.Service.SessionSecret)
	if nil != err {
		common.ErrPrintf("hex cannot decode SessionSecret: %v\n", err)
		os.Exit(errno.ErrorConfigDecodeHex.Code())
	}
	return ret
}

func SaveToYamlFile(filename string) {
	defaultConfig := config.DefaultConfig()
	bs, err := yaml.Marshal(defaultConfig)
	if nil != err {
		common.ErrPrintf("yaml cannot marshal DefaultConfig [%v]: %v\n", defaultConfig, err)
		os.Exit(errno.ErrorConfigMarshal.Code())
	}
	v := viper.New()
	v.SetConfigFile(filename)
	curDir := common.GetCurrentDirectory()
	v.AddConfigPath(curDir)
	v.SetConfigType("yaml")
	err = v.ReadConfig(bytes.NewBuffer(bs))
	if nil != err {
		common.ErrPrintf("viper cannot read in data [%v]: %v\n", string(bs), err)
		os.Exit(errno.ErrorConfigViperRead.Code())
	}
	err = v.SafeWriteConfig()
	if nil != err {
		common.ErrPrintf("viper cannot write config [%v]: %v\n", filename, err)
		os.Exit(errno.ErrorConfigWrite.Code())
	}
	common.OutPrintf("config template is written to [%s].", path.Join(curDir, filename))
	os.Exit(errno.ErrorSuccess.Code())
}
