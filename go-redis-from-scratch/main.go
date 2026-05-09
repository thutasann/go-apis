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

	// --- Testing
	client := client.New("localhost:5001")
	for i := 0; i < 10; i++ {
		// SET
		if err := client.Set(context.TODO(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second)

		// GET
		val, err := client.Get(context.TODO(), fmt.Sprintf("foo_%d", i))
		if err != nil {
			log.Fatal("GET error --> ", err)
		}
		fmt.Printf("val :>> %s\n", val)
	}

	// The main goroutine starts the server in a separate goroutine.
	// The select {} statement blocks the main goroutine forever, keeping the program running.
	// Without it, main() would exit immediately and kill the child goroutine.
	// select {}
}
