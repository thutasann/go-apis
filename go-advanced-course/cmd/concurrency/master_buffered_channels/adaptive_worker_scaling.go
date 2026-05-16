package main

import (
	"fmt"
	"sync"
	"time"
)

func scale_worker(id int, jobs chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d processing %d\n", id, job)
		time.Sleep(1 * time.Second)
	}
}

func Adaptive_Worker_Scaling() {
	jobs := make(chan int, 10) // buffered queue
	var wg sync.WaitGroup

	workerCount := 1

	// Auto scaler
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)

			// If queue is filling up, add worker
			if len(jobs) > 5 && workerCount < 3 {
				workerCount++
				wg.Add(1)
				go scale_worker(workerCount, jobs, &wg)
				fmt.Println("Scaling up worker:", workerCount)
			}
		}
	}()

	// Start initial worker
	wg.Add(1)
	go scale_worker(1, jobs, &wg)

	// Produce jobs
	for i := 1; i <= 15; i++ {
		jobs <- i
	}

	close(jobs)
	wg.Wait()
}
