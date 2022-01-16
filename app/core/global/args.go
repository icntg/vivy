package global

import (
	flag "github.com/spf13/pflag"
	"sync"
)

type SystemArgs struct {
	ConfigFilename *string
	ConfigTemplate *string
}

var (
	systemArgsInstance *SystemArgs = nil
	systemArgsOnce     sync.Once
)

// SystemArgsInstance /* 命令行参数处理 */
func SystemArgsInstance() *SystemArgs {
	systemArgsOnce.Do(func() {
		flagArgConfig := flag.StringP("config", "c", "config.yaml", "Using a custom config file")
		// 输出模板
		flagArgTemplate := flag.StringP("output", "o", "config-template.yaml", "Output a config file template")
		flag.Lookup("config").NoOptDefVal = "config.yaml"
		flag.Lookup("output").NoOptDefVal = ""
		if len(*flagArgTemplate) == 0 {
			flagArgTemplate = nil
		}
		systemArgsInstance = &SystemArgs{
			flagArgConfig,
			flagArgTemplate,
		}
	})
	return systemArgsInstance
}
