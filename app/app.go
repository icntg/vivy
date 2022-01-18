package main

import (
	"app/core/global"
	"app/utility/errno"
	"app/utility/logger"
	utilityPath "app/utility/path"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

/* 启动入口函数 */
func main() {
	// 命令行参数，获取配置文件等信息
	_ = global.SystemArgsInstance()
	// 加载配置文件
	gConfig := global.ConfigInstance()
	// 初始化日志
	_ = global.LoggersInstance()

	if gConfig.Dev.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.ForceConsoleColor()
	//engine := gin.Default()
	engine := gin.New()
	engine.Use(logger.Middleware()).Use(gin.Recovery())
	addRoute(engine)

	log.SetOutput(os.Stdout)
	log.Printf("gin-server is going to start on [%s] ...\n", gConfig.ServiceHostPort)
	_ = os.Stdout.Sync()

	err = engine.Run(gConfig.ServiceHostPort)
	if nil != err {
		os.Exit(errno.ErrorStartGinService)
		return
	}
}

func addRoute(router *gin.Engine) {
	//apiGroup := router.Group("/api")
	//loginGroup := apiGroup.Group("/login")
	if gin.Mode() == gin.DebugMode {
		router.StaticFS("/", http.Dir("/home/src/vivy/web/dist"))
	} else {
		binaryPath, err := utilityPath.GetBinaryPath()
		if nil != err {
			// log.debug
			router.StaticFS("/", http.Dir("./static"))
		} else {
			router.StaticFS("/", http.Dir(filepath.Join(binaryPath, "static")))
		}

	}
}

func initLogger() error {
	_ = logger.GetOutputLogger()
	_ = logger.GetAccessLogger()
	_ = logger.GetSecureLogger()
	return nil
}
