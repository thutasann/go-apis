package main

import (
	"fmt"
	"time"
)

func graceful_worker(jobs chan int, quit chan struct{}) {
	for {
		select {
		case job := <-jobs:
			fmt.Println("Processing job: ", job)
			time.Sleep(1 * time.Second)
		case <-quit:
			fmt.Println("Worker: shutdown signal received")
			return
		}
	}
}

func Graceful_Shutdown() {
	jobs := make(chan int)
	quit := make(chan struct{})

	go graceful_worker(jobs, quit)

	for i := 1; i <= 5; i++ {
		jobs <- i
	}

	time.Sleep(3 * time.Second)
	fmt.Println("Main: sending shutdown signal")
	close(quit)

	time.Sleep(1 * time.Second)
	fmt.Println("Main: exited cleanly")
}
