package aes256gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

func TestDecrypt(t *testing.T) {
	nonce, _ := hex.DecodeString("3d0bf2e2ef5e3f1b394049ee")
	encryptedWithTag, _ := hex.DecodeString("d8846b611e3daa4726de546393448e384da1fc84a290e7ea1f9fdad84236bd93a6c9479c543a575982a1c1377d46024e3d6672fef318f252508da182fccc3138dfaa1619093c0f844d091d32b28e1b71ea5bf9aa26e44152113fc53daf2181a83ee45c3de45d272bce6dd1e017f37a1d35a37e286ab810d2036516b8b7577365a2821d1885dc2d12d022c8413f480356be61f3b2ab5d10c1c75bb7393bc5d7efc1ef4894b75dee0b8f4deb697509a9996c16eff2f156ba6cf98bd0b07273d5e66c5b037b74dd5a23d097b7d94164dc2cd703b02477749efc6db8ed5c76601138602c7e099b8e18f158dc90af1b6c8d5fcbea4949b6cab6483732b84256d61f4f3ad1d03155ff92a7260a2dcd207f413cd192")
	sharedKey := []byte("12345678")

	encKey, macKey, _ := MakeKeys(sharedKey, nonce)
	_ = macKey
	block, _ := aes.NewCipher(encKey)
	aesGCM, _ := cipher.NewGCM(block)
	msg, err := aesGCM.Open(nil, nonce, encryptedWithTag, encKey)
	fmt.Println(err)
	s := string(msg)
	fmt.Println(s)
}

func TestEncrypt(t *testing.T) {
	data := "西城杨柳弄春柔，动离忧，泪难收。犹记多情、曾为系归舟。碧野朱桥当日事，人不见，水空流。韶华不为少年留，恨悠悠，几时休？飞絮落花时候、一登楼。便作春江都是泪，流不尽，许多愁。"
	dataBytes := []byte(data)
	sharedKey := []byte("12345678")
	enc, err := Encrypt(sharedKey, nil, dataBytes)
	fmt.Println(err)
	fmt.Println(hex.EncodeToString(enc))
}

func TestMakeKeys(t *testing.T) {
	type args struct {
		sharedKey []byte
		iv        []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		want1   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := MakeKeys(tt.args.sharedKey, tt.args.iv)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeKeys() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("MakeKeys() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
