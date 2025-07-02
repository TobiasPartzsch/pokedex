package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	newEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = newEntry
	return true
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.entries {
			if entry.createdAt.Add(c.interval).Before(time.Now()) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  map[string]cacheEntry{},
		interval: interval,
	}
	go c.reapLoop()
	return c
}
