package main

import (
	"context"
	"fmt"
	"time"
)

func cancellable_worker(ctx context.Context, id int, jobs chan int) {
	for {
		select {
		case job := <-jobs:
			fmt.Printf("Worker %d processing job %d\n", id, job)
			time.Sleep(1 * time.Second)
		case <-ctx.Done():
			fmt.Printf("Worker %d stopping: %v\n", id, ctx.Err())
			return
		}
	}
}

func Cancellable_Workers() {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan int)

	for i := 1; i <= 3; i++ {
		go cancellable_worker(ctx, i, jobs)
	}

	for i := 1; i <= 5; i++ {
		jobs <- i
	}

	time.Sleep(3 * time.Second)
	fmt.Println("Main: shutting down workers")
	cancel()

	time.Sleep(1 * time.Second)
	fmt.Println("Shutdown complete")
}
