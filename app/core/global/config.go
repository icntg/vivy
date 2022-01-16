package global

import (
	"app/core/global/config"
	"log"
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
		if nil != systemArgs {
			log.SetOutput(os.Stderr)
			log.Fatalf("TODO:")
		}
		if systemArgs.ConfigTemplate != nil && len(*systemArgs.ConfigTemplate) > 0 {
			// 输出配置模板
		} else {
			// 读取配置文件
		}

	})
	return configInstance
}
