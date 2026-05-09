package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(10)
	for i := range 10 {
		go calculateSquare(i, &wg)
	}
	elapsed := time.Since(start)
	wg.Wait()
	fmt.Println("Function took: ", elapsed)
}

func calculateSquare(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(i * i)
}
