package redis

import (
	"app/core/global/config"
	"app/core/utility/common"
	"app/core/utility/errno"
	"github.com/gin-contrib/sessions/redis"
	"os"
	"sync"
)

var (
	_redisStore *redis.Store = nil
	_once       sync.Once
)

func initStores() {

}

func SessionStoreInstance() *redis.Store {
	_once.Do(func() {
		var (
			gCfg = config.Instance()
		)
		rc := gCfg.DataSource.Redis
		_store, err := redis.NewStore(
			rc.MaxIdle,
			rc.Protocol,
			rc.Address,
			rc.Password,
			gCfg.Service.SessionSecretBytes,
		)
		if nil != err {
			common.ErrPrintf("redis cannot make new store: %v\n", err)
			os.Exit(errno.ErrorRedisStore.Code())
		}
		_redisStore = &_store
	})
	return _redisStore
}
