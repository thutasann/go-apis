package main

import (
	"fmt"
	"time"
)

func Priority_Queue() {
	vipQueue := make(chan string, 5)    // VIP waiting room
	normalQueue := make(chan string, 5) // Normal waiting room

	// Worker
	go func() {
		for {
			// always try VIP first (non-blocking)
			select {
			case vip := <-vipQueue:
				fmt.Println("Serving VIP: ", vip)
				time.Sleep(500 * time.Millisecond)
				continue
			default:
			}

			// If no VIP, serve normal (blocking)
			select {
			case normal := <-normalQueue:
				fmt.Println("Serving Normal:", normal)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	normalQueue <- "Alice"
	normalQueue <- "Bob"
	vipQueue <- "CEO"

	time.Sleep(3 * time.Second)
}
