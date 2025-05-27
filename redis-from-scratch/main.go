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
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	client := client.New("localhost:5001")
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("foo_%d", i)
		val := fmt.Sprintf("bar_%d", i)

		if err := client.Set(context.TODO(), key, val); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("data :>> ", server.kv.data)

	time.Sleep(time.Second)

	// The main goroutine starts the server in a separate goroutine.
	// The select {} statement blocks the main goroutine forever, keeping the program running.
	// Without it, main() would exit immediately and kill the child goroutine.
	// select {}
}
