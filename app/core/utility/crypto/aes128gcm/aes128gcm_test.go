package aes256gcm

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestDecrypt(t *testing.T) {
	encryptedWithTag, _ := hex.DecodeString("adcacfc5f318823511b0d918d28ee7f6ba6c76dfae0202958ee598c0c166c3cf762915ae91480ac5bba06467738ac062fd4754332c88c98f0b2e547cf8e6be825a90f16dbfb1f902f23e1003f29f1924b330475ddd5dc227826c07624b6add659ebf14e4ce99f36f15040979e134c55b83d05beec309a578d245a3488615cb9580890b9aa5cb76769dafc80b1a8aa9c233e44f9089e3a50f993cc34799eced36d221ef333fa932e4300cc415935aecec744e3782555520df29b78f38644f2697e7924c3c972b27f9e2f88bb4d39afa923b47ff1943c53fd6a2afb4dd3437d194f9c885ab80c4f437b5babf79e471ae87b663c7cef0b310c7327cd528c1888ad60ec1145d825ef8ee97ff491a417fcd0efeb8f07f23987827096469c7fe64")
	sharedKey := []byte("12345678")

	block := DataAES128GCM{
		SharedKey: sharedKey,
		Data:      encryptedWithTag,
	}

	msg, err := block.Decrypt()
	fmt.Println(err)
	s := string(msg)
	fmt.Println(s)
}

func TestEncrypt(t *testing.T) {
	data := "西城杨柳弄春柔，动离忧，泪难收。犹记多情、曾为系归舟。碧野朱桥当日事，人不见，水空流。韶华不为少年留，恨悠悠，几时休？飞絮落花时候、一登楼。便作春江都是泪，流不尽，许多愁。"
	dataBytes := []byte(data)
	sharedKey := []byte("12345678")
	block := DataAES128GCM{
		SharedKey: sharedKey,
		Data:      dataBytes,
	}
	//encKey, macKey, err := block.MakeKeys()

	fmt.Println(getNonceSize())
	enc, err := block.Encrypt()
	fmt.Println(err)
	fmt.Println(hex.EncodeToString(enc))
}
