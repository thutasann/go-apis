package main

import (
	"math/rand"
	"sync"
	"time"
)

// Producer pushes jobs into the channel-based queue
func producer(id int, queue chan<- Job, produced *int, mu *sync.Mutex) {
	for {
		mu.Lock()
		if *produced >= totalJobs {
			mu.Unlock()
			return
		}
		*produced++
		jobID := *produced
		mu.Unlock()

		queue <- Job{ID: jobID}
		time.Sleep(time.Duration(rand.Intn(400)+200) * time.Millisecond)
	}
}
