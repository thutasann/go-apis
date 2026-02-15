package worker

import (
	"context"
	"log"
	"sync"

	"github.com/thutasann/job-processor/internal/queue"
)

// Worker represents a single worker instance.
type Worker struct {
	ID    int
	Queue *queue.Queue
	WG    *sync.WaitGroup
}

// Start lunches the worker loop in a goroutine
func (w *Worker) Start(ctx context.Context) {
	go func() {
		defer w.WG.Done()

		log.Printf("[Worker %d] Started\n", w.ID)

		for {
			// Fetch job with context awareness
			job, err := w.Queue.Dequeue(ctx)
			if err != nil {
				log.Printf("[Worker %d] Shutting down\n", w.ID)
				return
			}

			w.process(job)
		}
	}()
}

// process executes the job safely.
func (w *Worker) process(j interface {
	String() string
	Execute() error
}) {
	defer func() {
		// Recover from panic to prevent worker crash
		if r := recover(); r != nil {
			log.Printf("[Worker %d] Recovered from panic in job %s: %v\n",
				w.ID, j.String(), r)
		}
	}()

	log.Printf("[Worker %d] Processing %s\n", w.ID, j.String())

	if err := j.Execute(); err != nil {
		log.Printf("[Worker %d] Job failed: %v\n", w.ID, err)
		return
	}

	log.Printf("[Worker %d] Completed %s\n", w.ID, j.String())
}
