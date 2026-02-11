package main

import (
	"fmt"
	"sync"
	"time"
)

func bp_worker(id int, jobs chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d handling job %d\n", id, job)
		time.Sleep(1 * time.Second)
	}
}

func Worker_Pool_With_Backpressure() {
	jobs := make(chan int, 4)
	var wg sync.WaitGroup

	// start two workers
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go bp_worker(i, jobs, &wg)
	}

	// Send jobs
	for i := 1; i <= 8; i++ {
		fmt.Println("Queueing job:", i)
		jobs <- i // blocks when queue is full
	}

	close(jobs)
	wg.Wait()
}
