/*
# Mini RunTime
1. Limit how many goroutines run at the same time (e.g., 3 max)

2. Queue extra jobs if pool is full

3. Worker goroutines reuse

4. Manual .Submit() function to send tasks
*/
package goroutines

import (
	"fmt"
	"sync"
	"time"
)

// Pool is the mini runtime
type Pool struct {
	taskQueue chan func()
	wg        sync.WaitGroup
}

// NewPool to create a new Pool
func NewPool(maxWorkers int) *Pool {
	p := &Pool{
		taskQueue: make(chan func(), 100), // buffer for queued tasks
	}

	// spawn fixed number of workers
	for i := 0; i < maxWorkers; i++ {
		go p.runtime_worker(i)
	}

	return p
}

// Submit adds a task to the pool
func (p *Pool) Submit(task func()) {
	p.wg.Add(1)
	p.taskQueue <- task
}

// worker function that run tasks
func (p *Pool) runtime_worker(id int) {
	for task := range p.taskQueue {
		fmt.Printf("👷 Worker %d running task...\n", id)
		task()
		p.wg.Done()
	}
}

// Wait blocks until all tasks are finished
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Mini RunTime Program
func MiniRuntimeProgram() {
	pool := NewPool(3)

	// simuate 10 tasks
	for i := 1; i <= 10; i++ {
		n := i
		pool.Submit(func() {
			fmt.Printf("🔨 Processing task %d\n", n)
			time.Sleep(1 * time.Second)
		})
	}

	pool.Wait()
	fmt.Println("✅ All tasks completed!")
}
