package session

import (
	"app/core/utility/crypto/aes128gcm"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

// 为了保证安全性，凭证必须使用secure cookie。
// 其他数据可用local_storage或session_storage保存

var (
	sessionSecret []byte = nil
)

type Session struct {
	StartTimeUnixStamp int64                  `json:"start_time"`
	UserId             string                 `json:"user_id"`
	Data               map[string]interface{} `json:"data"`
}

type EncodedSession string

func New(userId string) Session {
	return Session{
		StartTimeUnixStamp: time.Now().UTC().Unix(),
		UserId:             userId,
		Data:               make(map[string]interface{}),
	}
}

func (session *Session) encodeToString() string {
	j, err := json.Marshal(session)
	if nil != err {
		err = errors.Wrap(err, "session encode failed")
		// todo: log error
	}
	msg := aes128gcm.Message(j)
	enc, err := msg.EncryptToStream(sessionSecret, nil)
	if nil != err {
		err = errors.Wrap(err, "encrypt failed")
		// todo: log error
	}
	return base64.URLEncoding.EncodeToString(enc)
}

func (session *Session) decodeFromString(encSession string) error {
	enc0, err := base64.URLEncoding.DecodeString(encSession)
	if nil != err {
		err = errors.Wrap(err, "base64 decode failed")
		// todo: log err
		return err
	}
	enc1 := aes128gcm.Encrypted(enc0)
	dec, err := enc1.DecryptFromStream(sessionSecret)
	if nil != err {
		err = errors.Wrap(err, "decrypt failed")
		// todo: log err
		return err
	}
	err = json.Unmarshal(dec, session)
	if nil != err {
		err = errors.Wrap(err, "unmarshal failed")
		// todo: log err
		return err
	}
	return nil
}

func GetSession(name string, ctx *gin.Context) (*Session, error) {
	encSession := EncodedSession(ctx.Param(name))
	if len(encSession) == 0 {
		return nil, nil
	}

	return nil, nil
}

func SetSession(name string, session *Session, ctx *gin.Context) error {
	return nil
}
