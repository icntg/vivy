package main

import (
	"app/utility/config"
	"app/utility/errno"
	"app/utility/logger"
	utilityPath "app/utility/path"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

/* 启动入口函数 */
func main() {
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

/* 命令行参数处理 */
func initParam() error {
	// 配置文件路径
	var flagParamConfig = flag.String("c", "config.yaml", "Using a custom config file")
	// 输出模板
	var flagParamTemplate = flag.String("o", "", "Output a config file template")
	flag.Parse()

	// 输出模板
	n := len(*flagParamTemplate)
	if n > 0 {
		log.SetOutput(os.Stdout)
		log.Printf("Start to write config template [%s]\n", *flagParamTemplate)
		err := config.GenerateTemplate(*flagParamTemplate)
		if err == nil {
			log.SetOutput(os.Stdout)
			log.Printf("Write config template [%s] successfully!\n", *flagParamTemplate)
		} else {
			log.SetOutput(os.Stderr)
			log.Printf("Write config template [%s] failed: %v\n", *flagParamTemplate, err)
		}
		os.Exit(0)
	}

	// 读取配置文件。
	_, err := os.Stat(*flagParamConfig)
	if nil != err && os.IsNotExist(err) {
		// 文件不存在
		log.SetOutput(os.Stderr)
		log.Printf("config file [%s] doesnot not exist!\n", *flagParamConfig)
		return err
	}
	// 配置文件存在。尝试进行解析。
	conf, err := config.Read(*flagParamConfig)
	if err == nil {
		log.SetOutput(os.Stdout)
		log.Printf("Load config file [%s] successfully:\n%v\n", *flagParamConfig, conf)
		config.SetGlobalConfig(conf)
	} else {
		log.SetOutput(os.Stderr)
		log.Printf("Load config file [%s] failed: %v\n", *flagParamConfig, err)
	}
	return err
}

func initLogger() error {
	_ = logger.GetOutputLogger()
	_ = logger.GetAccessLogger()
	_ = logger.GetSecureLogger()
	return nil
}
