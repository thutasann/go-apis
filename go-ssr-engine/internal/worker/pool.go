package worker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// WorkerPool is a bounded goroutine pool.
//
// Goals:
// - Fixed number of workers
// - Bounded queue
// - Backpressure
// - Graceful shutdown
// - No dynamic goroutine spawn
type WorkerPool struct {
	queue chan Job

	workers int

	wg sync.WaitGroup

	// atomic counters (avoid mutex)
	active   int64
	rejected int64

	ctx    context.Context
	cancel context.CancelFunc
}

// New creates worker pool.
//
// workers: number of goroutines
// queueSize: bounded queue capacity
func New(workers, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		queue:   make(chan Job, queueSize),
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start launches fixed workers.
func (p *WorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.workerLoop()
	}
}

// workerLoop processes jobs sequentially
// No spawning inside. One goroutine per worker.
func (p *WorkerPool) workerLoop() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return

		case job := <-p.queue:
			atomic.AddInt64(&p.active, 1)

			// Render directly to response writer
			_ = job.Tpl.RenderTo(job.Res, &job.Ctx)

			atomic.AddInt64(&p.active, -1)
		}
	}
}

// Submit enqueues job with backpressure control.
//
// timeout defines how long caller waits before rejecting
func (p *WorkerPool) Submit(job Job, timeout time.Duration) error {
	select {
	case p.queue <- job:
		return nil
	case <-time.After(timeout):
		atomic.AddInt64(&p.rejected, 1)
		return errors.New("queue full")
	case <-p.ctx.Done():
		return errors.New("worker pool shutting down")
	}
}

// Shutdown grecefully stops workers.
//
// Waits for in-flight jobs to complete.
func (p *WorkerPool) Shutdown() {
	p.cancel()
	p.wg.Wait()
	close(p.queue)
}

// Active returns currently processing jobs
func (p *WorkerPool) Active() int64 {
	return atomic.LoadInt64(&p.active)
}

// Rejected returns number of rejectd jobs.
func (p *WorkerPool) Rejected() int64 {
	return atomic.LoadInt64(&p.rejected)
}
