package concurrencypatterns

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Periodically polls a 3rd-party API (e.g., to sync data).
// Must handle graceful shutdown via context.
// Logs metrics every 5 seconds.
// Supports concurrent workers.

// Simulates external API call
func fetchData(ctx context.Context, id int) {
	select {
	case <-ctx.Done():
		log.Printf("[Worker %d] Cancelled", id)
	case <-time.After(2 * time.Second):
		log.Printf("[Worker %d] Fetched data at %s", id, time.Now().Format(time.RFC3339))
	}
}

func startWorker(ctx context.Context, id int, wg *sync.WaitGroup, ticker *time.Ticker) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Workder %d] Exiting...", id)
		case <-ticker.C:
			fetchData(ctx, id)
		}
	}
}

func Ticker_DataFetch_Sample() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown support
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	const numWorker = 3
	var wg sync.WaitGroup

	log.Println("Starting workers...")

	for i := 1; i <= numWorker; i++ {
		wg.Add(1)
		go startWorker(ctx, i, &wg, ticker)
	}

	<-stop // Wait for CTRL+C
	log.Println("Shutdown signal received")

	cancel()  // Cancel context to stop workers
	wg.Wait() // Wait for all to finish

	log.Println("Gracefully shut down.")
}
