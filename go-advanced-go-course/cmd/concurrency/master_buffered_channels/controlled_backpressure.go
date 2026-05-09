package main

import (
	"fmt"
	"time"
)

func Controlled_Backpressure() {
	ch := make(chan int, 2) // buffered size = 2

	// Consumer
	go func() {
		for v := range ch {
			time.Sleep(1 * time.Second)
			fmt.Println("Consumed: ", v)
		}
	}()

	// Producer
	for i := range 5 {
		fmt.Println("Producing: ", i)
		ch <- i // blocks only when buffer is full
		fmt.Println("Buffered len: ", len(ch))
	}
}
