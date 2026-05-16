package main

import (
	"fmt"
	"time"
)

func source(name string, delay time.Duration, out chan string) {
	time.Sleep(delay)
	out <- name + " event"
}

func Select_Fan_In() {
	out := make(chan string)

	go source("HTTP", 1*time.Second, out)
	go source("Kafka", 3*time.Second, out)
	go source("Cron", 2*time.Second, out)

	for i := 0; i < 3; i++ {
		select {
		case msg := <-out:
			fmt.Println("Received:", msg)
		default:
			fmt.Println("Default event")
		}
	}
}
