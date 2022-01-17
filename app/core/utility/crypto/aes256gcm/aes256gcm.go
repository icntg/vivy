package aes256gcm

import (
	"app/core/utility/crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"github.com/pkg/errors"
)

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
	fakeSecret := [64]byte{0}
	block, err := aes.NewCipher(fakeSecret[:])
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: aes.NewCipher0")
	}
	aesGCM, err := cipher.NewGCM(block)
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: cipher.NewGCM0")
	}
	n := aesGCM.NonceSize()
	if nil == iv || len(iv) < n {
		nonce = crypto.Rand(n, true)
	} else {
		nonce = iv[:n]
	}

	encKey, macKey, err := MakeKeys(sharedKey, iv)
	if nil != err {
		return nil, err
	}
	block, err = aes.NewCipher(encKey)
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: aes.NewCipher1")
	}
	aesGCM, err = cipher.NewGCM(block)
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: cipher.NewGCM1")
	}
	encryptedWithTag := aesGCM.Seal(nil, nonce, message, macKey)
	return encryptedWithTag, nil
}

func Decrypt(sharedKey, encrypted []byte) ([]byte, error) {
	if nil == sharedKey {
		return nil, errors.Errorf("crypto.aes256gcm.Decrypt: sharedKey is nil")
	}
	if nil == encrypted {
		return nil, errors.Errorf("crypto.aes256gcm.Decrypt: encrypted is nil")
	}
	// TODO:
	return nil, nil
}

func MakeKeys(sharedKey, iv []byte) ([]byte, []byte, error) {
	if nil == sharedKey {
		return nil, nil, errors.Errorf("crypto.rc4md5.MakeKeys: sharedKey is nil")
	}
	if nil == iv || len(iv) != crypto.IVLen {
		return nil, nil, errors.Errorf("crypto.rc4md5.MakeKeys: IV is incorrect")
	}
	he := hmac.New(sha256.New, iv)
	he.Write(sharedKey)
	encKey := he.Sum(nil)
	hm := hmac.New(sha256.New, sharedKey)
	hm.Write(iv)
	macKey := hm.Sum(nil)
	return encKey, macKey, nil
}
