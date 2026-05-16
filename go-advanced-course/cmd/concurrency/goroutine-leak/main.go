package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go leak(&wg)
	wg.Wait()
}

func leak(wg *sync.WaitGroup) {
	ch := make(chan int)

	go func() {
		val := <-ch
		fmt.Println("Received: ", val)
		wg.Done()
	}()

	fmt.Println("Exiting leak method")
	wg.Done()
}
