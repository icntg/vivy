package main

import (
	"app/core/global"
	"app/core/system/initialize"
	"app/core/utility/common"
	"os"
)

/* 启动入口函数 */
func main() {
	// 命令行参数，获取配置文件等信息
	sysArgs := global.GetSysArgs()
	common.OutPrintf("success to get command line args.\n")
	// 加载配置文件
	_ = global.GetConfig()
	common.OutPrintf("success to read config.\n")
	// 初始化日志
	_ = global.GetLoggers()
	common.OutPrintf("success to initialize loggers.\n")
	// 根据启动参数初始化数据库
	if *sysArgs.FlagInitDatabase {
		initialize.InitDatabase()
		os.Exit(0)
	} else {
		_ = global.GetGORM()
		common.OutPrintf("success to initialize loggers.\n")

	}
	// 启动服务
	service := global.GetGinService()
	service.AddRoutes(
		global.AddStaticRoute,
	)
	service.Start()
}
