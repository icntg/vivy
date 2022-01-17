package rc4md5

import (
	"app/core/utility/crypto"
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rc4"
	"github.com/pkg/errors"
)

func Encrypt(sharedKey []byte, iv []byte, message []byte) ([]byte, error) {
	var (
		iv8 [8]byte
	)
	if nil == sharedKey {
		return nil, errors.Errorf("crypto.rc4md5.Encrypt: sharedKey is nil")
	}
	if nil == message {
		return nil, errors.Errorf("crypto.rc4md5.Encrypt: message is nil")
	}
	if nil == iv || len(iv) < crypto.IVLen {
		tmpIV := crypto.Rand(crypto.IVLen, true)
		copy(iv8[:], tmpIV)
	} else {
		copy(iv8[:], iv)
	}
	encKey, macKey, err := MakeKeys(sharedKey, iv)
	if nil != err {
		return nil, err
	}
	cipher, err := rc4.NewCipher(encKey)
	if nil != err {
		return nil, err
	}
	buffer := make([]byte, crypto.SumLen+crypto.IVLen+len(message))
	copy(buffer[crypto.SumLen:], iv8[:])
	cipher.XORKeyStream(buffer[crypto.SumLen+crypto.IVLen:], message)
	h := hmac.New(md5.New, macKey)
	h.Write(buffer[crypto.SumLen:])
	mac := h.Sum(nil)
	copy(buffer, mac)
	return buffer, nil
}

func Decrypt(sharedKey, encrypted []byte) ([]byte, error) {
	if nil == sharedKey {
		return nil, errors.Errorf("crypto.rc4md5.Decrypt: sharedKey is nil")
	}
	if nil == encrypted {
		return nil, errors.Errorf("crypto.rc4md5.Decrypt: encrypted is nil")
	}
	if len(encrypted) <= crypto.SumLen+crypto.IVLen+1 {
		return nil, errors.Errorf("crypto.rc4md5.Decrypt: encrypted is too short")
	}
	expectMac := encrypted[0:crypto.SumLen]
	iv := encrypted[crypto.SumLen : crypto.SumLen+crypto.IVLen]
	encKey, macKey, err := MakeKeys(sharedKey, iv)
	if nil != err {
		return nil, err
	}
	h := hmac.New(md5.New, macKey)
	h.Write(encrypted[crypto.SumLen:])
	realMac := h.Sum(nil)
	if !bytes.Equal(expectMac, realMac) {
		return nil, errors.Errorf("crypto.rc4md5.Decrypt: MAC verify failed")
	}

	cipher, err := rc4.NewCipher(encKey)
	if nil != err {
		return nil, err
	}
	buffer := make([]byte, len(encrypted)-crypto.SumLen-crypto.IVLen)
	cipher.XORKeyStream(buffer, encrypted[crypto.SumLen+crypto.IVLen:])
	return buffer, nil
}

func MakeKeys(sharedKey, iv []byte) ([]byte, []byte, error) {
	if nil == sharedKey {
		return nil, nil, errors.Errorf("crypto.rc4md5.MakeKeys: sharedKey is nil")
	}
	if nil == iv || len(iv) != crypto.IVLen {
		return nil, nil, errors.Errorf("crypto.rc4md5.MakeKeys: IV is incorrect")
	}
	he := hmac.New(md5.New, iv)
	he.Write(sharedKey)
	encKey := he.Sum(nil)
	hm := hmac.New(md5.New, sharedKey)
	hm.Write(iv)
	macKey := hm.Sum(nil)
	return encKey, macKey, nil
}
