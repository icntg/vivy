package global

import (
	"app/core/utility/common"
	flag "github.com/spf13/pflag"
	"sync"
)

type SystemArgs struct {
	ConfigFilename *string
	ConfigTemplate *string
	FlagUsage      string
}

var (
	systemArgsInstance *SystemArgs = nil
	systemArgsOnce     sync.Once
)

// SystemArgsInstance /* 命令行参数处理 */
func SystemArgsInstance() *SystemArgs {
	common.OutPrintf("VIVY backend app start ...\n")
	systemArgsOnce.Do(func() {
		flagArgConfig := flag.StringP("config", "c", "config.yaml", "Using a custom config file")
		// 输出模板
		flagArgTemplate := flag.StringP("output", "o", "", "Output a config file template")
		flag.Lookup("config").NoOptDefVal = "config.yaml"
		flag.Lookup("output").NoOptDefVal = ""
		if len(*flagArgTemplate) == 0 {
			flagArgTemplate = nil
		}
		systemArgsInstance = &SystemArgs{
			flagArgConfig,
			flagArgTemplate,
			flag.CommandLine.FlagUsages(),
		}
	})
	return systemArgsInstance
}
