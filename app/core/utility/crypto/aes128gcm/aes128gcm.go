package aes256gcm

import (
	"app/core/utility/crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"github.com/pkg/errors"
)

const (
	KeyBitLength  = 128
	KeyByteLength = KeyBitLength >> 3
)

var (
	nonceSize = -1
	HashFunc  = md5.New
)

func getNonceSize() (int, error) {
	if nonceSize <= 0 {
		fakeSecret := [KeyByteLength]byte{0}
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

type DataAES128GCM crypto.DataCrypto

func (d DataAES128GCM) MakeKeys() ([]byte, []byte, error) {
	if nil == d.SharedKey {
		return nil, nil, errors.Errorf("crypto.aes128gcm.MakeKeys: sharedKey is nil")
	}
	if nil == d.Nonce {
		return nil, nil, errors.Errorf("crypto.aes128gcm.MakeKeys: IV is nil")
	}
	he := hmac.New(HashFunc, d.Nonce)
	he.Write(d.SharedKey)
	encKey := he.Sum(nil)
	hm := hmac.New(HashFunc, d.SharedKey)
	hm.Write(d.Nonce)
	macKey := hm.Sum(nil)
	return encKey, macKey, nil
}

func (d DataAES128GCM) Encrypt() ([]byte, error) {
	if nil == d.SharedKey {
		return nil, errors.Errorf("crypto.aes128gcm.Encrypt: sharedKey is nil")
	}
	if nil == d.Data {
		return nil, errors.Errorf("crypto.aes128gcm.Encrypt: data is nil")
	}
	nonceSize, err := getNonceSize()
	if nil != err {
		return nil, errors.Errorf("crypto.aes128gcm.Encrypt: nonce size")
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
		return nil, errors.Errorf("crypto.aes128gcm.Encrypt: aes.NewCipher")
	}
	aesGCM, err := cipher.NewGCM(block)
	if nil != err {
		return nil, errors.Errorf("crypto.aes128gcm.Encrypt: cipher.NewGCM")
	}
	encWithTag := aesGCM.Seal(nil, d.Nonce, d.Data, macKey)
	encrypted := make([]byte, nonceSize+len(encWithTag))
	copy(encrypted, d.Nonce)
	copy(encrypted[nonceSize:], encWithTag)
	return encrypted, nil
}

func (d DataAES128GCM) Decrypt() ([]byte, error) {
	if nil == d.SharedKey {
		return nil, errors.Errorf("crypto.aes128gcm.Decrypt: sharedKey is nil")
	}
	if nil == d.Data {
		return nil, errors.Errorf("crypto.aes128gcm.Decrypt: data is nil")
	}
	nonceSize, err := getNonceSize()
	if nil != err {
		return nil, errors.Errorf("crypto.aes128gcm.Decrypt: nonce size")
	}
	d.Nonce = d.Data[:nonceSize]
	encWithTag := d.Data[nonceSize:]
	encKey, macKey, err := d.MakeKeys()
	if nil != err {
		return nil, err
	}
	block, err := aes.NewCipher(encKey)
	if nil != err {
		return nil, errors.Errorf("crypto.aes128gcm.Decrypt: aes.NewCipher")
	}
	aesGCM, err := cipher.NewGCM(block)
	if nil != err {
		return nil, errors.Errorf("crypto.aes128gcm.Decrypt: cipher.NewGCM")
	}
	message, err := aesGCM.Open(nil, d.Nonce, encWithTag, macKey)
	if nil != err {
		return nil, err
	}
	return message, nil
}
