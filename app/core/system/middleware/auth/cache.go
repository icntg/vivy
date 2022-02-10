package auth

import "sync"

var (
	_rwLock sync.RWMutex
)

type CacheInterface interface {
	SyncWithDatabase() // 从数据库中刷入Cache
	IsAccessible(userDbId uint32, method, url string) bool
}

type Cache map[string]map[string]string

func (c *Cache) SyncWithDatabase() {
	_rwLock.Lock()
	defer _rwLock.Unlock()

	//TODO implement me
	panic("implement me")
}

func (c *Cache) IsAccessible(userDbId uint32, method, url string) bool {
	_rwLock.RLock()
	defer _rwLock.RUnlock()

	//TODO implement me
	panic("implement me")
}
