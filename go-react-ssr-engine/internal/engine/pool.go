package engine

import (
	"fmt"
	"sync"
)

// Pool manages a fixed set of V8 workers via a buffered channel.
// Buffered channel = lock-free semaphore. No mutex in the hot path.
//
// Why Fixed pool:
//
// - Each V8 isolate costs ~ 10MB RAM. Unbounded = OOM under load.
//
// - Fixed pool = predictable memory: workerCount * 10MB.
//
// - Backpressure is free: when pool is empty, callers block on channel receive.
type Pool struct {
	workers chan *Worker
	size    int

	// all keeps references for shutdown cleanup.
	// Only touched at init and shutdown - no contention
	all  []*Worker
	once sync.Once
}

// NewPool creates `size` V8 workers and puts them in the channel.
// Returns error if any V8 isolate fails to initialize.
func NewPool(size int) (*Pool, error) {
	if size <= 0 {
		return nil, fmt.Errorf("pool: size must be > 0, got %d", size)
	}

	p := &Pool{
		workers: make(chan *Worker, size), // buffered = no lock on send/recv
		size:    size,
		all:     make([]*Worker, 0, size),
	}

	for i := 0; i < size; i++ {
		w, err := NewWorker(i)
		if err != nil {
			p.Shutdown() // Cleanup alrady-created workers on partial failure
			return nil, fmt.Errorf("pool: worker %d init failed: %w", i, err)
		}
		p.all = append(p.all, w)
		p.workers <- w // seed the pool
	}

	return p, nil
}

// Acquire blocks until a worker is available.
// Under normal load this returns instantly (buffered channel).
// Under heavy load this is the backpressure point — goroutines park here.
func (p *Pool) Acquire() *Worker {
	return <-p.workers
}

// Release returns a worker to the pool for reuse.
// Always call this in a defer after Acquire.
func (p *Pool) Release(w *Worker) {
	p.workers <- w
}

// Shutdown disposes all V8 isolates. Called once via sync.Once.
// Drains the channel first so no goroutine is stuck holding a dead worker.
func (p *Pool) Shutdown() {
	p.once.Do(func() {
		// Close channel so no new Acquires can happen
		close(p.workers)

		// Drain any works sitting in the channel
		for range p.workers {
			// Just drain
		}

		// Dispose every isolate we ever created
		for _, w := range p.all {
			w.Dispose()
		}
	})
}
