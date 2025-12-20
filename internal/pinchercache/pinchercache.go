// Package pinchercache provides a cache for the pincher-cli
package pinchercache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time // time entry was created
	data      []byte    // raw data we're caching
}

type Cache struct {
	cachedEntries map[string]cacheEntry
	interval      time.Duration
	mu            *sync.Mutex
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cachedEntries[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		data:      value,
	}
}

func (c *Cache) Get(key string) (entryData []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.cachedEntries[key]
	if !ok {
		return nil, false
	}

	return val.data, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap()
	}
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, val := range c.cachedEntries {
		if time.Since(val.createdAt) > c.interval {
			delete(c.cachedEntries, key)
		}
	}
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cachedEntries: make(map[string]cacheEntry),
		interval:      interval,
		mu:            &sync.Mutex{},
	}

	go cache.reapLoop()

	return cache
}
