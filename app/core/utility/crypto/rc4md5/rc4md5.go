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

var (
	hashFunc = md5.New
)

type Message []byte
type Encrypted []byte
type EncKey []byte
type MacKey []byte

// [MAC][IV][Encrypted]
// encKey = hmac_md5(iv, key)
// macKey = hmac_md5(key, iv)
// Encrypted = rc4(encKey, message)
// MAC = hmac_md5(macKey, iv + message)

func (message Message) EncryptToStream(sharedKey []byte, optionalNonce []byte) (Encrypted, error) {
	var (
		nonce []byte = nil
	)
	if nil == sharedKey {
		return nil, errors.Errorf("crypto.rc4md5.EncryptToStream: sharedKey is nil")
	}
	if nil == message {
		return nil, errors.Errorf("crypto.rc4md5.EncryptToStream: message is nil")
	}
	if nil == optionalNonce || len(optionalNonce) < IVLen {
		nonce = crypto.Rand(IVLen, true)
	} else {
		nonce = optionalNonce[:IVLen]
	}
	encKey, macKey, err := crypto.MakeKeys(hashFunc, sharedKey, nonce)
	if nil != err {
		return nil, err
	}
	cipher, err := rc4.NewCipher(encKey)
	if nil != err {
		return nil, err
	}
	b := BlockRC4MD5{
		MAC:   nil,
		Nonce: nonce,
		Data:  make(Message, len(message)),
	}

	cipher.XORKeyStream(b.Data, message)
	h := hmac.New(hashFunc, macKey)
	h.Write(b.Nonce)
	h.Write(b.Data)
	b.MAC = h.Sum(nil)
	encrypted, err := b.ToStream()
	if nil != err {
		return nil, err
	}
	return encrypted, nil
}

func (encrypted Encrypted) DecryptFromStream(sharedKey []byte) (Message, error) {
	if nil == encrypted {
		return nil, errors.Errorf("crypto.rc4md5.DecryptFromStream: encrypted is nil")
	}
	if nil == sharedKey {
		return nil, errors.Errorf("crypto.rc4md5.DecryptFromStream: sharedKey is nil")
	}
	b := BlockRC4MD5{}
	if err := b.FromStream(encrypted); nil != err {
		return nil, err
	}
	encKey, macKey, err := crypto.MakeKeys(hashFunc, sharedKey, b.Nonce)
	if nil != err {
		return nil, err
	}
	h := hmac.New(hashFunc, macKey)
	h.Write(b.Nonce)
	h.Write(b.Data)
	realMAC := h.Sum(nil)
	if !bytes.Equal(realMAC, b.MAC) {
		return nil, errors.Errorf("crypto.rc4md5.DecryptFromStream: MAC verify failed")
	}
	cipher, err := rc4.NewCipher(encKey)
	if nil != err {
		return nil, err
	}
	buffer := make([]byte, len(b.Data))
	cipher.XORKeyStream(buffer, b.Data)
	return buffer, nil
}

type BlockRC4MD5 crypto.Block

func (d *BlockRC4MD5) FromStream(encrypted Encrypted) error {
	var ()
	if len(encrypted) <= MACLen+IVLen {
		return errors.Errorf("crypto.rc4md5.FromStream: encrypted is too short")
	}
	d.MAC = encrypted[:MACLen]
	d.Nonce = encrypted[MACLen:][:IVLen]
	d.Data = encrypted[MACLen:][:IVLen]
	return nil
}

func (d *BlockRC4MD5) ToStream() (Encrypted, error) {
	if nil == d.Nonce {
		return nil, errors.Errorf("crypto.rc4md5.ToStream: Nonce is nil")
	}
	if nil == d.Data {
		return nil, errors.Errorf("crypto.rc4md5.ToStream: Data is nil")
	}
	buffer := make([]byte, MACLen+IVLen+len(d.Data))
	copy(buffer, d.MAC)
	copy(buffer[MACLen:], d.Nonce)
	copy(buffer[MACLen+IVLen:], d.Data)
	return buffer, nil
}
