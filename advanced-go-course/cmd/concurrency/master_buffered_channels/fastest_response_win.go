package main

import (
	"fmt"
	"time"
)

func fast_worker(name string, delay time.Duration, out chan string) {
	time.Sleep(delay)
	out <- name // buffered prevents blocking
}

func Fastest_Response_Win() {
	results := make(chan string, 3) // enough buffer for all workers

	go fast_worker("us-east", 3*time.Second, results)
	go fast_worker("us-west", 1*time.Second, results)
	go fast_worker("asia", 2*time.Second, results)

	select {
	case fastest := <-results:
		fmt.Println("Using fastest:", fastest)

	case <-time.After(2 * time.Second):
		fmt.Println("All services too slow")
	}

	time.Sleep(3 * time.Second)

}
