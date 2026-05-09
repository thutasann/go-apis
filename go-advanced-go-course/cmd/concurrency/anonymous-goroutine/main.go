package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println(runtime.NumCPU())
	go func() {
		fmt.Println("In anonymous function")
	}()
	time.Sleep(1 * time.Second)
}
