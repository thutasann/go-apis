package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// worker consumes jobs from the queue
func worker(id int, queue <-chan Job, consumed *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range queue {
		time.Sleep(time.Duration(rand.Intn(800)+400) * time.Millisecond)

		mu.Lock()
		*consumed++
		mu.Unlock()

		fmt.Printf("👷 Worker %d processed Job %d\n", id, job.ID)
	}
}
