package main

import "fmt"

func Sample_Buffered() {
	ch := make(chan int, 2)

	ch <- 1
	fmt.Println("Sent 1")

	ch <- 2
	fmt.Println("Sent 2")

	// ch <- 3 // DEADLOCK (buffer full)

	fmt.Println("Receiving...")
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
