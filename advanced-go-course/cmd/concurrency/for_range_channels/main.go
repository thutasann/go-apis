package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan int)
	go sell(ch, &wg)
	go buy(ch, &wg)
	wg.Wait()
}

func sell(ch chan int, wg *sync.WaitGroup) {
	ch <- 1
	ch <- 2
	fmt.Println("Sent all data")
	close(ch)
	wg.Done()
}

func buy(ch chan int, wg *sync.WaitGroup) {
	fmt.Println("Waiting for data")
	for val := range ch {
		fmt.Println("Received: ", val)
	}
	wg.Done()
}

func sample_for_range() {
	nums := []int{1, 2, 3, 4, 5}
	for i, item := range nums {
		fmt.Println(i, item)
	}
}
