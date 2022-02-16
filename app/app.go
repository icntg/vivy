package main

import (
	"app/core/global"
	"app/core/global/service"
	"app/core/system/initialize"
	"app/core/utility/common"
	"app/core/utility/errno"
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
		os.Exit(errno.ErrorSuccess.Code())
	} else {
		_ = global.GetGORM()
		common.OutPrintf("success to initialize data access object.\n")
		_ = global.GetRedis()
		common.OutPrintf("success to initialize redis store.\n")
	}
	// TODO: 启动mongodb
	// 启动服务
	httpService := global.GetGinService()
	httpService.AddRoutes(
		service.AddStaticRoute,
	)
	httpService.Start()
}
