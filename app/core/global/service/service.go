package service

import (
	"app/core/global/config"
	"app/core/global/logger"
	middleware2 "app/core/system/middleware"
	"app/core/utility/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	_serviceInstance *Service = nil
	_serviceOnce     sync.Once
)

func Instance() *Service {
	_serviceOnce.Do(func() {
		gConfig := config.Instance()
		oLog := logger.Instance().OutputLogger
		aLog := logger.Instance().AccessLogger

		if gConfig.Dev.Debug {
			oLog.Info("gin-server uses DebugMode.")
			gin.SetMode(gin.DebugMode)
		} else {
			common.OutPrintf("gin-server uses ReleaseMode.\n")
			oLog.Info("gin-server uses ReleaseMode.")
			gin.SetMode(gin.ReleaseMode)
		}
		gin.ForceConsoleColor()

		instance := Service{}
		instance.GinEngine = gin.New()
		instance.GinEngine.Use(middleware2.GinRecovery(aLog, true)).Use(middleware2.GinLogger(aLog))
		//instance.Start()
		_serviceInstance = &instance
	})
	return _serviceInstance
}

type Service struct {
	GinEngine *gin.Engine
}

type ServiceInterface interface {
	Start()
	AddRoutes(addRouteFunctions ...func(routes *gin.IRoutes))
}

func (ths Service) AddRoutes(addRouteFunctions ...func(engine *gin.Engine)) {
	for _, f := range addRouteFunctions {
		f(ths.GinEngine)
	}
}

func AddStaticRoute(engine *gin.Engine) {
	loggers := logger.Instance()

	if gin.Mode() == gin.DebugMode {
		engine.StaticFS("/", http.Dir("../web/dist"))
	} else {
		binaryPath, err := common.GetBinaryPath()
		if nil != err {
			loggers.OutputLogger.Error("cannot GetBinaryPath")
			engine.StaticFS("/", http.Dir("./static"))
		} else {
			engine.StaticFS("/", http.Dir(filepath.Join(binaryPath, "static")))
		}

	}
}

func (ths Service) Start() {
	oLog := logger.Instance().OutputLogger
	gConfig := config.Instance()

	startMsg := fmt.Sprintf("gin-server is going to start on [%s] ...",
		gConfig.Service.GetServiceAddress())
	oLog.Info(startMsg)
	if !gConfig.Dev.Debug { //
		log.SetOutput(os.Stdout)
		log.Println(startMsg)
		_ = os.Stdout.Sync()
	}
	err := ths.GinEngine.Run(gConfig.Service.GetServiceAddress())
	if nil != err {
		errMsg := fmt.Sprintf("gin-server failed to start: %v", err)
		if !gConfig.Dev.Debug {
			oLog.Error(errMsg)
			log.SetOutput(os.Stderr)
			log.Fatalf(errMsg)
		} else {
			oLog.Fatal(errMsg)
		}
	}
}

type WindowsService struct {
	*Service
}

// Start 不知道这些要干嘛
func (ths WindowsService) Start() {
	//TODO implement me
	panic("implement me")
}

type UnixService struct {
	*Service
}

// Start 不知道这些要干嘛
func (ths UnixService) Start() {
	//TODO implement me
	panic("implement me")
}
