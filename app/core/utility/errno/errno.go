package errno

import (
	"app/core/utility/common"
	"os"
	"strconv"
	"strings"
)

type errorCodeMsg string

const (
	ErrorSystemArgs      errorCodeMsg = "0x00010000 | SystemArgsError | 命令行参数错误"
	ErrorSystemArgsIsNil errorCodeMsg = "0x00010001 | SystemArgs is nil |命令行参数为空"

	ErrorConfig         errorCodeMsg = "0x00010010 | ConfigError | 配置文件错误"
	ErrorConfigNotExist errorCodeMsg = "0x00010010 | Config File does NOT exist | 配置文件不存在"

	//ErrorGenerateConfigTemplate              = 0x2001
	//ErrorReadConfig                          = 0x2002
	//ErrorInitLogger                          = 0x3003
	//ErrorStartGinService                     = 0x1010
	//ErrorConnectDatabase                     = 0x1020
)

type Exception struct {
	Code int
	Err  error
}

func (e errorCodeMsg) Code() int {
	a := strings.Split(string(e), "|")
	c, _ := strconv.ParseInt(strings.TrimSpace(a[0]), 16, 64)
	code := int(c)
	return code
}

func (e errorCodeMsg) Exit() {
	a := strings.Split(string(e), "|")
	c, _ := strconv.ParseInt(strings.TrimSpace(a[0]), 16, 64)
	code := int(c)
	msg := strings.TrimSpace(a[1])
	common.ErrPrintf(msg + "\n")
	os.Exit(code)
}
