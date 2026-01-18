package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents an item in the queue
type Job struct {
	ID int
}

const (
	queueSize    = 5
	numWorkers   = 3
	numProducers = 2
	totalJobs    = 20
)

func main() {
	rand.Seed(time.Now().UnixNano())

	jobQueue := make(chan Job, queueSize)
	done := make(chan struct{})

	var produced int
	var consumed int
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Start visualization
	go visualize(jobQueue, &produced, &consumed, &mu, done)

	// Start workers (consumers)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobQueue, &consumed, &mu, &wg)
	}

	// Start producers
	for p := 1; p <= numProducers; p++ {
		go producer(p, jobQueue, &produced, &mu)
	}

	// Let producers work
	time.Sleep(5 * time.Second)

	// Stop producing and close queue
	close(jobQueue)

	// Wait for workers to finish
	wg.Wait()

	// Stop visualization
	close(done)

	fmt.Println("\n✅ All jobs processed. Exiting.")
}
