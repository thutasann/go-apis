package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go goOne(ch1)
	go goTwo(ch2)

	select {
	case val1 := <-ch1:
		fmt.Println("val1: ", val1)
		break
		// fmt.Println("afer break") unreachable code
	case val2 := <-ch2:
		fmt.Println("val2: ", val2)
	default:
		fmt.Println("Executed default block")
	}

	time.Sleep(1 * time.Second)
}

func goOne(ch chan string) {
	ch <- "channel-1"
}

func goTwo(ch chan string) {
	ch <- "channel-2"
}
