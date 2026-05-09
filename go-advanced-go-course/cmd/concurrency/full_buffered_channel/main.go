package main

import (
	"fmt"
	"time"
)

// Example: buffer gets FULL → sender blocks
// Sending 1
// Sent 1
// Sending 2
// Sent 2
// Sending 3 (this WILL block)
// <-- PAUSE here for ~2 seconds -->
// Receiving one value
// Received: 1
// Sent 3 (unblocked)
// Receiving remaining values
// 2
// 3
func main() {
	ch := make(chan int, 2) // buffer size = 2

	fmt.Println("Sending 1")
	ch <- 1 // goes in buffer (1/2)
	fmt.Println("Sent 1")

	fmt.Println("Sending 2")
	ch <- 2 // goes in buffer (2/2, FULL)
	fmt.Println("Sent 2")

	// This goroutine tries to send AFTER buffer is full
	go func() {
		fmt.Println("Sending 3 (this WILL block)")
		ch <- 3 // BLOCKS here until someone receives
		fmt.Println("Sent 3 (unblocked)")
	}()

	// Give time to show that send is blocked
	time.Sleep(2 * time.Second)

	fmt.Println("Receiving one value")
	v := <-ch // frees 1 slot
	fmt.Println("Received:", v)

	// Now the blocked send can continue
	time.Sleep(1 * time.Second)

	fmt.Println("Receiving remaining values")
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
