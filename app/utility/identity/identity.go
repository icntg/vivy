package identity

import (
	"app/utility/logger"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

var (
	counterMutex sync.Mutex
	counterValue uint32
	macValue     [3]byte
	pidValue     [2]byte
)

func init() {
	counterValue = uint32(rand.Int()) & 0xffffff // 此处使用伪随机数即可
	macValue = mac()
	pidValue = pid()
}

func mac() [3]byte {
	interfaces, err := net.Interfaces()
	if nil != err {
		log, _ := logger.GetOutputLogger()
		log.Warnf("Cannot get MAC. Use random value instead. %v\n", err)
	}

	for _, inter := range interfaces {
		inter.HardwareAddr.String()
		fmt.Println(inter.Name, inter.HardwareAddr)
		break
	}
	return [3]byte{0, 0, 0}
}

// 如果检测到docker，则使用docker的id。否则才使用pid。
func pid() [2]byte {
	return [2]byte{0, 0}
}

func timestamp() [4]byte {
	now := uint32(time.Now().Unix())
	result := [4]byte{0, 0, 0, 0}
	result[3] = byte(now & 0xff)
	result[2] = byte((now >> 8) & 0xff)
	result[1] = byte((now >> 16) & 0xff)
	result[0] = byte((now >> 24) & 0xff)
	return result
}

func counter() [3]byte {
	var c uint32
	counterMutex.Lock()
	counterValue++
	counterValue = counterValue & 0xffffff
	c = counterValue
	counterMutex.Unlock()
	result := [3]byte{0, 0, 0}
	result[2] = byte(c & 0xff)
	result[1] = byte((c >> 8) & 0xff)
	result[0] = byte((c >> 16) & 0xff)
	return result
}

func ObjectId() [12]byte {
	return [12]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

func ObjectIdB64() string {
	objectId := ObjectId()
	return base64.StdEncoding.EncodeToString(objectId[:])
}

func ObjectIdHex() string {
	objectId := ObjectId()
	return hex.EncodeToString(objectId[:])
}
