package main

import (
	"fmt"
	"sync"
	"time"
)

// Visualize prints the queue stats periodically
func visualize(queue chan Job, produced, consumed *int, mu *sync.Mutex, done <-chan struct{}) {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			mu.Lock()
			p := *produced
			c := *consumed
			q := len(queue)
			mu.Unlock()

			drawDashboard(p, c, q)
		}
	}
}

// drawDashboard renders a simple terminal UI
func drawDashboard(produced, consumed, queueLen int) {
	fmt.Print("\033[H\033[2J") // clear terminal

	fmt.Println("📦 Channel-Based Queue Visualization")
	fmt.Println("===================================")
	fmt.Printf("Produced Jobs : %d\n", produced)
	fmt.Printf("Consumed Jobs : %d\n", consumed)
	fmt.Printf("In Queue      : %d\n", queueLen)

	fmt.Print("Queue: [")
	for i := 0; i < queueLen; i++ {
		fmt.Print("■")
	}
	for i := queueLen; i < queueSize; i++ {
		fmt.Print(" ")
	}
	fmt.Println("]")
}
