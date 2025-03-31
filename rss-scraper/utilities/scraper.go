package utilities

import (
	"log"
	"time"

	"github.com/thutasann/rssagg/internal/database"
)

// Function to Start Scraping
func StartScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s Duration", concurrency, timeBetweenRequest)
}
