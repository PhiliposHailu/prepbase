package infrastructure

import (
	"sync"
	"time"

	"github.com/philipos/prepbase/domain"
)

type cacheItem struct {
	value      interface{}
	expiration int64
}

type memoryCache struct {
	items map[string]cacheItem
	mu    sync.RWMutex // Protects the map from race conditions
}

func NewMemoryCache() domain.CacheService {
	return &memoryCache{
		items: make(map[string]cacheItem),
	}
}

func (c *memoryCache) Set(key string, value interface{}) {
	c.mu.Lock()         // LOCK: Nobody else can read or write right now
	defer c.mu.Unlock() // UNLOCK when done

	// Store the data and set it to expire in 1 minute
	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(1 * time.Minute).Unix(),
	}
}

func (c *memoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()         // READ LOCK: Multiple people can read, but nobody can write
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if it expired
	if time.Now().Unix() > item.expiration {
		// It will naturally be overwritten next time Set() is called.
		return nil, false 
	}

	return item.value, true
}