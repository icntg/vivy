//package logger
//
//import (
//	"app/core/config"
//	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
//	"go.uber.org/zap"
//	"go.uber.org/zap/zapcore"
//	"io"
//	"log"
//	"os"
//	"path"
//	"time"
//)
//
//const (
//	outputFilename = "output"
//	accessFilename = "access"
//	secureFilename = "secure"
//)
//
//var (
//	loglevel = map[string]zapcore.Level{
//		"debug": zapcore.DebugLevel,
//		"info":  zapcore.InfoLevel,
//		"warn":  zapcore.WarnLevel,
//		"error": zapcore.ErrorLevel,
//	}
//)
//
//type Loggers struct {
//	OutPutLogger *zap.SugaredLogger
//	AccessLogger *zap.Logger
//	SecureLogger *zap.SugaredLogger
//}
//
//func Init(cfg *config.Config) Loggers {
//	var (
//		err    error
//		tmpLog *zap.Logger
//	)
//	logPath := path.Join(cfg.Zap.Director)
//
//	// 由于AccessLogger内容较多，而且格式固定，故采用标准Logger
//	ths.AccessLogger, err = zap.NewProduction()
//	if nil != err {
//		log.SetOutput(os.Stderr)
//		log.Fatalf("cannot initialize AccessLogger: %v\n", err)
//	}
//
//	// 其他Logger格式不定，使用SugaredLogger
//	tmpLog, err = zap.NewProduction()
//	if nil != err {
//		log.SetOutput(os.Stderr)
//		log.Fatalf("cannot initialize OutPutLogger: %v\n", err)
//	}
//	ths.OutPutLogger = tmpLog.Sugar()
//
//	tmpLog, err = zap.NewProduction()
//	if nil != err {
//		log.SetOutput(os.Stderr)
//		log.Fatalf("cannot initialize SecureLogger: %v\n", err)
//	}
//	ths.SecureLogger = tmpLog.Sugar()
//}
//
////
//func initLogger(logFilePath string) (interface{}, error) {
//
//}
//
////初始化Encoder
//func initEncoder() zapcore.Encoder {
//	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
//		MessageKey:  "msg",
//		LevelKey:    "level",
//		TimeKey:     "time",
//		CallerKey:   "file",
//		EncodeLevel: zapcore.CapitalLevelEncoder, //基本zapcore.LowercaseLevelEncoder。将日志级别字符串转化为小写
//		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
//			enc.AppendString(t.Format("2006-01-02 15:04:05"))
//		},
//		EncodeCaller: zapcore.ShortCallerEncoder, //一般zapcore.ShortCallerEncoder，以包/文件:行号 格式化调用堆栈
//		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) { //一般zapcore.SecondsDurationEncoder,执行消耗的时间转化成浮点型的秒
//			enc.AppendInt64(int64(d) / 1000000)
//		},
//	})
//}
//
////日志文件切割
//func getWriter(filenamePrefix string, maxAgeDays int, rotationTimeDays int) io.Writer {
//	// 保存30天内的日志，每24小时(整点)分割一次日志
//	hook, err := rotatelogs.New(
//		filenamePrefix+"-%Y%m%d.log",
//		rotatelogs.WithLinkName(filenamePrefix+".log"),
//		rotatelogs.WithMaxAge(time.Duration(maxAgeDays*24)*time.Hour),
//		rotatelogs.WithRotationTime(time.Duration(rotationTimeDays*24)*time.Hour),
//	)
//	if err != nil {
//		log.SetOutput(os.Stderr)
//		log.Fatalf("cannot initialize SecureLogger: %v\n", err)
//	}
//	return hook
//}
