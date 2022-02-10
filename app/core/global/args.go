package global

import (
	flag "github.com/spf13/pflag"
	"sync"
)

type SystemArgs struct {
	ConfigFilename   *string
	ConfigTemplate   *string
	FlagInitDatabase *bool
	FlagUsage        string
}

var (
	_systemArgsInstance *SystemArgs = nil
	_systemArgsOnce     sync.Once
)

// systemArgsInstance /* 命令行参数处理 */
func systemArgsInstance() *SystemArgs {
	_systemArgsOnce.Do(func() {
		flagArgConfig := flag.StringP("config", "c", "config.yaml", "Using a custom config file")
		// 输出模板
		flagArgTemplate := flag.StringP("output", "o", "", "Output a config file template")
		flagInitDatabase := flag.BoolP("init-db", "", false, "Initialize MySQL(MariaDB) Database and tables")
		flag.Lookup("config").NoOptDefVal = "config.yaml"
		flag.Lookup("output").NoOptDefVal = ""
		flag.Lookup("init-db")
		flag.Parse()
		if len(*flagArgTemplate) == 0 {
			flagArgTemplate = nil
		}
		_systemArgsInstance = &SystemArgs{
			flagArgConfig,
			flagArgTemplate,
			flagInitDatabase,
			flag.CommandLine.FlagUsages(),
		}
	})

	return _systemArgsInstance
}
