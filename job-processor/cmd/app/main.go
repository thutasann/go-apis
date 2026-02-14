package main

import (
	"context"
	"fmt"
	"time"

	"github.com/thutasann/job-processor/internal/job"
	"github.com/thutasann/job-processor/internal/queue"
)

func main() {
	q := queue.New(2, 2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// enqueue
	q.Enqueue(job.New("1", job.Normal, nil))
	q.Enqueue(job.New("2", job.High, nil))
	q.Enqueue(job.New("3", job.Normal, nil))

	// dequeue
	for i := 0; i < 3; i++ {
		j, _ := q.Dequeue(ctx)
		fmt.Println("Dequeued:", j)
	}

	time.Sleep(time.Second)
}
