package crypto

import (
	cryptoRand "crypto/rand"
	"math/rand"
)

type DataCrypto struct {
	SharedKey []byte
	Data      []byte
	Nonce     []byte
}

type InterfaceCrypto interface {
	MakeKeys() ([]byte, []byte, error)
	Encrypt() ([]byte, error)
	Decrypt() ([]byte, error)
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
