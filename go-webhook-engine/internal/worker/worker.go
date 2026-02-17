package worker

import (
	"context"
	"log"

	"github.com/thutasann/go-webhook-engine/internal/domain"
)

func (p *Pool) startWorker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("worker shutting down")
			return
		default:
			eventID, err := p.queue.Dequeue(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				continue
			}

			event, err := p.repo.GetByID(ctx, eventID)
			if err != nil {
				continue
			}

			_ = p.repo.UpdateStatus(ctx, eventID, domain.StatusProcessing)

			// Simulate processing
			err = processEvent(event)

			if err != nil {
				_ = p.repo.IncrementRetry(ctx, eventID)
				_ = p.repo.UpdateStatus(ctx, eventID, domain.StatusFailed)
				continue
			}

			_ = p.repo.UpdateStatus(ctx, eventID, domain.StatusSuccess)
			log.Printf("processed event %s\n", event.ID.Hex())
		}
	}
}

func processEvent(event *domain.Event) error {
	// TODO: real business logic
	return nil
}
