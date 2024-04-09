package repository

import (
	"adt/model"
	"github.com/gookit/slog"
	"sync"
	"time"
)

type LocalCache struct {
	cache map[string]item
	mutex sync.RWMutex
}

type item struct {
	value    []model.Campaign
	expireAt time.Time
}

func NewLocalCache() *LocalCache {
	slog.Info("Setting up new Local Cache . . .")
	return &LocalCache{
		cache: make(map[string]item),
	}
}

func (c *LocalCache) Get(key string) ([]model.Campaign, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.cache[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expireAt) {
		return nil, false
	}
	return item.value, true
}

func (c *LocalCache) Put(key string, value []model.Campaign, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[key] = item{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}

}
