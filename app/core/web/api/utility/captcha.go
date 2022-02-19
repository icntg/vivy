package utility

import (
	"app/core/global"
	"app/core/global/config"
	"app/core/global/logger"
	"app/core/system/middleware/session_in_cookie"
	"app/core/utility/errno"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wenlng/go-captcha/captcha"
	"strconv"
	"strings"
)

//NewCaptcha 生成图形校验码
//@Response Session[captcha] = 生成的
func NewCaptcha(c *gin.Context) {
	var (
		timeout = config.Instance().Service.CookieTimeout
		oLog    = logger.Instance().OutputLogger
	)
	capt := captcha.GetCaptcha()
	dots, b64, tb64, _, err := capt.Generate()
	if nil != err {
		oLog.Errorf("ErrorCaptchaGenerate: %v\n", err)
		e := errno.ErrorCaptchaGenerate.Fold()
		errMsg := map[string]interface{}{
			"code":    e.Code,
			"message": e.Error.Error(),
		}
		c.JSON(errno.HttpInternalServerError, errMsg)
		return
	}
	// save DOTS and KEY to session
	session := make(map[string]interface{})
	session["dots"] = dots
	encSession, err := session_in_cookie.EncodeSecureCookie(session)
	if nil != err {
		oLog.Errorf("EncodeSecureCookieError: %v\n", err)
		e := errno.ErrorSecureCookie.Fold()
		errMsg := map[string]interface{}{
			"code":    e.Code,
			"message": e.Error.Error(),
		}
		c.JSON(errno.HttpInternalServerError, errMsg)
		return
	}
	c.SetCookie(
		session_in_cookie.SecureCookieName,
		encSession,
		timeout,
		"/",
		global.ProductName,
		false,
		true,
	)
	bt := map[string]interface{}{
		"code":         errno.NoError,
		"image_base64": b64,
		"thumb_base64": tb64,
	}
	c.JSON(errno.HttpOK, bt)
}

// VerifyCaptcha 校验
func VerifyCaptcha(dotsInput, dotsExpect string) bool {
	if len(dotsInput) == 0 || len(dotsExpect) == 0 {
		return false
	}
	input := strings.Split(dotsInput, ",")
	var dct map[int]captcha.CharDot
	if err := json.Unmarshal([]byte(dotsExpect), &dct); err != nil {
		return false
	}
	chkRet := false
	if (len(dct) * 2) == len(input) {
		for i, dot := range dct {
			j := i * 2
			k := i*2 + 1
			sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", input[j]), 64)
			sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", input[k]), 64)

			// 检测点位置
			// chkRet = captcha.CheckPointDist(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height))

			// 校验点的位置,在原有的区域上添加额外边距进行扩张计算区域,不推荐设置过大的padding
			// 例如：文本的宽和高为30，校验范围x为10-40，y为15-45，此时扩充5像素后校验范围宽和高为40，则校验范围x为5-45，位置y为10-50
			chkRet = captcha.CheckPointDistWithPadding(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height), 5)
			if !chkRet {
				break
			}
		}
	}
	return chkRet
}

//func VerifyCaptcha1(c *gin.Context) {
//	code := 1
//	_ = r.ParseForm()
//	dots := r.Form.Get("dots")
//	key := r.Form.Get("key")
//	if dots == "" || key == "" {
//		bt, _ := json.Marshal(map[string]interface{}{
//			"code":    code,
//			"message": "dots or key param is empty",
//		})
//		_, _ = fmt.Fprintf(w, string(bt))
//		return
//	}
//
//	cacheData := readCache(key)
//	if cacheData == "" {
//		bt, _ := json.Marshal(map[string]interface{}{
//			"code":    code,
//			"message": "illegal key",
//		})
//		_, _ = fmt.Fprintf(w, string(bt))
//		return
//	}
//	src := strings.Split(dots, ",")
//
//	var dct map[int]captcha.CharDot
//	if err := json.Unmarshal([]byte(cacheData), &dct); err != nil {
//		bt, _ := json.Marshal(map[string]interface{}{
//			"code":    code,
//			"message": "illegal key",
//		})
//		_, _ = fmt.Fprintf(w, string(bt))
//		return
//	}
//
//	chkRet := false
//	if (len(dct) * 2) == len(src) {
//		for i, dot := range dct {
//			j := i * 2
//			k := i*2 + 1
//			sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[j]), 64)
//			sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[k]), 64)
//
//			// 检测点位置
//			// chkRet = captcha.CheckPointDist(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height))
//
//			// 校验点的位置,在原有的区域上添加额外边距进行扩张计算区域,不推荐设置过大的padding
//			// 例如：文本的宽和高为30，校验范围x为10-40，y为15-45，此时扩充5像素后校验范围宽和高为40，则校验范围x为5-45，位置y为10-50
//			chkRet = captcha.CheckPointDistWithPadding(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height), 5)
//			if !chkRet {
//				break
//			}
//		}
//	}
//
//	if chkRet {
//		// 通过校验
//		code = 0
//	}
//
//	bt, _ := json.Marshal(map[string]interface{}{
//		"code": code,
//	})
//	_, _ = fmt.Fprintf(w, string(bt))
//	return
//}
