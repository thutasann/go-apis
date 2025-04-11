package fundamentals

import "fmt"

var Hello string = "Hello"

// Defer Sample One
func DeferSampleOne() {
	fmt.Println("----> Mutex Examples")
	defer fmt.Println("A")
	fmt.Println("B")
}

// Defer Inside Loop
func DeferInsideLoop() {
	fmt.Println("---> Loop with Defer:")
	for i := 0; i < 3; i++ {
		defer fmt.Println("-> deferred:", i)
	}
}
