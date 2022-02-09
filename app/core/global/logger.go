package global

import (
	"app/core/system/config"
	"app/core/utility/common"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"path"
	"runtime"
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
	GormLogger   *logger.Interface
}

var (
	_loggersInstance *Loggers = nil
	_loggersOnce     sync.Once
)

func loggersInstance() *Loggers {
	_loggersOnce.Do(func() {
		cfgInst := configInstance()
		if nil == cfgInst {
			log.SetOutput(os.Stderr)
			log.Fatalln("require config information.")
		}
		loggers := initLoggers(cfgInst)
		_loggersInstance = &loggers
	})
	return _loggersInstance
}

func initLoggers(cfgInst *config.Config) Loggers {
	var (
		logWriter    io.Writer
		outputWriter io.Writer
	)
	logDir := path.Join(cfgInst.Zap.Director)
	if !common.FileExists(logDir) {
		err := os.MkdirAll(logDir, 0o755)
		if nil != err {
			log.SetOutput(os.Stderr)
			log.Fatalf("cannot create Log Directory [%s]: %v\n", logDir, err)
		}
	}
	encoder := initEncoder()

	loggers := Loggers{nil, nil, nil, nil}
	// accessLogger
	{
		logPath := path.Join(logDir, accessFilename)
		logWriter = getWriter(logPath, 6*31, 7)
		if cfgInst.Dev.Debug {
			logWriter = io.MultiWriter(logWriter, os.Stderr)
		}
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(logWriter), logLevel[cfgInst.Zap.Level]),
		)
		tmpLog := zap.New(core, zap.AddCaller())
		loggers.AccessLogger = tmpLog
	}
	// outputLogger
	outputWriter = getWriter(path.Join(logDir, outputFilename), 6*31, 7)
	if cfgInst.Dev.Debug {
		outputWriter = io.MultiWriter(outputWriter, os.Stderr)
	}
	{
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(outputWriter), logLevel[cfgInst.Zap.Level]),
		)
		tmpLog := zap.New(core, zap.AddCaller())
		loggers.OutPutLogger = tmpLog.Sugar()
	}
	// secureLogger
	{
		logPath := path.Join(logDir, secureFilename)
		logWriter = getWriter(logPath, 6*31, 31)
		if cfgInst.Dev.Debug {
			logWriter = io.MultiWriter(logWriter, os.Stderr)
		}
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(logWriter), logLevel[cfgInst.Zap.Level]),
		)
		tmpLog := zap.New(core, zap.AddCaller())
		loggers.SecureLogger = tmpLog.Sugar()
	}
	// dataLogger
	{
		tmpLog := logger.New(
			log.New(outputWriter, "\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second,   // 慢 SQL 阈值
				LogLevel:                  logger.Silent, // 日志级别
				IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,         // 禁用彩色打印
			},
		)
		loggers.GormLogger = &tmpLog
		//loggers.DataLogger = log.New(io.MultiWriter(os.Stderr, ))
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
	var (
		hook *rotatelogs.RotateLogs
		err  error
	)
	if runtime.GOOS == "windows" {
		hook, err = rotatelogs.New(
			filenamePrefix+"-%Y%m%d.log",
			rotatelogs.WithMaxAge(time.Duration(maxAgeDays*24)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(rotationTimeDays*24)*time.Hour),
		)
	} else {
		hook, err = rotatelogs.New(
			filenamePrefix+"-%Y%m%d.log",
			rotatelogs.WithLinkName(filenamePrefix+".log"),
			rotatelogs.WithMaxAge(time.Duration(maxAgeDays*24)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(rotationTimeDays*24)*time.Hour),
		)
	}
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Fatalf("cannot initialize SecureLogger: %v\n", err)
	}
	return hook
}
