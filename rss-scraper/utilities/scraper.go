package utilities

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/thutasann/rssagg/internal/database"
)

// Function to Start Scraping
func StartScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s Duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched: ", err)
		return
	}

	rssFeed, err := UrlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post: ", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
