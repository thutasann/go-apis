package main

import (
	"fmt"
	"time"
)

func Sample_Unbuffered() {
	ch := make(chan int)

	go func() {
		fmt.Println("Goroutine: waiting to send...")
		ch <- 10 // BLOCKS until someone receives
		fmt.Println("Goroutine: send completed")
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("Main: ready to receive")

	val := <-ch
	fmt.Println("Main: received", val)
}
