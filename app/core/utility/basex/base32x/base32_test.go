package base32x

import (
	"fmt"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	var (
		a   []byte
		b   string
		err error
	)
	a = []byte("a")
	b = Encode(a, true)
	fmt.Println(b)
	a, err = Decode(b)
	fmt.Println(err)
	fmt.Println(string(a))

	a = []byte("aa")
	b = Encode(a, true)
	fmt.Println(b)
	a, err = Decode(b)
	fmt.Println(err)
	fmt.Println(string(a))

	a = []byte("aaa")
	b = Encode(a, true)
	fmt.Println(b)
	a, err = Decode(b)
	fmt.Println(err)
	fmt.Println(string(a))

	a = []byte("aaaa")
	b = Encode(a, true)
	fmt.Println(b)
	a, err = Decode(b)
	fmt.Println(err)
	fmt.Println(string(a))

	a = []byte("aaaaa")
	b = Encode(a, true)
	fmt.Println(b)
	a, err = Decode(b)
	fmt.Println(err)
	fmt.Println(string(a))
}
