package concurrencypatterns

import (
	"fmt"
	"sync"
	"time"
)

// SafeCache is a concurrent in-memory key-value store.
type SafeCache struct {
	mu    sync.Mutex
	store map[string]string
}

// Set stores a key-value pair safely
func (c *SafeCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
	fmt.Printf("✅ Set %s = %s\n", key, value)
}

// Get safely retrieves a value by key
func (c *SafeCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, exists := c.store[key]
	return val, exists
}

// Thread Safe Sample
func ThreadSafeSample() {
	cache := SafeCache{
		store: make(map[string]string),
	}

	var wg sync.WaitGroup

	// simulate multiple goroutines writing to the cache
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("user:%d", id)
			cache.Set(key, fmt.Sprintf("User%dData", id))
		}(i)
	}

	// simulate multiple goroutines reading from the cache
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			key := fmt.Sprintf("user:%d", id)
			if val, ok := cache.Get(key); ok {
				fmt.Printf("🔎 Got %s = %s\n", key, val)
			} else {
				fmt.Printf("❌ %s not found\n", key)
			}
		}(i)
	}

	wg.Wait()
}
