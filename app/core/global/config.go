package global

import (
	"app/core/system/config"
	"app/core/utility/common"
	"app/core/utility/errno"
	"bytes"
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

func configInstance() *config.Config {
	_configOnce.Do(func() {
		systemArgs := systemArgsInstance()
		if nil == systemArgs {
			common.ErrPrintf("require system args.\n")
			os.Exit(errno.ErrorSystemArgs)
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
	systemArgs := systemArgsInstance()

	if !common.FileExists(filename) {
		common.ErrPrintf("config file [%s] does not exist.\n", filename)
		common.OutPrintf(systemArgs.FlagUsage)
		os.Exit(errno.ErrorReadConfig)
	}
	v := viper.New()
	v.SetConfigFile(filename)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if nil != err {
		common.ErrPrintf("viper cannot read config [%s]: %v\n", filename, err)
		common.OutPrintf(systemArgs.FlagUsage)
		os.Exit(errno.ErrorReadConfig)
	}
	ret := config.DefaultConfig()
	err = v.Unmarshal(&ret)
	if nil != err {
		common.ErrPrintf("viper cannot unmarshal data [%v]: %v\n", v, err)
		common.OutPrintf(systemArgs.FlagUsage)
		os.Exit(errno.ErrorReadConfig)
	}
	return ret
}

func SaveToYamlFile(filename string) {
	bs, err := yaml.Marshal(config.DefaultConfig())
	if nil != err {
		common.ErrPrintf("yaml cannot marshal DefaultConfig: %v\n", err)
		os.Exit(errno.ErrorGenerateConfigTemplate)
	}
	v := viper.New()
	v.SetConfigFile(filename)
	curDir := common.GetCurrentDirectory()
	v.AddConfigPath(curDir)
	v.SetConfigType("yaml")
	err = v.ReadConfig(bytes.NewBuffer(bs))
	if nil != err {
		common.ErrPrintf("viper cannot read in data [%v]: %v\n", string(bs), err)
		os.Exit(errno.ErrorGenerateConfigTemplate)
	}
	err = v.SafeWriteConfig()
	if nil != err {
		common.ErrPrintf("viper cannot write config [%v]: %v\n", filename, err)
		os.Exit(errno.ErrorGenerateConfigTemplate)
	}
	common.OutPrintf("config template is written to [%s].", path.Join(curDir, filename))
	os.Exit(0)
}
