package main

import (
	"fmt"
	"time"
)

func logger(logs chan string) {
	for msg := range logs {
		fmt.Println("LOG:", msg)
		time.Sleep(500 * time.Millisecond)
	}
}

// You don’t want your app to slow down just because logs are slow.
func Loggin_System() {
	// Buffer allows app to continue without waiting on logger
	logs := make(chan string, 5)

	go logger(logs)

	for i := 1; i <= 10; i++ {
		logs <- fmt.Sprintf("event #%d", i)
		fmt.Println("Generated event:", i)
	}

	close(logs)
	time.Sleep(3 * time.Second)
}
