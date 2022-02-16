package errno

import (
	"app/core/utility/common"
	"os"
	"strconv"
	"strings"
)

type errCodeMsg string

const (
	// 负数为系统错误
	// 0 无错误
	// 0x0000xxxx 为基础组件错误
	// 0x0001xxxx 为业务逻辑错误

	ErrorSuccess errCodeMsg = "0x00000000 | ErrorSuccess | 无错误"

	ErrorSystemArgs      errCodeMsg = "0x00000010 | SystemArgsError | 命令行参数错误"
	ErrorSystemArgsIsNil errCodeMsg = "0x00000011 | SystemArgs is nil | 命令行参数为空"

	ErrorConfig          errCodeMsg = "0x00000020 | ConfigError | 配置文件错误"
	ErrorConfigIsNil     errCodeMsg = "0x00000021 | Config is nil | 配置信息为空"
	ErrorConfigNotExist  errCodeMsg = "0x00000022 | Config File does NOT exist | 配置文件不存在"
	ErrorConfigRead      errCodeMsg = "0x00000023 | Config File Read error | 配置文件读取错误"
	ErrorConfigUnmarshal errCodeMsg = "0x00000024 | Config File Unmarshal error | 配置文件解析错误"
	ErrorConfigMarshal   errCodeMsg = "0x00000025 | Config Marshal error | 配置文件生成yaml数据字节流错误"
	ErrorConfigViperRead errCodeMsg = "0x00000026 | Config Viper Read error | 配置文件yaml字节流读入viper错误"
	ErrorConfigWrite     errCodeMsg = "0x00000027 | Config File Write error | 配置文件写入错误"
	ErrorConfigDecodeHex errCodeMsg = "0x00000028 | Config File Decode hex error | 配置文件十六进制字符串解码错误"

	ErrorLogger        errCodeMsg = "0x00000030 | LoggerError | 日志错误"
	ErrorLoggerMakeDir errCodeMsg = "0x00000031 | Logger cannot create Log Directory | 日志无法创建目录"
	ErrorLoggerInit    errCodeMsg = "0x00000032 | Logger Init error | 日志初始化错误"

	ErrorGORM          errCodeMsg = "0x00000040 | GORMError | MySQL错误"
	ErrorGORMOpen      errCodeMsg = "0x00000041 | GORM cannot open database | GORM打开数据库错误"
	ErrorGORMDBHandler errCodeMsg = "0x00000042 | GORM cannot get database handler | GORM获取数据库handler错误"
	ErrorGORMDBStats   errCodeMsg = "0x00000043 | GORM cannot get database stats | GORM获取数据库状态错误"
	ErrorSQLOpen       errCodeMsg = "0x00000044 | SQL cannot open database | SQL打开数据库错误"
	ErrorSQLExec       errCodeMsg = "0x00000045 | SQL Connection cannot exec sql | SQL执行语句错误"

	ErrorRedis      errCodeMsg = "0x00000050 | RedisError | Redis错误"
	ErrorRedisStore errCodeMsg = "0x00000051 | Redis New Store Error | Redis连接错误"

	//ErrorStartGinService                     = 0x1010
	//ErrorConnectDatabase                     = 0x1020
)

func (e errCodeMsg) Code() int {
	a := strings.Split(string(e), "|")
	c, _ := strconv.ParseInt(strings.TrimSpace(a[0]), 16, 64)
	code := int(c)
	return code
}

func (e errCodeMsg) Exit() {
	a := strings.Split(string(e), "|")
	c, _ := strconv.ParseInt(strings.TrimSpace(a[0]), 16, 64)
	code := int(c)
	msg := strings.TrimSpace(a[1])
	common.ErrPrintf(msg + "\n")
	os.Exit(code)
}
