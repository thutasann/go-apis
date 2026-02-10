package main

import (
	"fmt"
	"time"
)

func Backpressure_With_Unbuffered_Channel() {
	ch := make(chan int) // unbuffered

	// consumer (slow)
	go func() {
		for v := range ch {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Consumer: ", v)
		}
	}()

	// Producer (fast)
	for i := range 5 {
		fmt.Println("Producing: ", i)
		ch <- i // BLOCKS until consumer receives
	}
}
