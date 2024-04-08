package repository

import (
	"github.com/gookit/slog"
	"sync"
	"time"
)

type LocalCache struct {
	cache map[string]item
	mutex sync.RWMutex
}

type item struct {
	value    []byte
	expireAt time.Time
}

func NewLocalCache() *LocalCache {
	slog.Info("Setting up new Local Cache . . .")
	return &LocalCache{
		cache: make(map[string]item),
	}
}

func (c *LocalCache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.cache[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expireAt) {
		c.mutex.RUnlock()
		c.mutex.Lock()
		delete(c.cache, key)
		c.mutex.Unlock()
		c.mutex.RLock()
		return nil, false
	}
	return item.value, true
}

func (c *LocalCache) Put(key string, value []byte, ttl time.Duration) {
	c.mutex.Lock()

	c.cache[key] = item{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}

	c.mutex.Unlock()
}
