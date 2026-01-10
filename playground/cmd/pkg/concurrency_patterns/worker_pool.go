package concurrencypatterns

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents work to be done
type WorkerJob struct {
	ID int
}

// Worker Result represents processed job output
type WorkerResult struct {
	JobID  int
	Output string
}

// worker processes jobs from the jobs channel
func wp_worker(
	workerID int,
	jobs <-chan WorkerJob,
	results chan<- WorkerResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)

		results <- WorkerResult{
			JobID:  job.ID,
			Output: fmt.Sprintf("Processed by worker %d", workerID),
		}
	}
}

func Worker_Pool_Sample() {
	const (
		numWorkers = 3
		numJobs    = 10
	)

	jobs := make(chan WorkerJob, numJobs)
	results := make(chan WorkerResult, numJobs)

	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go wp_worker(i, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- WorkerJob{ID: j}
	}
	close(jobs)

	// Close results when workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for res := range results {
		fmt.Printf("Job %d -> %s\n", res.JobID, res.Output)
	}

}
