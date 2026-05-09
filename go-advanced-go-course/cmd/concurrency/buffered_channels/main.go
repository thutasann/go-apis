package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan int, 2)
	go sell(ch, &wg)
	wg.Wait()
}

func sell(ch chan int, wg *sync.WaitGroup) {
	ch <- 10
	ch <- 11
	ch <- 12
	go buy(ch, wg)
	fmt.Println("Sent all data to the channel")
	wg.Done()
}

func buy(ch chan int, wg *sync.WaitGroup) {
	fmt.Println("Watiging for data")
	fmt.Println("Received data: ", <-ch)
	wg.Done()
}
