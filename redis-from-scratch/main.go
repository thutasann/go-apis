package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/thutasann/redisfromscratch/client"
)

// REDIS FROM SCRATCH
func main() {
	go func() {
		server := NewServer(Config{})
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		client := client.New("localhost:5001")
		key := fmt.Sprintf("foo-%d", i)
		val := fmt.Sprintf("bar-%d", i)

		if err := client.Set(context.TODO(), key, val); err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(time.Second)

	// The main goroutine starts the server in a separate goroutine.
	// The select {} statement blocks the main goroutine forever, keeping the program running.
	// Without it, main() would exit immediately and kill the child goroutine.
	// select {}
}
