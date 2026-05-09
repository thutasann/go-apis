package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/thutasann/job-processor/internal/job"
	"github.com/thutasann/job-processor/internal/queue"
	"github.com/thutasann/job-processor/internal/worker"
)

func main() {
	q := queue.New(5, 5)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start 3 workers
	for i := 1; i <= 3; i++ {
		w := &worker.Worker{
			ID:    i,
			Queue: q,
			WG:    &wg,
		}
		wg.Add(1)
		w.Start(ctx)
	}

	// Enqueue jobs
	for i := 1; i <= 6; i++ {
		j := job.New(
			string(rune('A'+i)),
			job.Normal,
			func() error {
				time.Sleep(1 * time.Second)
				return nil
			},
		)
		q.Enqueue(j)
	}

	time.Sleep(3 * time.Second)

	// Graceful shutdown
	log.Println("Stopping system...")
	cancel()

	wg.Wait()

}
