package concurrencypatterns

import (
	"fmt"
	"time"
)

func someFunc(num string) {
	time.Sleep(3 * time.Second)
	fmt.Println(num)
}

// GoRoutines Sample One
func GoRoutineSampleOne() {
	go someFunc("this is num")
	fmt.Println("hi")
}
