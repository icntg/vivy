package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"time"
)

func Middleware() gin.HandlerFunc {

	logger := GetAccessLogger()
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logger.Out,
		logrus.FatalLevel: logger.Out,
		logrus.DebugLevel: logger.Out,
		logrus.WarnLevel:  logger.Out,
		logrus.ErrorLevel: logger.Out,
		logrus.PanicLevel: logger.Out,
	}

	logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))

	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIP := c.ClientIP()

		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()
	}
}
