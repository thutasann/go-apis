package worker

import (
	"context"
	"log"

	"github.com/thutasann/go-webhook-engine/internal/domain"
	"github.com/thutasann/go-webhook-engine/internal/queue"
)

func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("worker shutting down")
			return
		case eventID, ok := <-p.jobs:
			if !ok {
				return
			}

			p.handleEvent(ctx, eventID)
		}
	}
}

func (p *Pool) handleEvent(ctx context.Context, eventID string) {
	// Wait for rate limiter token
	if err := p.rateLimiter.Wait(ctx); err != nil {
		return
	}

	event, err := p.repo.GetByID(ctx, eventID)
	if err != nil {
		return
	}

	_ = p.repo.UpdateStatus(ctx, eventID, domain.StatusProcessing)

	err = processEvent(event)

	if err != nil {
		p.retryOrDLQ(ctx, eventID, event)
		return
	}

	_ = p.repo.UpdateStatus(ctx, eventID, domain.StatusSuccess)
}

func (p *Pool) retryOrDLQ(ctx context.Context, eventID string, event *domain.Event) {
	_ = p.repo.IncrementRetry(ctx, eventID)

	if event.RetryCount+1 >= event.MaxRetries {
		_ = p.repo.UpdateStatus(ctx, eventID, domain.StatusFailed)

		if redisQueue, ok := p.queue.(*queue.RedisQueue); ok {
			_ = redisQueue.EnqueueDLQ(ctx, eventID)
		}

		log.Printf("event %s moved to DLQ\n", eventID)
		return
	}

	_ = p.queue.Enqueue(ctx, eventID)
}

func processEvent(event *domain.Event) error {
	// simulate real processing
	return nil
}
