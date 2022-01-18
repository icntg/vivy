package rc4md5

import (
	"app/core/utility/crypto"
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rc4"
	"github.com/pkg/errors"
)

const (
	MACLen = 16
	IVLen  = 8
)

// [MAC][IV][Encrypted]
// encKey = hmac_md5(iv, key)
// macKey = hmac_md5(key, iv)
// Encrypted = rc4(encKey, message)
// MAC = hmac_md5(macKey, iv + message)

type DataRC4MD5 crypto.DataCrypto

func (d DataRC4MD5) MakeKeys() ([]byte, []byte, error) {
	if nil == d.SharedKey {
		return nil, nil, errors.Errorf("crypto.rc4md5.MakeKeys: sharedKey is nil")
	}
	if nil == d.Nonce || len(d.Nonce) != IVLen {
		return nil, nil, errors.Errorf("crypto.rc4md5.MakeKeys: IV is incorrect")
	}
	he := hmac.New(md5.New, d.Nonce)
	he.Write(d.SharedKey)
	encKey := he.Sum(nil)
	hm := hmac.New(md5.New, d.SharedKey)
	hm.Write(d.Nonce)
	macKey := hm.Sum(nil)
	return encKey, macKey, nil
}

func (d DataRC4MD5) Encrypt() ([]byte, error) {
	if nil == d.SharedKey {
		return nil, errors.Errorf("crypto.rc4md5.Encrypt: sharedKey is nil")
	}
	if nil == d.Data {
		return nil, errors.Errorf("crypto.rc4md5.Encrypt: data is nil")
	}
	if nil == d.Nonce || len(d.Nonce) < IVLen {
		d.Nonce = crypto.Rand(IVLen, true)
	} else {
		d.Nonce = d.Nonce[:IVLen]
	}
	encKey, macKey, err := d.MakeKeys()
	if nil != err {
		return nil, err
	}
	cipher, err := rc4.NewCipher(encKey)
	if nil != err {
		return nil, err
	}
	buffer := make([]byte, MACLen+IVLen+len(d.Data))
	copy(buffer[MACLen:], d.Nonce)
	cipher.XORKeyStream(buffer[MACLen+IVLen:], d.Data)
	h := hmac.New(md5.New, macKey)
	h.Write(buffer[MACLen:])
	mac := h.Sum(nil)
	copy(buffer, mac)
	return buffer, nil
}

func (d DataRC4MD5) Decrypt() ([]byte, error) {
	if nil == d.SharedKey {
		return nil, errors.Errorf("crypto.rc4md5.Decrypt: sharedKey is nil")
	}
	if nil == d.Data {
		return nil, errors.Errorf("crypto.rc4md5.Decrypt: encrypted is nil")
	}
	if len(d.Data) <= MACLen+IVLen+1 {
		return nil, errors.Errorf("crypto.rc4md5.Decrypt: encrypted is too short")
	}
	expectMac := d.Data[0:MACLen]
	d.Nonce = d.Data[MACLen : MACLen+IVLen]
	encKey, macKey, err := d.MakeKeys()
	if nil != err {
		return nil, err
	}
	h := hmac.New(md5.New, macKey)
	h.Write(d.Data[MACLen:])
	calcMac := h.Sum(nil)
	if !bytes.Equal(expectMac, calcMac) {
		return nil, errors.Errorf("crypto.rc4md5.Decrypt: MAC verify failed")
	}

	cipher, err := rc4.NewCipher(encKey)
	if nil != err {
		return nil, err
	}
	buffer := make([]byte, len(d.Data)-MACLen-IVLen)
	cipher.XORKeyStream(buffer, d.Data[MACLen+IVLen:])
	return buffer, nil
}
