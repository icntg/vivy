package initialize

import (
	"app/core/config"
	"app/core/utility/errno"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

type SystemArgs struct {
	ConfigFilename string
}

// InitArgsFactory /* 命令行参数处理 */
func InitArgsFactory() SystemArgs {
	// 配置文件路径
	flagArgConfig := flag.StringP("config", "c", "config.yaml", "Using a custom config file")
	// 输出模板
	flagArgTemplate := flag.StringP("output", "o", "config-template.yaml", "Output a config file template")
	flag.Lookup("config").NoOptDefVal = "config.yaml"
	flag.Lookup("output").NoOptDefVal = ""

	// 输出模板
	n := len(*flagArgTemplate)
	if n > 0 {
		outputConfigTemplate(*flagArgTemplate)
	}
	return SystemArgs{*flagArgConfig}
}

func readConfig(filename string) {
	// 读取配置文件。
	_, err := os.Stat(filename)
	if nil != err && os.IsNotExist(err) {
		// 文件不存在
		log.SetOutput(os.Stderr)
		log.Printf("config file [%s] doesnot not exist!\n", filename)
		os.Exit(errno.ErrorReadConfig)
	}
	// 配置文件存在。尝试进行解析。
	conf, err := config.Read(filename)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Printf("Load config file [%s] failed: %v\n", filename, err)
		os.Exit(errno.ErrorReadConfig)

	}
	log.SetOutput(os.Stdout)
	log.Printf("Load config file [%s] successfully:\n%v\n", filename, conf)
	config.SetGlobalConfig(conf)
}

func outputConfigTemplate(filename string) {
	log.SetOutput(os.Stdout)
	log.Printf("Start to write config template [%s]\n", filename)
	err := config.GenerateTemplate(filename)
	if err == nil {
		log.SetOutput(os.Stdout)
		log.Printf("Write config template [%s] successfully!\n", filename)
		os.Exit(0)
	} else {
		log.SetOutput(os.Stderr)
		log.Printf("Write config template [%s] failed: %v\n", filename, err)
		os.Exit(errno.ErrorGenerateConfigTemplate)
	}
}
