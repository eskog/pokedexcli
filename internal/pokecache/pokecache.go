package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	data      []byte
}

type Cache struct {
	entry map[string]cacheEntry
	mu    sync.Mutex
}

func NewCache(cacheDuration time.Duration) *Cache {
	cache := Cache{
		entry: make(map[string]cacheEntry),
		mu:    sync.Mutex{},
	}
	go cache.cleanLoop(cacheDuration)
	return &cache

}

func (c *Cache) cleanLoop(cacheDuration time.Duration) {
	ticker := time.NewTicker(cacheDuration)
	defer ticker.Stop()
	for {
		select {
		case tick := <-ticker.C:
			for url, cacheEntry := range c.entry {
				if tick.Sub(cacheEntry.createdAt) < time.Second*30 {
					delete(c.entry, url)
				}
			}
		}
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		data:      val,
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	data, exists := c.entry[key]
	if !exists {
		return []byte{}, false
	}
	fmt.Println("Using cache")
	return data.data, true
}
