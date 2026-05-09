package main

import (
	"fmt"
	"time"
)

func Smart_Dropper() {
	events := make(chan string, 3) // limited capacity

	// consumer
	go func() {
		for e := range events {
			fmt.Println("Processing...", e)
			time.Sleep(1 * time.Second) // slow processor
		}
	}()

	send := func(event string) {
		select {
		case events <- event:
			fmt.Println("Queued:", event)
		default:
			fmt.Println("Dropped (overloaded):", event)
		}
	}

	for i := 1; i <= 6; i++ {
		send(fmt.Sprintf("event-%d", i))
	}

	time.Sleep(5 * time.Second)
}
