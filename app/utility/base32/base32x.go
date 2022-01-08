package base32

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

const (
	Base32xPadding     = "_ILOV"
	Base32xEncodeTable = "0123456789ABCDEFGHJKMNPQRSTUWXYZ"
	Base32xDecodeTable = "" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xff\xff\xff\xff\xff\xff" +
		"\xff\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\xff\x12\x13\xff\x14\x15\xff" +
		"\x16\x17\x18\x19\x1a\x1b\xff\x1c\x1d\x1e\x1f\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"
)

func Encode(in []byte, withPadding bool) string {
	var (
		buffer     []byte
		inBitLen   int
		inByteLen  int
		encByteLen int
		inBytes    []byte
		idx        byte
		paddingIdx int = 0
	)
	inByteLen = len(in)
	inBitLen = inByteLen << 3
	if inBitLen%5 == 0 {
		encByteLen = inBitLen / 5
		inBytes = in
	} else {
		encByteLen = inBitLen/5 + 1
		inBytes = make([]byte, inByteLen+1)
		copy(inBytes, in)
		inBytes[inByteLen] = 0
		if withPadding {
			paddingIdx = (inBitLen/5+1)*5 - inBitLen
		}
	}
	if withPadding && paddingIdx > 0 {
		buffer = make([]byte, encByteLen+1)
	} else {
		buffer = make([]byte, encByteLen)
	}

	for i := 0; i < encByteLen; i++ {
		k := i * 5
		j := k + 4
		idx0 := k / 8
		idx1 := j / 8
		bit0 := k % 8
		bit1 := j % 8
		if idx0 == idx1 {
			idx = (inBytes[idx0] >> (7 - bit1)) & 0x1f
		} else {
			idx = ((inBytes[idx0] << (5 - (8 - bit0))) | ((inBytes[idx1] >> (7 - bit1)) & 0x1f)) & 0x1f
		}
		buffer[i] = Base32xEncodeTable[idx]
	}
	if withPadding && paddingIdx > 0 {
		buffer[encByteLen] = Base32xPadding[paddingIdx]
	}
	return string(buffer)
}

func EncodeId(in []byte) string {
	var (
		buffer     []byte
		inBitLen   int
		inByteLen  int
		encBitLen  int
		encByteLen int
		inBytes    []byte
		idx        byte
	)
	inByteLen = len(in)
	inBitLen = inByteLen << 3
	if inBitLen%5 == 0 {
		encBitLen = inBitLen
		encByteLen = inBitLen / 5
		inBytes = in
	} else {
		encByteLen = inBitLen/5 + 1
		encBitLen = inBitLen + 8
		inBytes = []byte{0}
		inBytes = append(inBytes[:], in...)
	}
	buffer = make([]byte, encByteLen)
	for i := 0; i < encByteLen; i++ {
		k := encBitLen - i*5 - 1
		j := k - 4
		idx0 := j / 8
		idx1 := k / 8
		bit0 := j % 8
		bit1 := k % 8
		if idx0 == idx1 {
			idx = (inBytes[idx0] >> (7 - bit1)) & 0x1f
		} else {
			idx = ((inBytes[idx0] << (5 - (8 - bit0))) | ((inBytes[idx1] >> (7 - bit1)) & 0x1f)) & 0x1f
		}
		buffer[encByteLen-i-1] = Base32xEncodeTable[idx]
	}
	return string(buffer)
}

func Decode(in string) ([]byte, error) {
	var (
		padding    int
		decBitLen  int
		encByteLen int
		encBuf     []byte
		decBuf     []byte
	)
	c := strings.ToUpper(in[len(in)-1 : len(in)])
	padding = strings.Index(Base32xPadding, c)
	if padding <= 0 {
		padding = 0
		decBitLen = len(in) * 5
		encByteLen = len(in)
	} else {
		decBitLen = (len(in)-1)*5 - padding
		encByteLen = len(in) - 1
	}
	if decBitLen%8 != 0 {
		return nil, errors.Wrap(nil, "base32x decode length")
	}

	encBuf = make([]byte, encByteLen)
	copy(encBuf, in[0:encByteLen])

	decByteLen := decBitLen / 8
	decBuf = make([]byte, decByteLen+1)
	for i := 0; i < encByteLen; i++ {
		e := encBuf[i]
		d := Base32xDecodeTable[e]
		if d > 31 {
			return nil, errors.Wrap(nil, fmt.Sprintf("decode error at %d: %v", i, e))
		}
		bitIdx := i * 5
		idx0 := bitIdx / 8
		idx1 := (bitIdx + 4) / 8
		bit0 := bitIdx % 8
		bit1 := (bitIdx + 4) % 8
		if idx0 == idx1 {
			decBuf[idx0] |= d << (7 - bit1)
		} else {
			decBuf[idx0] |= d >> (5 - (8 - bit0))
			decBuf[idx1] |= byte((d << (7 - bit1)) & 0xff)
		}
	}
	return decBuf[0:decByteLen], nil
}
