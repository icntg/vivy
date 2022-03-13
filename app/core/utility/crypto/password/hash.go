package password

import (
	"app/core/utility/crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"github.com/pkg/errors"
	"hash"
	"strings"
)

const b64table = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	algorithms map[string]func() hash.Hash = nil
)

func init() {
	algorithms = make(map[string]func() hash.Hash)
	algorithms["1"] = md5.New
	algorithms["5"] = sha256.New
	algorithms["6"] = sha512.New
}

type HashOptions struct {
	Nonce []byte
	Algo  string
}

type HashResult struct {
	Algo                 string
	Nonce                []byte
	NonceString          string
	HashedPassword       []byte
	HashedPasswordString string
}

func (r HashResult) Error() string {
	return "PasswordHashResultError"
}

func b64encode(in []byte) string {
	return strings.ReplaceAll(base64.NewEncoding(b64table).EncodeToString(in), "=", "")
}

func b64decode(in string) ([]byte, error) {
	var (
		s string
	)
	switch len(in) % 4 {
	case 0:
		s = in
		break
	case 2:
		s = in + "=="
		break
	case 3:
		s = in + "="
		break
	default:
		return nil, errors.New("base64decode padding size")
	}
	return base64.NewEncoding(b64table).DecodeString(s)
}

func Hash(pwd string, options *HashOptions) string {
	var (
		opt *HashOptions
	)
	if nil != options {
		opt = options
	} else {
		opt = &HashOptions{nil, "5"}
	}
	if nil == opt.Nonce {
		opt.Nonce = crypto.Rand(12, true)
	}
	if _, ok := algorithms[opt.Algo]; !ok {
		opt.Algo = "5"
	}
	m := hmac.New(algorithms[opt.Algo], opt.Nonce)
	m.Write([]byte(pwd))
	enc := m.Sum(nil)

	buffer := strings.Builder{}
	buffer.WriteString("$")
	buffer.WriteString(opt.Algo)
	buffer.WriteString("$")
	buffer.WriteString(b64encode(opt.Nonce))
	buffer.WriteString("$")
	buffer.WriteString(b64encode(enc))

	return buffer.String()
}

///*
//模拟密码校验第一阶段
//数据库中保存的hashed校验失败，或者无法提取nonce
//*/
//func fakeVerify0(pwd string) (bool, error) {
//}
//
///*
//*/
//func fakeVerify1(pwd string, )

func Verify(pwd string, hashed string) (bool, error) {
	var (
		result          = false
		err      error  = nil
		nonce    []byte = nil
		inHashed        = ""
	)
	arr := strings.Split(strings.TrimSpace(hashed), "$")
	if 4 != len(arr) {
		err = errors.Wrap(HashResult{}, "hashed format size")
		goto FAKE
	}
	if 0 != len(arr[0]) {
		err = errors.Wrap(HashResult{}, "hashed array 0")
		goto FAKE
	}
	if _, ok := algorithms[arr[1]]; !ok {
		err = errors.Wrap(HashResult{}, "hashed algorithm not exist")
		goto FAKE
	}
	nonce, err = b64decode(arr[2])
	if nil != err {
		err = errors.Wrap(err, "hashed nonce decode")
		goto FAKE
	}

	{
		m := hmac.New(algorithms[arr[1]], nonce)
		m.Write([]byte(pwd))
		enc := m.Sum(nil)
		inHashed = b64encode(enc)
	}

	if len(arr[3]) != len(inHashed) {
		err = errors.Wrap(err, "hashed length not equal")
		goto FAKE
	}
	result = true
	for i := 0; i < len(inHashed); i++ {
		result = result && arr[3][i] == inHashed[i]
	}
	return result, nil
FAKE:
	// TODO: mock the delay time of password verify
	return false, err
}
