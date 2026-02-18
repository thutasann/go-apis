package worker

import (
	"context"
	"sync"

	"github.com/thutasann/go-webhook-engine/internal/queue"
	"github.com/thutasann/go-webhook-engine/internal/repository"
)

type Pool struct {
	workerCount int
	queue       queue.Queue
	repo        repository.EventRepository

	jobs chan string

	wg sync.WaitGroup

	rateLimiter *RateLimiter
}

func NewPool(workerCount int, q queue.Queue, r repository.EventRepository, rl *RateLimiter) *Pool {
	return &Pool{
		workerCount: workerCount,
		queue:       q,
		repo:        r,
		jobs:        make(chan string, 100),
		rateLimiter: rl,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		// Dispatcher
		p.wg.Add(1)
		go p.dispatcher(ctx)

		// Workers
		for i := 0; i < p.workerCount; i++ {
			p.wg.Add(1)
			go p.worker(ctx)
		}
	}
}

func (p *Pool) Shutdown() {
	p.wg.Wait()
}

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
