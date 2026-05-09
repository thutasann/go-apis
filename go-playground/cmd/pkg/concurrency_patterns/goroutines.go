package concurrencypatterns

import (
	"fmt"
	"sync"
)

func someFunc(num string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(num)
}

// GoRoutines Sample One
func GoRoutineSampleOne() {
	var wg sync.WaitGroup

	wg.Add(3) // launching 3 goroutines
	go someFunc("1", &wg)
	go someFunc("2", &wg)
	go someFunc("3", &wg)

	wg.Wait() // wait for all goroutines to finish
	fmt.Println("hi")
}
