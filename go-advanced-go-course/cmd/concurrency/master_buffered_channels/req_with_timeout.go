package main

import (
	"fmt"
	"time"
)

// select lets us:
// - Wait for a response
// - OR give up after a timeout
func Request_With_Timeout() {
	// small waiting room (cap 1)
	responseCh := make(chan string, 1)

	// worker (slow responder)
	go func() {
		time.Sleep(3 * time.Second)
		responseCh <- "here is your order"
	}()

	select {
	case res := <-responseCh:
		fmt.Println("Received: ", res)
	case <-time.After(2 * time.Second):
		fmt.Println("Request timed out")
	}
}
