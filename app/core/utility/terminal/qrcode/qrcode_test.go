package qrcode

import (
	"fmt"
	"testing"
)

func TestConsoleQRCode(t *testing.T) {
	a := "otpauth://totp/google:admin?algorithm=SHA1&digits=6&period=30&issuer=google&secret=github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e"
	x := ConsoleQRCode(a)
	fmt.Print(x)
}
