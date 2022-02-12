package qrcode

import (
	QRCode "github.com/skip2/go-qrcode"
	"strings"
)

const (
	WindowsWhiteBlock = "▇"
	WindowsBlackBlock = " "
	WindowsNewLine    = "\r\n"
	XNixWhiteBlock    = "\033[0;37;40m  "
	XNixBlackBlock    = "\033[0;37;47m  "
	XNixNewLine       = "\033[0m\n"
)

type ConsoleStyle struct {
	WhiteBlock string
	BlackBlock string
	NewLine    string
}

var (
	Styles = [...]ConsoleStyle{
		{XNixWhiteBlock, XNixBlackBlock, XNixNewLine},
		{XNixBlackBlock, XNixWhiteBlock, XNixNewLine},
		{WindowsWhiteBlock, WindowsBlackBlock, WindowsNewLine},
		{WindowsWhiteBlock + WindowsWhiteBlock, WindowsBlackBlock + WindowsBlackBlock, WindowsNewLine},
	}
)

func ConsoleQRCode(in string) []string {
	q, _ := QRCode.New(in, QRCode.Medium)
	bitmap := q.Bitmap()
	ret := make([]string, len(Styles))
	for i, s := range Styles {
		buffer := strings.Builder{}
		for i := 0; i < len(bitmap); i++ {
			for j := 0; j < len(bitmap[i]); j++ {
				if bitmap[i][j] {
					buffer.WriteString(s.WhiteBlock)
				} else {
					buffer.WriteString(s.BlackBlock)
				}
			}
			buffer.WriteString(s.NewLine)
		}
		ret[i] = buffer.String()
	}
	return ret
}
