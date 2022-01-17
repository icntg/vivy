package global

import (
	"app/core/global/config"
	"app/core/utility/common"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

const (
	outputFilename = "output"
	accessFilename = "access"
	secureFilename = "secure"
)

var (
	logLevel = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}
)

type Loggers struct {
	OutPutLogger *zap.SugaredLogger
	AccessLogger *zap.Logger
	SecureLogger *zap.SugaredLogger
}

var (
	loggersInstance *Loggers = nil
	loggersOnce     sync.Once
)

func LoggersInstance() *Loggers {
	loggersOnce.Do(func() {
		cfgInst := ConfigInstance()
		if nil != cfgInst {
			log.SetOutput(os.Stderr)
			log.Fatalln("require config information.")
		}
		loggers := initLoggers(cfgInst)
		loggersInstance = &loggers
	})
	return loggersInstance
}

func initLoggers(cfgInst *config.Config) Loggers {
	var (
		logWriter io.Writer
	)
	logDir := path.Join(cfgInst.Zap.Director)
	if !common.FileExists(logDir) {
		file, err := os.Create(logDir)
		defer func() {
			err := file.Close()
			if nil != err {
				log.SetOutput(os.Stderr)
				log.Fatalf("cannot close Log Directory [%s]: %v\n", logDir, err)
			}
		}()
		if nil != err {
			log.SetOutput(os.Stderr)
			log.Fatalf("cannot create Log Directory [%s]: %v\n", logDir, err)
		}
	}
	encoder := initEncoder()

	loggers := Loggers{nil, nil, nil}
	// accessLogger
	{
		logPath := path.Join(logDir, accessFilename)
		logWriter = getWriter(logPath, 6*31, 7)
		if cfgInst.Dev.Debug {
			logWriter = io.MultiWriter(logWriter, os.Stdout)
		}
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(logWriter), logLevel[cfgInst.Zap.Level]),
		)
		tmpLog := zap.New(core, zap.AddCaller())
		loggers.AccessLogger = tmpLog
	}
	// outputLogger
	{
		logPath := path.Join(logDir, outputFilename)
		logWriter = getWriter(logPath, 6*31, 7)
		if cfgInst.Dev.Debug {
			logWriter = io.MultiWriter(logWriter, os.Stdout)
		}
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(logWriter), logLevel[cfgInst.Zap.Level]),
		)
		tmpLog := zap.New(core, zap.AddCaller())
		loggers.OutPutLogger = tmpLog.Sugar()
	}
	// secureLogger
	{
		logPath := path.Join(logDir, secureFilename)
		logWriter = getWriter(logPath, 6*31, 31)
		if cfgInst.Dev.Debug {
			logWriter = io.MultiWriter(logWriter, os.Stdout)
		}
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(logWriter), logLevel[cfgInst.Zap.Level]),
		)
		tmpLog := zap.New(core, zap.AddCaller())
		loggers.SecureLogger = tmpLog.Sugar()
	}
	return loggers
}

//初始化Encoder
func initEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "time",
		CallerKey:   "file",
		EncodeLevel: zapcore.CapitalLevelEncoder, //基本zapcore.LowercaseLevelEncoder。将日志级别字符串转化为小写
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeCaller: zapcore.ShortCallerEncoder, //一般zapcore.ShortCallerEncoder，以包/文件:行号 格式化调用堆栈
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) { //一般zapcore.SecondsDurationEncoder,执行消耗的时间转化成浮点型的秒
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
}

//日志文件切割
func getWriter(filenamePrefix string, maxAgeDays int, rotationTimeDays int) io.Writer {
	// 保存30天内的日志，每24小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filenamePrefix+"-%Y%m%d.log",
		rotatelogs.WithLinkName(filenamePrefix+".log"),
		rotatelogs.WithMaxAge(time.Duration(maxAgeDays*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(rotationTimeDays*24)*time.Hour),
	)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatalf("cannot initialize SecureLogger: %v\n", err)
	}
	return hook
}
