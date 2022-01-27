package crypto

import (
	"crypto/hmac"
	cryptoRand "crypto/rand"
	"github.com/pkg/errors"
	"hash"
	"math/rand"
)

type Block struct {
	MAC   []byte
	Nonce []byte
	Data  []byte
}

type IBlock interface {
	FromStream(encrypted Encrypted) error
	ToStream() (Encrypted, error)
}

type Message []byte
type Encrypted []byte

type IEncrypt interface {
	EncryptToStream(sharedKey []byte, optionalNonce []byte) (Encrypted, error)
}

type IDecrypt interface {
	DecryptFromStream(sharedKey []byte) (Message, error)
}

type EncKey []byte
type MacKey []byte

func MakeKeys(hashFunc func() hash.Hash, sharedKey []byte, nonce []byte) (EncKey, MacKey, error) {
	if nil == sharedKey {
		return nil, nil, errors.Errorf("crypto.MakeKeys: sharedKey is nil")
	}
	if nil == nonce {
		return nil, nil, errors.Errorf("crypto.MakeKeys: nonce is nil")
	}
	he := hmac.New(hashFunc, nonce)
	he.Write(sharedKey)
	encKey := he.Sum(nil)
	hm := hmac.New(hashFunc, sharedKey)
	hm.Write(nonce)
	macKey := hm.Sum(nil)
	return encKey, macKey, nil
}

func Rand(n int, trySafe bool) []byte {
	buffer := make([]byte, n)
	if trySafe {
		m, err := cryptoRand.Read(buffer)
		if nil == err && m == n {
			return buffer
		}
	}
	m, err := rand.Read(buffer)
	if nil == err && m == n {
		return buffer
	}
	x := n / 4
	y := n % 4

	a := rand.Uint64()
	for i := 0; i < y; i++ {
		buffer[i] = byte((a >> (8 * i)) & 0xff)
	}
	for i := 0; i < x; i++ {
		j := i*4 + y
		a := rand.Uint64()
		buffer[j] = byte((a >> 24) & 0xff)
		buffer[j+1] = byte((a >> 16) & 0xff)
		buffer[j+2] = byte((a >> 8) & 0xff)
		buffer[j+3] = byte(a & 0xff)
	}
	return buffer
}
