// 说明：
// 未登录前，使用cookie加密保存session
// 登录之后，使用redis session

package session

import (
	"app/core/global/config"
	"app/core/utility/crypto/aes128gcm"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

// 为了保证安全性，凭证必须使用secure cookie。
// 其他数据可用local_storage或session_storage保存

const (
	SecureCookieName = "_s5n_"
	InContextName    = "SESSION"
	TimestampKey     = "TIMESTAMP"
)

func decode(encB64Cookie string) (map[string]interface{}, error) {
	var (
		secret    = config.Instance().Service.SessionSecretBytes
		encCookie []byte
		err       error
		message   []byte
	)
	if encCookie, err = base64.URLEncoding.DecodeString(encB64Cookie); nil != err {
		return nil, err
	}
	enc := aes128gcm.Encrypted(encCookie)
	if message, err = enc.DecryptFromStream(secret); nil != err {
		return nil, err
	}
	ret := make(map[string]interface{})
	if err = json.Unmarshal(message, &ret); nil != err {
		return nil, err
	}
	return ret, nil
}

func GetSessionInCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			session       map[string]interface{}
			encB64Session string
			tmpJson       map[string]interface{}
			tmpInterface  interface{}
			err           error
		)
		// 首先从cookie中提取
		if encB64Session, err = c.Cookie(SecureCookieName); nil != err {
			// log
			goto TryJSON
		}
		if session, err = decode(encB64Session); nil != err {
			// log
			goto TryJSON
		}
		goto SUCCESS
		// 其次从JSON参数中提取
	TryJSON:
		tmpJson = make(map[string]interface{})
		if err = c.BindJSON(&tmpJson); nil != err {
			// log
			goto TryPostForm
		}
		tmpInterface = tmpJson[SecureCookieName]
		switch tmpInterface.(type) {
		case string:
			encB64Session = tmpInterface.(string)
		default:
			// log
			goto TryPostForm
		}
		if session, err = decode(encB64Session); nil != err {
			// log
			goto TryPostForm
		}
		goto SUCCESS
		// 最后从PostForm中提取
	TryPostForm:
		if encB64Session = c.PostForm(SecureCookieName); len(encB64Session) == 0 {
			goto EXIT
		}
		if session, err = decode(encB64Session); nil != err {
			goto EXIT
		}
	SUCCESS:
		// TODO: check session[TimestampKey]
		_ = session[TimestampKey]
		c.Set(InContextName, session)
	EXIT:
		c.Next()
	}
}

func SetSessionInCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			secret              = config.Instance().Service.SessionSecretBytes
			timeout             = config.Instance().Service.CookieTimeout
			tmpSessionInterface interface{}
			tmpSession          map[string]interface{}
			sessionExist        bool
			clearSession        aes128gcm.Message
			encSession          aes128gcm.Encrypted
			encB64Session       string
			err                 error
		)
		if tmpSessionInterface, sessionExist = c.Get(InContextName); !sessionExist {
			// TODO: log
			c.SetCookie(SecureCookieName, "", -1, "/", "", false, false)
			goto EXIT
		}
		switch tmpSessionInterface.(type) {
		case map[string]interface{}:
			tmpSession = tmpSessionInterface.(map[string]interface{})
		default:
			// log
			goto EXIT
		}
		tmpSession[TimestampKey] = time.Now().UTC().Unix()
		if clearSession, err = json.Marshal(tmpSession); nil != err {
			// log
			goto EXIT
		}
		if encSession, err = clearSession.EncryptToStream(secret, nil); nil != err {
			// log
			goto EXIT
		}
		encB64Session = base64.URLEncoding.EncodeToString(encSession)
		c.SetCookie(SecureCookieName, encB64Session, timeout, "/", "", false, false)
	EXIT:
		c.Next()
	}
}

//
//type Session struct {
//	StartTimeUnixStamp int64                  `json:"start_time"`
//	UserId             string                 `json:"user_id"`
//	Data               map[string]interface{} `json:"data"`
//}
//
//type EncodedSession string
//
//func New(userId string) Session {
//	return Session{
//		StartTimeUnixStamp: time.Now().UTC().Unix(),
//		UserId:             userId,
//		Data:               make(map[string]interface{}),
//	}
//}
//
//func (session *Session) encodeToString() string {
//	j, err := json.Marshal(session)
//	if nil != err {
//		err = errors.Wrap(err, "session encode failed")
//		// todo: log error
//	}
//	msg := aes128gcm.Message(j)
//	enc, err := msg.EncryptToStream(sessionSecret, nil)
//	if nil != err {
//		err = errors.Wrap(err, "encrypt failed")
//		// todo: log error
//	}
//	return base64.URLEncoding.EncodeToString(enc)
//}
//
//func (session *Session) decodeFromString(encSession string) error {
//	enc0, err := base64.URLEncoding.DecodeString(encSession)
//	if nil != err {
//		err = errors.Wrap(err, "base64 decode failed")
//		// todo: log err
//		return err
//	}
//	enc1 := aes128gcm.Encrypted(enc0)
//	dec, err := enc1.DecryptFromStream(sessionSecret)
//	if nil != err {
//		err = errors.Wrap(err, "decrypt failed")
//		// todo: log err
//		return err
//	}
//	err = json.Unmarshal(dec, session)
//	if nil != err {
//		err = errors.Wrap(err, "unmarshal failed")
//		// todo: log err
//		return err
//	}
//	return nil
//}
//
//func GetSession(name string, ctx *gin.Context) (*Session, error) {
//	encSession := EncodedSession(ctx.Param(name))
//	if len(encSession) == 0 {
//		return nil, nil
//	}
//
//	return nil, nil
//}
//
//func SetSession(name string, session *Session, ctx *gin.Context) error {
//	return nil
//}
