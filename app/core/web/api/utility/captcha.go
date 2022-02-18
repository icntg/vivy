package utility

import (
	"app/core/global/config"
	"app/core/global/logger"
	"app/core/utility/errno"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wenlng/go-captcha/captcha"
)

// NewCaptcha 生成图形校验码
// @Response Session[captcha] = 生成的
func NewCaptcha(c *gin.Context) {
	var (
		gCfg = config.Instance()
		oLog = logger.Instance().OutputLogger
	)
	session := sessions.Default(c)

	capt := captcha.GetCaptcha()
	dots, b64, tb64, key, err := capt.Generate()
	if nil != err {
		_ = errno.ErrorCaptchaGenerate.Fold()
		oLog.Errorf("")
		//c.JSON()
		return
	}
	// save DOTS and KEY to session
}

// VerifyCaptcha 校验
func VerifyCaptcha(c *gin.Context) {

}
