package client

import (
	"maps"
	"strings"
	"sync"
	"time"
)

type cacheEntry struct {
	CreatedAt time.Time `json:"created_at"`
	Data      []byte    `json:"data"`
	// entry is protected from reaping
	Protected bool `json:"protected"`
}

type Cache struct {
	CachedEntries map[string]cacheEntry `json:"cached_entries"`
	interval      time.Duration
	mu            *sync.Mutex
}

func (c *Cache) Set(entries map[string]cacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	maps.Copy(c.CachedEntries, entries)
}

func (c *Cache) Add(key string, value []byte, protect bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.CachedEntries[key] = cacheEntry{
		CreatedAt: time.Now().UTC(),
		Data:      value,
		Protected: protect,
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
		if !val.Protected && (time.Since(val.CreatedAt) > c.interval) {
			delete(c.CachedEntries, key)
		}
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.CachedEntries, key)
}

// DeleteAllStartsWith removes all cached entries whose
// key starts with the given prefix.
// This can be used to invalidate all entries related
// to a resource that may have had an instance updated or removed.
func (c *Cache) DeleteAllStartsWith(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.CachedEntries {
		if strings.HasPrefix(key, prefix) {
			delete(c.CachedEntries, key)
		}
	}
}

// Clear deletes all cached entries.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.CachedEntries)
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		CachedEntries: make(map[string]cacheEntry),
		interval:      interval,
		mu:            &sync.Mutex{},
	}

	go cache.reapLoop()

	return &cache
}
