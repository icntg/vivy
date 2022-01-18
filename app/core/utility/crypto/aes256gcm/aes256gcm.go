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
	nonceSize = -1
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

type DataAES256GCM crypto.DataCrypto

func (d DataAES256GCM) MakeKeys() ([]byte, []byte, error) {
	if nil == d.SharedKey {
		return nil, nil, errors.Errorf("crypto.rc4md5.MakeKeys: sharedKey is nil")
	}
	if nil == d.Nonce {
		return nil, nil, errors.Errorf("crypto.rc4md5.MakeKeys: IV is nil")
	}
	he := hmac.New(sha256.New, d.Nonce)
	he.Write(d.SharedKey)
	encKey := he.Sum(nil)
	hm := hmac.New(sha256.New, d.SharedKey)
	hm.Write(d.Nonce)
	macKey := hm.Sum(nil)
	return encKey, macKey, nil
}

func (d DataAES256GCM) Encrypt() ([]byte, error) {
	if nil == d.SharedKey {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: sharedKey is nil")
	}
	if nil == d.Data {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: data is nil")
	}
	nonceSize, err := getNonceSize()
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Encrypt: nonce size")
	}
	if nil == d.Nonce || len(d.Nonce) < nonceSize {
		d.Nonce = crypto.Rand(nonceSize, true)
	} else {
		d.Nonce = d.Nonce[:nonceSize]
	}

	encKey, macKey, err := d.MakeKeys()
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
	encWithTag := aesGCM.Seal(nil, d.Nonce, d.Data, macKey)
	encrypted := make([]byte, nonceSize+len(encWithTag))
	copy(encrypted, d.Nonce)
	copy(encrypted[nonceSize:], encWithTag)
	return encrypted, nil
}

func (d DataAES256GCM) Decrypt() ([]byte, error) {
	if nil == d.SharedKey {
		return nil, errors.Errorf("crypto.aes256gcm.Decrypt: sharedKey is nil")
	}
	if nil == d.Data {
		return nil, errors.Errorf("crypto.aes256gcm.Decrypt: data is nil")
	}
	nonceSize, err := getNonceSize()
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.Decrypt: nonce size")
	}
	d.Nonce = d.Data[:nonceSize]
	encWithTag := d.Data[nonceSize:]
	encKey, macKey, err := d.MakeKeys()
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
	message, err := aesGCM.Open(nil, d.Nonce, encWithTag, macKey)
	if nil != err {
		return nil, err
	}
	return message, nil
}
