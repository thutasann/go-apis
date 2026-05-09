package main

import (
	"fmt"
	"time"
)

func Worker_Pool_With_Cancellation() {
	jobs := make(chan int, 5)   // job queue
	stop := make(chan struct{}) // stop signal

	// worker
	go func() {
		for {
			select {
			case job := <-jobs:
				fmt.Println("Processing job: ", job)
				time.Sleep(1 * time.Second)

			case <-stop:
				fmt.Println("Worker stopping")
				return
			}
		}
	}()

	// send jobs
	for i := 1; i <= 3; i++ {
		jobs <- i
	}

	time.Sleep(2 * time.Second)

	// cancel everything
	close(stop)

	time.Sleep(1 * time.Second)
}
