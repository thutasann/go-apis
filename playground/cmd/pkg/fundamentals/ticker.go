package fundamentals

import (
	"fmt"
	"time"
)

func TickerSampleOne() {
	fmt.Println("---> Ticker Sample One")

	// Create a ticker that ticks every 2 seconds
	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at:", t)
		}
	}()

	// Run for 10 seconds before stoppping
	time.Sleep(10 * time.Second)
	ticker.Stop()
	fmt.Println("Ticker Stopped")
}
