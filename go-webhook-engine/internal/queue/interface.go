package queue

import "context"

type Queue interface {
	Enqueue(ctx context.Context, eventID string) error

	// Blocking dequeue
	Dequeue(ctx context.Context) (string, error)
}
