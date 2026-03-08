package main

import (
	"fmt"
	"time"
)

/*
LESSON 5: Batch Processing with Slices

Goal:
Collect events into batches before processing them.

Benefits:
- fewer DB calls
- fewer network calls
- better CPU cache usage
- much higher throughput

Pattern used in:
Kafka consumers, log processors, analytics pipelines.
*/

type Event struct {
	ID   int
	Data string
}

/*
simulateDBInsert simulates a database bulk insert.

In real system this might be:

INSERT INTO events VALUES (...), (...), (...);
*/
func simulateDBInsert(batch []Event) {
	fmt.Printf("DB Insert: batch size = %d\n", len(batch))

	time.Sleep(100 * time.Millisecond)
}

/*
processEvents processes events using batching

batchSize controls when we flush
*/
func processEvents(events []Event, batchSize int) {

	// preallocate batch size
	batch := make([]Event, 0, batchSize)

	for _, e := range events {
		batch = append(batch, e)

		// when batch is full -> process it
		if len(batch) == batchSize {
			simulateDBInsert(batch)

			/*
				IMPORTANT:
				reset slice length but keep capacity

				This avoids allocation of a new slice.
			*/
			batch = batch[:0]
		}
	}

	/*
		Flush remaining events
	*/
	if len(batch) > 0 {

		simulateDBInsert(batch)
	}
}

func Slice_Batch_Processing() {
	/*
		Simulate incoming events
	*/

	totalEvents := 23

	events := make([]Event, totalEvents)

	for i := 0; i < totalEvents; i++ {

		events[i] = Event{
			ID:   i + 1,
			Data: fmt.Sprintf("event-%d", i+1),
		}
	}

	fmt.Println("Starting batch processing\n")

	processEvents(events, 5)

	fmt.Println("\nProcessing complete")
}
