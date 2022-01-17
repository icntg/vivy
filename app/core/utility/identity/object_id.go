package identity

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

var (
	counterMutex   sync.Mutex
	counterValue   uint32
	macValue       [3]byte
	pidValue       [2]byte
	pidDockerCache bool
)

func init() {
	counterValue = uint32(rand.Int()) & 0xffffff // 此处使用伪随机数即可
	macValue = getMAC()
	pidDockerCache = false
	pidValue = getPid()
}

func getMAC() [3]byte {
	netInterfaces, err := net.Interfaces()
	if nil == err {
		for _, netInterface := range netInterfaces {
			macAddrStr := netInterface.HardwareAddr.String()
			if len(macAddrStr) > 0 {
				m := netInterface.HardwareAddr
				mac := [...]byte{m[0] ^ m[3], m[1] ^ m[4], m[2] ^ m[5]}
				return mac
			}
		}
	}
	// 如果无法取得网卡信息，或者没有mac可用，使用MD5(hostname)。
	log.Printf("Cannot get MAC information. To use HOSTNAME instead.\n")
	hostname, err := os.Hostname()
	if nil == err {
		h := md5.New()
		h.Write([]byte(hostname))
		m := h.Sum(nil)
		mac := [...]byte{m[0], m[1], m[2]}
		return mac
	}
	// 如果主机名也无法取得，使用随机数
	log.Printf("Cannot get HOSTNAME either. To use RANDOM instead.\n")
	r := rand.Int()
	return [...]byte{byte((r >> 16) & 0xff), byte((r >> 8) & 0xff), byte(r & 0xff)}
}

// 如果检测到docker，则使用docker的id。否则才使用pid。
func getPid() [2]byte {
	pid := os.Getpid()
	if pid == 1 {
		if pidDockerCache {
			return pidValue
		}
		// 某些情况下，docker容器中pid会为1，这样就失去的这个字段的作用。
		// 尝试读取docker id替代。
		const dockerFile = "/proc/self/cgroup"
		const sep = ":/docker/"
		_, err := os.Stat(dockerFile)
		r := rand.Uint64()
		if nil != err {
			return [2]byte{byte(r & 0xff), byte((r >> 8) & 0xff)}
		}
		fc, err := ioutil.ReadFile(dockerFile)
		if nil != err {
			return [2]byte{byte((r >> 16) & 0xff), byte((r >> 24) & 0xff)}
		}
		idx := bytes.Index(fc, []byte(sep))
		if idx < 0 {
			return [2]byte{byte((r >> 32) & 0xff), byte((r >> 40) & 0xff)}
		}
		idx += len(sep)
		buf := string(fc[idx : idx+4])
		result, err := hex.DecodeString(buf)
		if nil != err {
			return [2]byte{byte((r >> 48) & 0xff), byte((r >> 56) & 0xff)}
		}
		pidDockerCache = true
		pidValue = [...]byte{result[0], result[1]}
		return pidValue
	}
	return [...]byte{byte((pid >> 8) & 0xff), byte(pid & 0xff)}
}

func getTimestamp() [4]byte {
	now := uint32(time.Now().Unix())
	result := [4]byte{0, 0, 0, 0}
	result[3] = byte(now & 0xff)
	result[2] = byte((now >> 8) & 0xff)
	result[1] = byte((now >> 16) & 0xff)
	result[0] = byte((now >> 24) & 0xff)
	return result
}

func getCounter() [3]byte {
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
	ts := getTimestamp()
	c := getCounter()
	pid := getPid()
	return [12]byte{ts[0], ts[1], ts[2], ts[3],
		macValue[0], macValue[1], macValue[2],
		pid[0], pid[1],
		c[0], c[1], c[2]}
}

func ObjectIdB64() string {
	objectId := ObjectId()
	return base64.StdEncoding.EncodeToString(objectId[:])
}

func ObjectIdHex() string {
	objectId := ObjectId()
	return hex.EncodeToString(objectId[:])
}
