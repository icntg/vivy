package crypto

import (
	"fmt"
	"testing"
)

func TestComparePassword(t *testing.T) {
	hashedPassword := "c4nu5d0d1rqza8rpq8ky6kbx6c2p846t"
	salt := "uyeprhlu47zxmrzsbyh6sqnrvhhy4glo"
	var x bool
	x = ComparePassword("adfasdfasdf", hashedPassword, salt)
	fmt.Println(x)
	x = ComparePassword("9271354313d094536328", hashedPassword, salt)
	fmt.Println(x)
}
