package cache

import (
	"fmt"
	"sync"
	"testing"
)

func TestLRUBasic(t *testing.T) {
	c := NewLRU(3)

	c.Set("a", "html-a")
	c.Set("b", "html-b")
	c.Set("c", "html-c")

	// All present
	for _, key := range []string{"a", "b", "c"} {
		if _, ok := c.Get(key); !ok {
			t.Errorf("expected %s to be cached", key)
		}
	}

	// Insert 4th — evicts "a" (least recently used)
	c.Set("d", "html-d")

	if _, ok := c.Get("a"); ok {
		t.Error("expected 'a' to be evicted")
	}
	if _, ok := c.Get("d"); !ok {
		t.Error("expected 'd' to be cached")
	}
}

func TestLRUAccessPromotes(t *testing.T) {
	c := NewLRU(3)

	c.Set("a", "1")
	c.Set("b", "2")
	c.Set("c", "3")

	// Access "a" — moves it to front, "b" is now oldest
	c.Get("a")

	// Insert "d" — should evict "b" not "a"
	c.Set("d", "4")

	if _, ok := c.Get("a"); !ok {
		t.Error("expected 'a' to survive (was accessed recently)")
	}
	if _, ok := c.Get("b"); ok {
		t.Error("expected 'b' to be evicted (oldest)")
	}
}

func TestLRUDelete(t *testing.T) {
	c := NewLRU(10)
	c.Set("x", "data")
	c.Delete("x")

	if _, ok := c.Get("x"); ok {
		t.Error("expected 'x' to be deleted")
	}
}

func TestLRUFlush(t *testing.T) {
	c := NewLRU(10)
	c.Set("a", "1")
	c.Set("b", "2")
	c.Flush()

	if c.Len() != 0 {
		t.Errorf("expected empty cache after flush, got %d", c.Len())
	}
}

func TestLRUZeroSize(t *testing.T) {
	// Zero-size cache is a no-op — used in dev mode
	c := NewLRU(0)
	c.Set("a", "data")

	if _, ok := c.Get("a"); ok {
		t.Error("zero-size cache should always miss")
	}
}

func TestLRUConcurrent(t *testing.T) {
	// Hammer the cache from 100 goroutines to detect races.
	// Run with -race flag: go test -race ./internal/cache/
	c := NewLRU(100)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			c.Set(key, "value")
			c.Get(key)
			c.Delete(key)
		}(i)
	}

	wg.Wait()
}

func BenchmarkLRUGet(b *testing.B) {
	c := NewLRU(10000)
	// Pre-fill
	for i := 0; i < 10000; i++ {
		c.Set(fmt.Sprintf("k%d", i), "html-content-here")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			c.Get(fmt.Sprintf("k%d", i%10000))
			i++
		}
	})
}
