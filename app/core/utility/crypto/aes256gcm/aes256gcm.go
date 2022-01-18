package aes256gcm

import (
	"app/core/utility/crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"github.com/pkg/errors"
)

var (
	nonceSize int = -1
)

func getNonceSize() (int, error) {
	if nonceSize <= 0 {
		fakeSecret := [32]byte{0}
		block, err := aes.NewCipher(fakeSecret[:])
		if nil != err {
			return -1, err
		}
		aesGCM, err := cipher.NewGCM(block)
		if nil != err {
			return -1, err
		}
		nonceSize = aesGCM.NonceSize()
	}
	return nonceSize, nil
}

func Encrypt(sharedKey []byte, iv []byte, message []byte) ([]byte, error) {
	var (
		nonce []byte = nil
	)
	if nil == sharedKey {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: sharedKey is nil")
	}
	if nil == message {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: message is nil")
	}
	nonceSize, err := getNonceSize()
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: nonce size")
	}
	if nil == iv || len(iv) < nonceSize {
		nonce = crypto.Rand(nonceSize, true)
	} else {
		nonce = iv[:nonceSize]
	}

	encKey, macKey, err := MakeKeys(sharedKey, nonce)
	if nil != err {
		return nil, err
	}
	block, err := aes.NewCipher(encKey)
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: aes.NewCipher")
	}
	aesGCM, err := cipher.NewGCM(block)
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: cipher.NewGCM")
	}
	encWithTag := aesGCM.Seal(nil, nonce, message, macKey)
	encrypted := make([]byte, nonceSize+len(encWithTag))
	copy(encrypted, nonce)
	copy(encrypted[nonceSize:], encWithTag)
	return encrypted, nil
}

func Decrypt(sharedKey, encrypted []byte) ([]byte, error) {
	if nil == sharedKey {
		return nil, errors.Errorf("crypto.aes256gcm.Decrypt: sharedKey is nil")
	}
	if nil == encrypted {
		return nil, errors.Errorf("crypto.aes256gcm.Decrypt: encrypted is nil")
	}
	nonceSize, err := getNonceSize()
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Decrypt: nonce size")
	}
	nonce := encrypted[:nonceSize]
	encWithTag := encrypted[nonceSize:]
	encKey, macKey, err := MakeKeys(sharedKey, nonce)
	if nil != err {
		return nil, err
	}
	block, err := aes.NewCipher(encKey)
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Decrypt: aes.NewCipher")
	}
	aesGCM, err := cipher.NewGCM(block)
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Decrypt: cipher.NewGCM")
	}
	message, err := aesGCM.Open(nil, nonce, encWithTag, macKey)
	if nil != err {
		return nil, err
	}
	return message, nil
}

func MakeKeys(sharedKey, iv []byte) ([]byte, []byte, error) {
	if nil == sharedKey {
		return nil, nil, errors.Errorf("crypto.rc4md5.MakeKeys: sharedKey is nil")
	}
	if nil == iv {
		return nil, nil, errors.Errorf("crypto.rc4md5.MakeKeys: IV is nil")
	}
	he := hmac.New(sha256.New, iv)
	he.Write(sharedKey)
	encKey := he.Sum(nil)
	hm := hmac.New(sha256.New, sharedKey)
	hm.Write(iv)
	macKey := hm.Sum(nil)
	return encKey, macKey, nil
}
