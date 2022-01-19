package global

import (
	"app/core/global/config"
	"app/core/utility/common"
	"app/core/utility/errno"
	"os"
	"sync"
)

var (
	configInstance *config.Config = nil
	configOnce     sync.Once
)

func ConfigInstance() *config.Config {
	configOnce.Do(func() {
		systemArgs := SystemArgsInstance()
		if nil == systemArgs {
			common.ErrPrintf("require system args.\n")
			os.Exit(errno.ErrorSystemArgs)
		}
		if systemArgs.ConfigTemplate != nil && len(*systemArgs.ConfigTemplate) > 0 {
			// 输出配置模板
			config.SaveToYamlFile(*systemArgs.ConfigTemplate)
			// exit
		} else {
			// 读取配置文件
			cfg := config.ReadFromYamlFile(*systemArgs.ConfigFilename)
			configInstance = &cfg
		}
	})
	return configInstance
}
