package goroutines

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Basic Event Loop Style Using Select
func BasicEventLoopStyleUsingSelect() {
	ticker := time.NewTicker(1 * time.Second)
	events := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		events <- "🚀 Event A"
		time.Sleep(3 * time.Second)
		events <- "🚀 Event B"
	}()

	for {
		select {
		case e := <-events:
			fmt.Println("Got event: ", e)
		case t := <-ticker.C:
			fmt.Println("⏰ Tick at", t.Format("15:04:05"))
		default:
			// this keeps loop non-blocking
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func do_concurrent_tasks(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("starting::", name)
	time.Sleep(2 * time.Second)
	fmt.Println("Done:", name)
}

// Concurrent Task Queue (like Promise.all)
func ConcurrentTaskQueue() {
	var wg sync.WaitGroup
	tasks := []string{"Task1", "Task2", "Task3"}

	for _, task := range tasks {
		wg.Add(1)
		go do_concurrent_tasks(task, &wg)
	}

	wg.Wait()
	fmt.Println("All tasks complete")
}

func message_queue_worker(id int, jobs <-chan string) {
	for job := range jobs {
		fmt.Printf("Worker %d handlng job: %s\n", id, job)
		time.Sleep(1 * time.Second)
	}
}

// Simulating a Message Queue System
func MessageQueueSystem() {
	jobQueue := make(chan string, 5)

	// Spawn workers
	for i := 1; i <= 3; i++ {
		go message_queue_worker(i, jobQueue)
	}

	// Send jobs
	for j := 1; j <= 6; j++ {
		jobQueue <- fmt.Sprintf("📦 Job %d", j)
	}

	close(jobQueue)

	time.Sleep(5 * time.Second) // wait for all workers
}

func context_timeout_task(ctx context.Context) {
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Task finished")
	case <-ctx.Done():
		fmt.Println("🛑 Task cancelled:", ctx.Err())
	}
}

// Manual Context Timeout (like AbortController)
func ManualContextTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go context_timeout_task(ctx)
	time.Sleep(3 * time.Second)
	fmt.Println("Main exit")
}

func long_job() {
	time.Sleep(10 * time.Second)
}

// Debugging for GoRoutines
func GoRoutineDebugging() {
	for i := 0; i < 5; i++ {
		go long_job()
	}

	time.Sleep(1 * time.Second)
	fmt.Println("🧠 Number of goroutines:", runtime.NumGoroutine())
}

func worker_square(n int, wg *sync.WaitGroup, out chan<- int) {
	defer wg.Done()
	out <- n * n
}

// Goroutine Fan-Out, Fan-In Pattern
//
// - use `wg` to coordinate when all the `worker_square` goroutines are done
//
// - without `wg.Wait()`, you could close `out` too early
//
// - goroutine might still be writing to `out` - panic: "send on closed channel!"
//
// - so you wait for all goroutines to finish (`wg.Wait()`), then you `close(out)` safely
//
// - so `wg` is there to block the main goroutine untile all workers are done
//
// - WG Under the hood
//
// - When you Add(1), it increments the counter.
//
// - When you Done(), it decrements the counter.
//
// - When you Wait(), it blocks until the counter == 0.
//
// - `Wait()` uses a blocking mechanism (atomic counters + semaphores) efficiently
//
// - It is blocking, but only the goroutine that calls `Wait()` is blocked - other goroutines keep running freely
func GoRoutinesFanOutFanIn() {
	numbers := []int{2, 3, 4, 5}
	out := make(chan int, len(numbers))
	var wg sync.WaitGroup

	// fan-out: spawan goroutines
	for _, num := range numbers {
		wg.Add(1)
		go worker_square(num, &wg, out)
	}

	wg.Wait()
	close(out)

	// fan-in: collect results
	sum := 0
	for result := range out {
		sum += result
	}

	fmt.Println("sum of squares: ", sum)
}

func slow_server(response chan<- string) {
	time.Sleep(3 * time.Second)
	response <- "server responsed after 3s"
}

// Timeout Using `select` (simulate request timeouts)
func TimeoutUsingSelect() {
	response := make(chan string)
	go slow_server(response)

	select {
	case res := <-response:
		fmt.Println("Received:", res)
	case <-time.After(2 * time.Second):
		fmt.Println("⏰ Timeout! No response in time")
	}
}
