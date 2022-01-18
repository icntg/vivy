package crypto

import (
	cryptoRand "crypto/rand"
	"github.com/pkg/errors"
	"math/rand"
)

const (
	SumLen = 16
	IVLen  = 8
)

type Block struct {
	Sum  [SumLen]byte
	IV   [IVLen]byte
	Data []byte
}

func StreamToBlock(stream []byte) (Block, error) {
	block := Block{}
	if nil == stream {
		return block, errors.Errorf("crypto.@crypto.StreamToBlock: stream is nil")
	}
	if len(stream) <= SumLen+IVLen+1 {
		return block, errors.Errorf("crypto.@crypto.StreamToBlock: stream is too short")
	}
	copy(block.Sum[:], stream[:SumLen])
	copy(block.IV[:], stream[SumLen:SumLen+IVLen])
	block.Data = make([]byte, len(stream)-SumLen-IVLen)
	copy(block.Data, stream[SumLen+IVLen:])
	return block, nil
}

func (ths *Block) ToStream() ([]byte, error) {
	if nil == ths {
		return nil, errors.Errorf("crypto.@crypto.ToStream: this is nil")
	}
	n := SumLen + IVLen + len(ths.Data)
	buffer := make([]byte, n)
	copy(buffer, ths.Sum[:])
	copy(buffer[SumLen:], ths.IV[:])
	copy(buffer[SumLen+IVLen:], ths.Data)
	return buffer, nil
}

func (ths *Block) FromStream(stream []byte) error {
	if nil == stream {
		return errors.Errorf("crypto.@crypto.FromStream: stream is nil")
	}
	if len(stream) < SumLen+IVLen+1 {
		return errors.Errorf("crypto.@crypto.FromStream: stream is too short")
	}
	copy(ths.Sum[:], stream[:SumLen])
	copy(ths.IV[:], stream[SumLen:SumLen+IVLen])
	ths.Data = make([]byte, len(stream)-SumLen-IVLen)
	copy(ths.Data, stream[SumLen+IVLen:])
	return nil
}

type Crypto interface {
	MakeKeys(sharedKey, iv []byte) ([]byte, []byte, error)
	Encrypt(sharedKey []byte, iv []byte, message []byte) ([]byte, error)
	Decrypt(sharedKey, encrypted []byte) ([]byte, error)
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
