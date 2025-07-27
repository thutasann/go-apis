package concurrencypatterns

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Config struct {
	sync.RWMutex
	data map[string]string
}

// Get retrieves a config value (concurrent safe)
func (c *Config) Get(key string) (string, bool) {
	c.RLock()         // Multiple readers allowed
	defer c.RUnlock() // Must unlock after read
	val, ok := c.data[key]
	return val, ok
}

// Set updates a config value (exclusive access)
func (c *Config) Set(key, value string) {
	c.Lock()         // Exclusive write access
	defer c.Unlock() // Must unlock after write
	log.Printf("[WRITER] Updating %s to %s", key, value)
	c.data[key] = value
}

// Simulate many readers
func simulateReader(id int, c *Config, key string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		val, ok := c.Get(key)
		if ok {
			fmt.Printf("[Reader %d] %s = %s\n", id, key, val)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// Simulate one writer
func simulateWriter(c *Config, key string, wg *sync.WaitGroup) {
	defer wg.Done()
	values := []string{"A", "B", "C", "D"}
	for _, v := range values {
		c.Set(key, v)
		time.Sleep(2 * time.Second) // Slow updates
	}
}

// [WRITER] Updating mode to init
// [Reader 1] mode = init
// [Reader 2] mode = init
// [Reader 3] mode = init
// [WRITER] Updating mode to A
// [Reader 1] mode = A
// [Reader 2] mode = A
func Config_Manager_RWMutex() {
	config := &Config{
		data: make(map[string]string),
	}
	config.Set("mode", "init")

	var wg sync.WaitGroup

	// Start readers
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go simulateReader(i, config, "mode", &wg)
	}

	// Start writer
	wg.Add(1)
	go simulateWriter(config, "mode", &wg)

	wg.Wait()
	fmt.Println("All goroutines completed.")
}
