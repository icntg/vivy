package qrcode

import "runtime"

const (
	WindowsWhiteBlock = "▇"
	WindowsBlackBlock = "  "
	WindowsNewLine    = "\n"
	XNixWhiteBlock    = "\033[0;37;47m  "
	XNixBlackBlock    = "\033[0;37;40m  "
	XNixNewLine       = "\033[0m\n"
)

var (
	whiteBlock = ""
	blackBlock = ""
	newLine    = ""
)

func init() {
	if runtime.GOOS == "windows" {
		whiteBlock = WindowsWhiteBlock
		blackBlock = WindowsBlackBlock
		newLine = WindowsNewLine
	} else {
		whiteBlock = XNixWhiteBlock
		blackBlock = XNixBlackBlock
		newLine = XNixNewLine
	}
}
