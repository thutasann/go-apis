package fundamentals

import (
	"fmt"
	"sync"
)

var mu sync.Mutex
var counter int

func increment() {
	mu.Lock()
	counter++
	defer mu.Unlock()
	fmt.Println("counter --> ", counter)
}

func MutexSamples() {
	fmt.Println("----> Mutex Examples")
	increment()
}
