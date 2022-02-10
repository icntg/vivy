package base36

import (
	"fmt"
	"testing"
)

func TestDecodeString(t *testing.T) {
	a := "1dhpndum9ltrtpk8hnrqp5o3is"
	x, err := DecodeString(a)
	fmt.Println(err)
	fmt.Println(x)
	y := EncodeToStringLc(x)
	fmt.Println(a)
	fmt.Println(y)
}
