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

	wg sync.WaitGroup
}

func NewPool(workerCount int, q queue.Queue, r repository.EventRepository) *Pool {
	return &Pool{
		workerCount: workerCount,
		queue:       q,
		repo:        r,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.startWorker(ctx)
	}
}

func (p *Pool) Shutdown() {
	p.wg.Wait()
}
