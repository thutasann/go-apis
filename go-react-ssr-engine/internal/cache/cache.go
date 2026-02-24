package cache

import (
	"container/list"
	"sync"
)

// LRU is a concurrent-safe Least Recently Used cache for rendered HTML.
// Under load, many requests hit the same pages. Caching avoids redundant
// V8 renders which are the most expensive operation in the pipeline.
//
// Design: sync.RWMutex + doubly-linked list + map.
// - Reads (cache hits) take a read lock — many goroutines read simultaneously.
// - Writes (cache misses, evictions) take a write lock — brief, infrequent.
// - LRU eviction ensures hot pages stay cached, cold pages get evicted.
type LRU struct {
	mu      sync.Mutex
	maxSize int
	items   map[string]*list.Element
	order   *list.List // front = most recent, back = least recent
}

type entry struct {
	key   string
	value string // rendered HTML
}

// NewLRU creates a cache with a max entry count.
// maxSize=0 creates a no-op cache (all Gets miss, all Sets are dropped).
func NewLRU(maxSize int) *LRU {
	return &LRU{
		maxSize: maxSize,
		items:   make(map[string]*list.Element, maxSize),
		order:   list.New(),
	}
}

// Get retrieves cached HTML for a cache key.
// Returns the HTML and true on hit, empty string and false on miss.
// Promotes the entry to front (most recent) on hit.
func (c *LRU) Get(key string) (string, bool) {
	if c.maxSize == 0 {
		return "", false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	elem, exists := c.items[key]
	if !exists {
		return "", false
	}

	// Move to front — this page is hot, keep it in cache
	c.order.MoveToFront(elem)
	return elem.Value.(*entry).value, true
}

// Set stores rendered HTML. Evicts the least recently used entry
// if cache is full. Overwrites if key already exists.
func (c *LRU) Set(key string, value string) {
	if c.maxSize == 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Update existing entry
	if elem, exists := c.items[key]; exists {
		c.order.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return
	}

	// Evict oldest if full
	if c.order.Len() >= c.maxSize {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.items, oldest.Value.(*entry).key)
		}
	}

	// Insert new entry at front
	e := &entry{key: key, value: value}
	elem := c.order.PushFront(e)
	c.items[key] = elem
}

// Delete removes a specific key. Used by invalidator on file change.
func (c *LRU) Delete(key string) {
	if c.maxSize == 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.order.Remove(elem)
		delete(c.items, key)
	}
}

// Flush clears the entire cache. Used on full rebuild in dev mode.
func (c *LRU) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element, c.maxSize)
	c.order.Init()
}

// Len returns current number of cached entries. For metrics/debug.
func (c *LRU) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}
