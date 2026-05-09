package queue

import (
	"context"

	"github.com/thutasann/job-processor/internal/job"
)

// Queue holds two buffered channels for priority handling
type Queue struct {
	highCh   chan *job.Job
	normalCh chan *job.Job
}

// New creates a new priority queue.
// highSize and normalSize define buffer capacity.
func New(highSize, normalSize int) *Queue {
	return &Queue{
		highCh:   make(chan *job.Job, highSize),
		normalCh: make(chan *job.Job, normalSize),
	}
}

// Enqueue inserts a job into the appropriate buffer.
// this BLOCKS when the buffer is full (intentional backpressure)
func (q *Queue) Enqueue(j *job.Job) {
	switch j.Priority {
	case job.High:
		q.highCh <- j
	default:
		q.normalCh <- j
	}
}

// Dequeue selects a job with priority awareness.
// It respects context cancellation for graceful shutdown.
func (q *Queue) Dequeue(ctx context.Context) (*job.Job, error) {
	for {
		// First, try to pull HIGH without blocking.
		select {
		case j := <-q.highCh:
			return j, nil
		default:
		}

		// If no HIGH job immediately available,
		// block waiting for either HIGH or NORMAL.
		select {
		case j := <-q.highCh:
			return j, nil
		case j := <-q.normalCh:
			return j, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
