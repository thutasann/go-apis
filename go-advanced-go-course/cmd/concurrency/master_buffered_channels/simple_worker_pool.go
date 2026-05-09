package main

import (
	"fmt"
	"time"
)

func simple_worker(id int, jobs <-chan int) {
	for job := range jobs {
		fmt.Printf("Worker %d is processing job %d\n", id, job)
		time.Sleep(1 * time.Second)
	}
}

func Simple_Worker_Pool() {
	jobs := make(chan int, 5) // job queue

	// start 3 workers
	for i := 1; i <= 3; i++ {
		go simple_worker(i, jobs)
	}

	// send 10 jobs
	for j := 1; j <= 10; j++ {
		jobs <- j
	}

	close(jobs)

	time.Sleep(5 * time.Second)
}
