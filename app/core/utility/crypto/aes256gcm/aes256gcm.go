package aes256gcm

import (
	"app/core/utility/crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"github.com/pkg/errors"
)

const (
	KeyBitLength    = 256
	KeyByteLength   = KeyBitLength >> 3
	BlockByteLength = 128 >> 3
)

var (
	nonceSize = -1
	hashFunc  = sha256.New
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

type Message []byte
type Encrypted []byte

func (message Message) EncryptToStream(sharedKey []byte, optionalNonce []byte) (Encrypted, error) {
	var (
		nonce     []byte = nil
		nonceSize int
		err       error
	)
	if nil == message {
		return nil, errors.Errorf("crypto.aes256gcm.EncryptToStream: message is nil")
	}
	if nil == sharedKey {
		return nil, errors.Errorf("crypto.aes256gcm.EncryptToStream: sharedKey is nil")

	}
	if nonceSize, err = getNonceSize(); nil != err {
		return nil, err
	}
	if nil != optionalNonce && len(optionalNonce) >= nonceSize {
		nonce = optionalNonce[:nonceSize]
	} else {
		nonce = crypto.Rand(nonceSize, true)
	}

	encKey, macKey, err := crypto.MakeKeys(hashFunc, sharedKey, nonce)
	if nil != err {
		return nil, err
	}
	block, err := aes.NewCipher(encKey)
	if nil != err {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if nil != err {
		return nil, err
	}
	b := BlockAES128GCM{Nonce: nonce}
	b.Data = aesGCM.Seal(nil, nonce, message, macKey)
	encrypted, err := b.ToStream()
	if nil != err {
		return nil, err
	}
	return encrypted, nil
}

func (encrypted Encrypted) DecryptFromStream(sharedKey []byte) (Message, error) {
	if nil == encrypted {
		return nil, errors.Errorf("crypto.aes256gcm.DecryptFromStream: encrypted is nil")
	}
	if nil == sharedKey {
		return nil, errors.Errorf("crypto.aes256gcm.DecryptFromStream: sharedKey is nil")

	}

	b := BlockAES128GCM{}
	if err := b.FromStream(encrypted); nil != err {
		return nil, err
	}
	encKey, macKey, err := crypto.MakeKeys(hashFunc, sharedKey, b.Nonce)
	if nil != err {
		return nil, err
	}

	block, err := aes.NewCipher(encKey)
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.DecryptFromStream: aes.NewCipher")
	}
	aesGCM, err := cipher.NewGCM(block)
	if nil != err {
		return nil, errors.Errorf("crypto.aes256gcm.DecryptFromStream: cipher.NewGCM")
	}
	message, err := aesGCM.Open(nil, b.Nonce, b.Data, macKey)
	if nil != err {
		return nil, err
	}
	return message, nil
}

type BlockAES128GCM crypto.Block

func (b *BlockAES128GCM) FromStream(encrypted Encrypted) error {
	var (
		nonceSize int
		err       error
	)
	if nonceSize, err = getNonceSize(); nil != err {
		return err
	}
	if len(encrypted) <= nonceSize {
		return errors.Errorf("crypto.aes256gcm.FromStream: encrypted is too short")
	}
	b.Nonce = encrypted[:nonceSize]
	b.Data = encrypted[nonceSize:]
	//if len(b.Data)%BlockByteLength != 0 {
	//	return errors.Errorf("crypto.aes256gcm.FromStream: encrypted length wrong")
	//}
	return nil
}

func (b *BlockAES128GCM) ToStream() (Encrypted, error) {
	if nil == b.Nonce {
		return nil, errors.Errorf("crypto.aes256gcm.ToStream: Nonce is nil")
	}
	if nil == b.Data {
		return nil, errors.Errorf("crypto.aes256gcm.ToStream: Data is nil")
	}
	buffer := make([]byte, len(b.Nonce)+len(b.Data))
	copy(buffer, b.Nonce)
	copy(buffer[len(b.Nonce):], b.Data)
	return buffer, nil
}
