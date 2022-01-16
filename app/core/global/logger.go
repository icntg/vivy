package global

import (
	"app/core/utility/logger"
	"log"
	"os"
	"sync"
)

var (
	loggersInstance *logger.Loggers = nil
	loggersOnce     sync.Once
)

func LoggersInstance() *logger.Loggers {
	loggersOnce.Do(func() {
		config := ConfigInstance()
		if nil != config {
			log.SetOutput(os.Stderr)
			log.Fatalf("TODO:")
		}
	})
	return loggersInstance
}
