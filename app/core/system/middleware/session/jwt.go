package session

import (
	"app/core/utility/basex/base36"
	"app/core/utility/crypto"
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rc4"
	"github.com/pkg/errors"
)

// 仅保持用户ID
// Set-Cookie: PHPSESSID=1dhpndum9ltrtpk8hnrqp5o3is; path=/

type PHPSessionId struct {
	UserIntId uint32
	StartTime uint32
}

type JWT string

func (s *PHPSessionId) Encode(sharedKey []byte) string {
	nonce := crypto.Rand(4, true)
	encKey, macKey, _ := crypto.MakeKeys(md5.New, sharedKey, nonce)
	src := make([]byte, 8)
	for i := 0; i < 4; i++ {
		src[i] = byte((s.UserIntId >> (8 * i)) & 0xff)
		src[i+4] = byte((s.StartTime >> (8 * i)) & 0xff)
	}
	dst := make([]byte, 17)
	copy(dst[5:], nonce)
	cipher, _ := rc4.NewCipher(encKey)
	cipher.XORKeyStream(dst[9:], src)
	h := hmac.New(md5.New, macKey)
	h.Write(dst[5:])
	mac := h.Sum(nil)
	copy(dst, mac[:5])
	return base36.EncodeToStringLc(dst)
}

func (s *JWT) Decode(sharedKey []byte) (*PHPSessionId, error) {
	stream, err := base36.DecodeString(string(*s))
	if nil != err {
		return nil, err
	}
	if len(stream) != 17 {
		return nil, errors.Errorf("jwt decode length error")
	}
	nonce := stream[5:9]
	encKey, macKey, _ := crypto.MakeKeys(md5.New, sharedKey, nonce)
	expectMac := stream[:5]
	h := hmac.New(md5.New, macKey)
	h.Write(stream[5:])
	realMac := h.Sum(nil)
	if !bytes.Equal(expectMac, realMac) {
		return nil, errors.Errorf("jwt verify failed")
	}
	buffer := make([]byte, 8)
	cipher, _ := rc4.NewCipher(encKey)
	cipher.XORKeyStream(buffer, stream[9:])

	session := PHPSessionId{0, 0}
	for i := 0; i < 4; i++ {
		session.UserIntId = (session.UserIntId << 8) | uint32(buffer[3-i])
		session.StartTime = (session.StartTime << 8) | uint32(buffer[7-i])
	}
	return &session, nil
}
