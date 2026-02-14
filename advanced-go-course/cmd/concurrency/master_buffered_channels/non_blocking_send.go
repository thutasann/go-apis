package main

import "fmt"

func Non_Blocking_Send() {
	// Suggestion box with space for 2 messages
	events := make(chan string, 2)

	sendEvent := func(event string) {
		select {
		case events <- event:
			fmt.Println("Stored: ", event)
		default:
			fmt.Println("Dropped: ", event)
		}
	}

	sendEvent("A")
	sendEvent("B")
	sendEvent("C") // this one will be dropped
}
