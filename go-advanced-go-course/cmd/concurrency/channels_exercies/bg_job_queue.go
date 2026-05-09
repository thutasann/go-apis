package main

import (
	"fmt"
	"time"
)

func queue_worker(id int, jobs chan string) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job: %s\n", id, job)
		time.Sleep(2 * time.Second)
		fmt.Printf("Worker %d finished job: %s\n", id, job)
	}
}

func Background_Job_Queue() {
	jobs := make(chan string, 5) // buffered queue

	// start workers
	for i := 1; i <= 3; i++ {
		go queue_worker(i, jobs)
	}

	// simulate incoming requests
	requests := []string{
		"resize-image",
		"send-email",
		"generate-pdf",
		"compress-video",
		"backup-db",
	}

	for _, req := range requests {
		fmt.Println("Main: queued job:", req)
		jobs <- req // does NOT block until buffer is full
	}

	close(jobs)

	time.Sleep(8 * time.Second)
	fmt.Println("All jobs processed")
}
