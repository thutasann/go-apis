package main

import (
	"fmt"
	"time"
)

// basically shows that goroutines do not have parents or children
// the output is non-deterministic
func main() {
	go start()
	time.Sleep(1 * time.Second)
}

func start() {
	go process()
	fmt.Println("In start")
}

func process() {
	fmt.Println("In process")
}
