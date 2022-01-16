package main

import (
	"app/core/initialize"
	"app/utility/config"
	"app/utility/errno"
	"app/utility/logger"
	utilityPath "app/utility/path"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

/* 启动入口函数 */
func main() {
	container := dig.New()
	// init args 参数
	container.Provide(initialize.InitArgs)
	// init config (load / write a template) 读取配置。或生成配置模板。
	container.Provide()
	// init logger 根据配置生成logger

	// run app server forever 开始运行后台服务。

	//defer
	var err error
	// 命令行参数处理、加载配置文件
	err = initParam()
	if nil != err {
		os.Exit(errno.ErrorInitParam)
		return
	}
	gConfig := config.GetGlobalConfigEx()

	// 初始化日志模块
	_ = initLogger()

	if gConfig.Debug {
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
