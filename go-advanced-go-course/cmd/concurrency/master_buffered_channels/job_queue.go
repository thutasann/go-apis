package main

import (
	"fmt"
	"time"
)

// Producer faster than consumer
func Producer_Faster_Than_Consumer() {
	// Buffered channel with capacity 3
	// This means up to 3 jobs can wait in the queue
	jobs := make(chan int, 3)

	// Consumer (worker)
	go func() {
		for job := range jobs {
			fmt.Println("[consumer] Processing job: ", job)
			time.Sleep(1 * time.Second) // slow worker
		}
	}()

	// Producer
	for i := 1; i <= 5; i++ {
		fmt.Println("[producer] Sending job: ", i)
		jobs <- i // blocks only when buffer is full
	}

	close(jobs)
	time.Sleep(6 * time.Second)
}
