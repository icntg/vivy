package logger

import (
	"app/utility/config"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	outputLogger *logrus.Logger = nil
	accessLogger *logrus.Logger = nil
	secureLogger *logrus.Logger = nil
)

func initLogger(name string, days int) (*logrus.Logger, error) {
	logNameSplit := fmt.Sprintf("%s-%s.log", name, "%Y%m%d")
	logNameLink := name + ".log"
	writer, err := rotatelogs.New(
		filepath.Join(config.GetGlobalConfig().Logger.Path, logNameSplit),
		rotatelogs.WithLinkName(filepath.Join(config.GetGlobalConfig().Logger.Path, logNameLink)),
		rotatelogs.WithMaxAge(time.Duration(days*24)*time.Hour),
		rotatelogs.WithRotationCount(30),
		rotatelogs.WithRotationTime(time.Minute),
	)
	if nil != err {
		log.SetOutput(os.Stderr)
		log.Printf("Failed to create rotatelogs [%s]: %v\n", logNameLink, err)
		return nil, err
	}
	logger := logrus.New()
	// TODO: 修改机制，如果文件初始化失败，使用标准输出。而不是返回错误。
	logger.SetOutput(writer)
	return logger, nil
}

func GetOutputLogger() (*logrus.Logger, error) {
	if nil == outputLogger {
		var err error
		outputLogger, err = initLogger("output", 7)
		if nil != err {
			return nil, err
		}
	}
	return outputLogger, nil
}

func GetAccessLogger() (*logrus.Logger, error) {
	if nil == accessLogger {
		var err error
		accessLogger, err = initLogger("access", 7)
		if nil != err {
			return nil, err
		}
	}
	return accessLogger, nil
}

func GetSecureLogger() (*logrus.Logger, error) {
	if nil == secureLogger {
		var err error
		secureLogger, err = initLogger("secure", 30)
		if nil != err {
			return nil, err
		}
	}
	return secureLogger, nil
}

//
//// LoggerToFile 日志记录到文件
//func LoggerToFile() gin.HandlerFunc {
//
//	logFilePath := config.Log_FILE_PATH
//	logFileName := config.LOG_FILE_NAME
//
//	//日志文件
//	fileName := path.Join(logFilePath, logFileName)
//
//	//写入文件
//	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
//	if err != nil {
//		fmt.Println("err", err)
//	}
//
//	//实例化
//	logger := logrus.New()
//
//	//设置输出
//	logger.Out = src
//
//	//设置日志级别
//	logger.SetLevel(logrus.DebugLevel)
//
//	//设置日志格式
//	logger.SetFormatter(&logrus.TextFormatter{})
//
//	return func(c *gin.Context) {
//		// 开始时间
//		startTime := time.Now()
//
//		// 处理请求
//		c.Next()
//
//		// 结束时间
//		endTime := time.Now()
//
//		// 执行时间
//		latencyTime := endTime.Sub(startTime)
//
//		// 请求方式
//		reqMethod := c.Request.Method
//
//		// 请求路由
//		reqUri := c.Request.RequestURI
//
//		// 状态码
//		statusCode := c.Writer.Status()
//
//		// 请求IP
//		clientIP := c.ClientIP()
//
//		// 日志格式
//		logger.Infof("| %3d | %13v | %15s | %s | %s |",
//			statusCode,
//			latencyTime,
//			clientIP,
//			reqMethod,
//			reqUri,
//		)
//	}
//}

//
//// LoggerToMongo 日志记录到 MongoDB
//func LoggerToMongo() gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//	}
//}
//
//// LoggerToES 日志记录到 ES
//func LoggerToES() gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//	}
//}
//
//// LoggerToMQ 日志记录到 MQ
//func LoggerToMQ() gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//	}
//}
