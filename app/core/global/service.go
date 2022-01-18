package global

import (
	"app/core/middleware"
	"github.com/gin-gonic/gin"
)

type Service struct {
	GinEngine *gin.Engine
}

type WindowsService struct {
	Service
}

type UnixService struct {
	Service
}

type ServiceInterface interface {
	Start()
	AddRoutes(addRouteFunctions ...func(routes *gin.IRoutes))
}

func (ths Service) AddRoutes(addRouteFunctions ...func()) {
	//TODO implement me
	panic("implement me")
}

func (ths Service) Start() {
	gConfig := ConfigInstance()
	if gConfig.Dev.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.ForceConsoleColor()

	ths.GinEngine = gin.New()
	ths.GinEngine.Use(middleware.GinRecovery(LoggersInstance().AccessLogger, true)).Use(middleware.GinLogger(LoggersInstance().AccessLogger))
	ths.AddRoutes(func() {

	})
}

func (ths WindowsService) Start() {
	//TODO implement me
	panic("implement me")
}

func (ths UnixService) Start() {
	//TODO implement me
	panic("implement me")
}
