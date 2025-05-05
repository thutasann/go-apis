package concurrencypatterns

import (
	"fmt"
	"sync"
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

func ThreadSafeSample() {

}
