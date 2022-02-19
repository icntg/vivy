// 说明：
// 未登录前，使用cookie加密保存session
// 登录之后，使用redis session

package session_in_cookie

import (
	"app/core/global/config"
	"app/core/global/logger"
	"app/core/utility/crypto/aes128gcm"
	"app/core/utility/errno"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"
)

// 为了保证安全性，凭证必须使用secure cookie。
// 其他数据可用local_storage或session_storage保存

const (
	SecureCookieName = "_"
	TimestampKey     = "TIMESTAMP"
)

func EncodeSecureCookie(cookieMap map[string]interface{}) (string, error) {
	var (
		oLog                              = logger.Instance().OutputLogger
		secret                            = config.Instance().Service.SessionSecretInBytes
		clearSession  aes128gcm.Message   = nil
		encSession    aes128gcm.Encrypted = nil
		b64EncSession string
		err           error = nil
	)
	if cookieMap == nil || len(cookieMap) == 0 {
		oLog.Debugf("cookieMap [%v] is empty\n", cookieMap)
		return "", errno.ErrorParamIsNil.Fold().Error
	}
	cookieMap[TimestampKey] = strconv.FormatInt(time.Now().UTC().Unix(), 16)
	if clearSession, err = json.Marshal(cookieMap); nil != err {
		oLog.Debugf("jsonify cookieMap [%v] failed: %v\n", cookieMap, err)
		return "", err
	}
	if encSession, err = clearSession.EncryptToStream(secret, nil); nil != err {
		oLog.Debugf("clearSession [%v] encrypt failed: %v\n", clearSession, err)
		return "", err
	}
	b64EncSession = base64.URLEncoding.EncodeToString(encSession)
	if len(b64EncSession) > 1024*3.5 {
		oLog.Debugf("b64EncSession [%v] is too long: %v\n", b64EncSession, len(b64EncSession))
		return "", errno.ErrorSecureCookieTooLarge.Fold().Error
	}
	return b64EncSession, nil
}

func DecodeSecureCookie(b64EncSession string) (map[string]interface{}, error) {
	var (
		oLog         = logger.Instance().OutputLogger
		secret       = config.Instance().Service.SessionSecretInBytes
		timeout      = config.Instance().Service.CookieTimeout
		timestamp    int64
		clearSession aes128gcm.Message   = nil
		encSession   aes128gcm.Encrypted = nil
		err          error               = nil
	)
	if encSession, err = base64.URLEncoding.DecodeString(b64EncSession); nil != err {
		oLog.Debugf("b64EncSession [%v] decode failed: %v\n", b64EncSession, err)
		return nil, err
	}
	if clearSession, err = encSession.DecryptFromStream(secret); nil != err {
		oLog.Debugf("encSession [%v] decrypt failed: %v\n", encSession, err)
		return nil, err
	}
	session := make(map[string]interface{})
	if err = json.Unmarshal(clearSession, &session); nil != err {
		oLog.Debugf("clearSession [%v] unmarshal to map failed: %v\n", clearSession, err)
		return nil, err
	}
	if _, ok := session[TimestampKey]; !ok {
		oLog.Debugf("key [%s] not in session [%v]\n", TimestampKey, session)
		return nil, errno.ErrorSecureCookieMissTimeData.Fold().Error
	}
	if timestamp, err = strconv.ParseInt(session[TimestampKey].(string), 16, 64); nil != err {
		oLog.Debugf("cannot convert key [%s][%v] in session [%v] to int\n", TimestampKey, session[TimestampKey], session)
		return nil, err
	}
	now := time.Now().UTC().Unix()
	if now-timestamp > int64(timeout) {
		oLog.Debugf("session timeout now[%v] - sessionTimestamp[%v] > timeoutValue[%v]\n", now, timestamp, timeout)
		return nil, errno.ErrorSecureCookieTimeout.Fold().Error
	}
	delete(session, TimestampKey)
	return session, nil
}
