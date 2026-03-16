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
	cache := &memoryCache{
		items: make(map[string]cacheItem),
	}
	go cache.startGarbageCollector(5 * time.Minute)
	return cache
}

func (c *memoryCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mu.Lock()         // LOCK: Nobody else can read or write right now
	defer c.mu.Unlock() // UNLOCK when done

	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration).Unix(), // Dynamic expiration time
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
		return nil, false 
	}

	return item.value, true
}

// THE GARBAGE COLLECTOR
func (c *memoryCache) startGarbageCollector(interval time.Duration) {
	// A Ticker fires an event every X minutes
	ticker := time.NewTicker(interval)
	
	// Infinite loop running in the background
	for {
		<-ticker.C // Wait for the ticker to fire
		
		now := time.Now().Unix()

		c.mu.Lock() // We need a full Write Lock to delete items
		for key, item := range c.items {
			if now > item.expiration {
				delete(c.items, key) // Remove from RAM!
			}
		}
		c.mu.Unlock()
	}
}