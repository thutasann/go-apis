package main

import "fmt"

// No Backpressure (memory risk)
func Without_BackPressure() {
	queue := make([]int, 0)

	// Producer (fast)
	for i := range 1_000_000 {
		queue = append(queue, i) // unbounded growth
	}

	fmt.Println("Queue Size: ", len(queue))
}
