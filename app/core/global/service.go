package global

import (
	"app/core/middleware"
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
	serviceInstance *Service = nil
	serviceOnce     sync.Once
)

func ServiceInstance() *Service {
	serviceOnce.Do(func() {
		gConfig := ConfigInstance()
		oLog := LoggersInstance().OutPutLogger

		if gConfig.Dev.Debug {
			oLog.Info("gin-server uses DebugMode.")
			gin.SetMode(gin.DebugMode)
		} else {
			oLog.Info("gin-server uses ReleaseMode.")
			gin.SetMode(gin.ReleaseMode)
		}
		gin.ForceConsoleColor()

		instance := Service{}
		instance.GinEngine = gin.New()
		instance.GinEngine.Use(middleware.GinRecovery(LoggersInstance().AccessLogger, true)).Use(middleware.GinLogger(LoggersInstance().AccessLogger))
		//instance.Start()
		serviceInstance = &instance
	})
	return serviceInstance
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
	loggers := LoggersInstance()

	if gin.Mode() == gin.DebugMode {
		engine.StaticFS("/", http.Dir("../web/dist"))
	} else {
		binaryPath, err := common.GetBinaryPath()
		if nil != err {
			loggers.OutPutLogger.Error("cannot GetBinaryPath")
			engine.StaticFS("/", http.Dir("./static"))
		} else {
			engine.StaticFS("/", http.Dir(filepath.Join(binaryPath, "static")))
		}

	}
}

func (ths Service) Start() {
	oLog := LoggersInstance().OutPutLogger
	gConfig := ConfigInstance()

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
