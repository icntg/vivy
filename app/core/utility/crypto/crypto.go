package crypto

import (
	"app/core/utility/basex/base32x"
	"crypto/hmac"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"github.com/pkg/errors"
	"hash"
	"math/rand"
	"strings"
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

const PasswordHashLength = 32

func EncPassword(clearPassword string, b32lowerSalt string) (string, error) {
	// 使用标准的base32，便于其他程序兼容
	salt, err := base32.StdEncoding.DecodeString(strings.ToUpper(b32lowerSalt))
	if nil != err {
		return "", errors.Wrap(err, "salt decode failed")
	}
	pb := []byte(clearPassword)
	h := hmac.New(sha256.New, salt)
	h.Write(pb)
	ep := h.Sum(nil)
	ret := strings.ToLower(base32x.Encode(ep, true)[:PasswordHashLength])
	return ret, nil
}

func EncPasswordInit(clearPassword string) (string, string) {
	salt := Rand(20, true)
	pb := []byte(clearPassword)
	h := hmac.New(sha256.New, salt)
	h.Write(pb)
	ep := h.Sum(nil)
	ret := strings.ToLower(base32x.Encode(ep, true)[:PasswordHashLength])
	return ret, strings.ToLower(base32.StdEncoding.EncodeToString(salt))
}

func ComparePassword(clearPassword, storedPassword, b32lowerSalt string) bool {
	expect, err := EncPassword(clearPassword, b32lowerSalt)
	if nil != err {
		// 随机生成。此处特殊符号。
		const x = "`~!@#$%^&*()_-=+\\|[]{},./<>?"
		expect = strings.ToUpper(hex.EncodeToString(Rand(PasswordHashLength/2, true)))
		expect = expect[:len(x)] + x
	}
	fullFill := func(in string) [PasswordHashLength]byte {
		buffer := [PasswordHashLength]byte{}
		bin := []byte(in)
		if len(bin) < PasswordHashLength {
			n := PasswordHashLength - len(bin)
			rb := Rand(n, true)
			copy(buffer[:], bin)
			copy(buffer[len(bin):], rb)
		} else if len(bin) > PasswordHashLength {
			copy(buffer[:], bin[:PasswordHashLength])
		} else {
			copy(buffer[:], bin)
		}
		return buffer
	}
	a := fullFill(expect)
	b := fullFill(storedPassword)
	ret := true
	for i := 0; i < PasswordHashLength; i++ {
		ret = ret && a[i] == b[i]
	}
	return ret
}
