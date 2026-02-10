package main

import (
	"fmt"
	"time"
)

func Rate_Limiter() {
	// buffer size defines max allowed bursts
	tokens := make(chan struct{}, 3)

	// Refil tokens every second
	go func() {
		ticker := time.NewTicker(1 * time.Second)

		for range ticker.C {
			select {
			case tokens <- struct{}{}:
				// token added
			default:
				// buffer full, skip
			}
		}
	}()

	for i := 1; i <= 10; i++ {
		<-tokens // wait for token
		fmt.Println("Request allowed:", i, time.Now().Format("15:04:05"))
	}
}
