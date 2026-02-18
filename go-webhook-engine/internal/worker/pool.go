package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/thutasann/go-webhook-engine/internal/queue"
	"github.com/thutasann/go-webhook-engine/internal/repository"
)

type Pool struct {
	workerCount int
	minWorkers  int
	maxWorkers  int
	queue       queue.Queue
	repo        repository.EventRepository
	rateLimiter *RateLimiter

	jobs chan string

	wg sync.WaitGroup

	scaleMu      sync.Mutex
	activeWorker int
}

func NewPool(
	workerCount int,
	minWorkers int,
	maxWorkers int,
	q queue.Queue,
	r repository.EventRepository,
	rl *RateLimiter,
) *Pool {

	if workerCount < minWorkers {
		workerCount = minWorkers
	}

	if workerCount > maxWorkers {
		workerCount = maxWorkers
	}

	return &Pool{
		workerCount:  workerCount,
		minWorkers:   minWorkers,
		maxWorkers:   maxWorkers,
		queue:        q,
		repo:         r,
		rateLimiter:  rl,
		jobs:         make(chan string, 100),
		activeWorker: workerCount,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		// Dispatcher
		p.wg.Add(1)
		go p.dispatcher(ctx)

		// Start initial workers
		p.scaleWorkers(ctx, p.workerCount)

		// start scaler
		p.wg.Add(1)
		go p.scaler(ctx)
	}
}

// scaler monitors queue length and adjusts workers dynamically
func (p *Pool) scaler(ctx context.Context) {
	defer p.wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			queueLen := len(p.jobs)
			if queueLen > 50 && p.activeWorker < p.maxWorkers {
				p.scaleWorkers(ctx, 1) // add another one
				log.Printf("Scaling up: active workers=%d\n", p.activeWorker)
			} else if queueLen < 10 && p.activeWorker > p.minWorkers {
				// remove one worker by sending cancel context
				p.scaleWorkers(ctx, -1)
				log.Printf("Scaling down: active workers=%d\n", p.activeWorker)
			}
		}
	}
}

// scaleWorkers adjusts worker count dynamically (+ve = add, -ve = remove)
func (p *Pool) scaleWorkers(ctx context.Context, delta int) {
	p.scaleMu.Lock()
	defer p.scaleMu.Unlock()

	if delta > 0 {
		for i := 0; i < delta && p.activeWorker < p.maxWorkers; i++ {
			p.wg.Add(1)
			go p.worker(ctx)
			p.activeWorker++
		}
	} else if delta < 0 {
		// Just reduce the active worker counter; actual worker will exit naturally when jobs channel closes or context cancels
		p.activeWorker += delta
		if p.activeWorker < p.minWorkers {
			p.activeWorker = p.minWorkers
		}
	}
}

// Shutdown
func (p *Pool) Shutdown() {
	p.wg.Wait()
}

// dispatcher
func (p *Pool) dispatcher(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			close(p.jobs)
			return
		default:
			eventID, err := p.queue.Dequeue(ctx)
			if err != nil {
				continue
			}

			select {
			case p.jobs <- eventID:
			case <-ctx.Done():
				return
			}
		}
	}
}
