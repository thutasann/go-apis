package main

import (
	"fmt"
	"time"
)

// two goroutines and use a channel to communicate between them.
func main() {
	ch := make(chan string)

	go sell(ch)
	go buy(ch)

	time.Sleep(2 * time.Second)
}

// sends data to the channel
func sell(ch chan string) {
	ch <- "Furniture"
	fmt.Println("Sent data to the channel")
}

// receive data from the channel
func buy(ch chan string) {
	fmt.Println("Waiting for Data...")
	val := <-ch
	fmt.Println("Received data - ", val)
}
