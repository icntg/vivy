package main

import (
	"app/core/global"
)

/* 启动入口函数 */
func main() {
	// 命令行参数，获取配置文件等信息
	_ = global.GetSysArgs()
	// 加载配置文件
	_ = global.GetConfig()
	// 初始化日志
	_ = global.GetLoggers()
	// 启动服务
	service := global.GetGinService()
	service.AddRoutes(
		global.AddStaticRoute,
	)
	service.Start()
}
