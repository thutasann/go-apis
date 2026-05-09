package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val       []byte    // cache entry value
	createdAt time.Time // cache entry created time
}

type Cache struct {
	cache map[string]cacheEntry // Cache Entry (Key, Value) Map
	mux   *sync.Mutex           // For Locking goroutines
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mux:   &sync.Mutex{},
	}
	go c.ReapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now().UTC(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	cacheE, ok := c.cache[key]
	return cacheE.val, ok
}

func (c *Cache) ReapLoop(interval time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.Reap(interval)
	}
}

func (c *Cache) Reap(interval time.Duration) {
	timeAgo := time.Now().UTC().Add(-interval)
	for k, v := range c.cache {
		if v.createdAt.Before(timeAgo) {
			delete(c.cache, k)
		}
	}
}
