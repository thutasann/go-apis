/*
1. `for-select` is a Go idiom to listen on multiple channels continuously

2. react when one of them sends a value, It's used for concurrent control loops
*/
package forselect

import (
	"fmt"
	"time"
)

func TimeoutWithSelect() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "done"
	}()

	for {
		select {
		case msg := <-ch:
			fmt.Println("Received: ", msg)
		case <-time.After(3 * time.Second):
			fmt.Println("Timeout!")
			return
		}
	}
}

func GracefulShutdownWithQuitChannel() {
	messages := make(chan string)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case msg := <-messages:
				fmt.Println("Received: ", msg)
			case <-quit:
				fmt.Println("Shutting down....")
				return
			}
		}
	}()

	messages <- "hello"
	quit <- struct{}{}
}
