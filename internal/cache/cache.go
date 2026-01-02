// Package cache provides a cache for the pincher-cli
package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	CreatedAt time.Time `json:"created_at"`
	Data      []byte    `json:"data"`
}

type Cache struct {
	CachedEntries map[string]cacheEntry `json:"cached_entries"`
	interval      time.Duration
	mu            *sync.Mutex
}

func (c *Cache) Set(entries map[string]cacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CachedEntries = entries
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CachedEntries[key] = cacheEntry{
		CreatedAt: time.Now().UTC(),
		Data:      value,
	}
}

func (c *Cache) Get(key string) (entryData []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.CachedEntries[key]
	if !ok {
		return nil, false
	}

	return val.Data, true
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

	for key, val := range c.CachedEntries {
		if time.Since(val.CreatedAt) > c.interval {
			delete(c.CachedEntries, key)
		}
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.CachedEntries, key)
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		CachedEntries: make(map[string]cacheEntry),
		interval:      interval,
		mu:            &sync.Mutex{},
	}

	go cache.reapLoop()

	return cache
}
