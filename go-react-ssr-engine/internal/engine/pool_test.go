package engine

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestPoolAcquireRelease(t *testing.T) {
	pool, err := NewPool(4)
	if err != nil {
		t.Fatalf("pool init: %v", err)
	}
	defer pool.Shutdown()

	// Acquire all 4 workers — should not block
	workers := make([]*Worker, 4)
	for i := 0; i < 4; i++ {
		workers[i] = pool.Acquire()
		if workers[i] == nil {
			t.Fatalf("worker %d is nil", i)
		}
	}

	// Release all back
	for _, w := range workers {
		pool.Release(w)
	}
}

func TestPoolConcurrency(t *testing.T) {
	// Verifies the pool handles high concurrency without deadlocks.
	// 4 workers, 100 goroutines — each goroutine acquires, does work, releases.
	pool, err := NewPool(4)
	if err != nil {
		t.Fatalf("pool init: %v", err)
	}
	defer pool.Shutdown()

	var wg sync.WaitGroup
	var completed atomic.Int64

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := pool.Acquire()
			// Simulate render work — the worker is exclusively ours here
			completed.Add(1)
			pool.Release(w)
		}()
	}

	wg.Wait()

	if completed.Load() != 100 {
		t.Fatalf("expected 100 completions, got %d", completed.Load())
	}
}

func TestPoolShutdownIdempotent(t *testing.T) {
	pool, err := NewPool(2)
	if err != nil {
		t.Fatalf("pool init: %v", err)
	}

	// Should not panic on double shutdown
	pool.Shutdown()
	pool.Shutdown()
}
