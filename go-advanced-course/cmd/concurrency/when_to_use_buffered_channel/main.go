// - You have 3 workers preparing food
//
// - You have 1 cashier packing orders
//
// - If workers must wait for the cashier every time → everyone gets blocked
//
// - Instead, workers put food on a counter (buffer) and continue cooking
//
// - Cashier picks items from the counter when ready
//
// - That “counter” is a buffered channel
package main

import (
	"fmt"
	"time"
)

func main() {
	// BUFFERED CHANNEL with capacity 3
	// This means:
	// - Up to 3 values can be sent WITHOUT blocking
	// - Sender blocks only when buffer is FULL
	results := make(chan string, 3)

	// Start 3 worker goroutines
	for i := 1; i <= 3; i++ {
		go worker(i, results)
	}

	// Consumer (main goroutine)
	// Reads results slowly to simulate heavy processing
	for i := 1; i <= 3; i++ {
		// THIS RECEIVE will block only if buffer is EMPTY
		result := <-results
		fmt.Println("Consumer received: ", result)

		// Slow consumer
		time.Sleep(2 * time.Second)
	}

	fmt.Println("All Work done")
}

func worker(id int, results chan<- string) {
	fmt.Println("Worker", id, "started work")

	// Simulate fast work
	time.Sleep(500 * time.Millisecond)

	// SEND into buffered channel
	// This will NOT block unless buffer is full
	results <- fmt.Sprintf("result from worker %d", id)

	fmt.Println("Worker", id, "sent result and continues")
}
