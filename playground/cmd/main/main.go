package main

import (
	"fmt"
	"time"

	"github.com/thutasann/playground/cmd/pkg/fundamentals"
)

func sayhello() {
	fmt.Println("Hello, Go!")
}

// Playground Main
func main() {
	fmt.Println("----- Playground -----")
	fmt.Println(fundamentals.Hello)
	go sayhello()
	time.Sleep(time.Second)
}
