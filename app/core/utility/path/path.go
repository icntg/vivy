package path

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetBinaryPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if nil != err {
		return "", err
	}
	bp, err := filepath.Abs(file)
	if nil != err {
		return "", err
	}
	idx := strings.LastIndex(bp, string(os.PathSeparator))
	ret := bp[:idx]
	return strings.Replace(ret, "\\", "/", -1), nil
}
