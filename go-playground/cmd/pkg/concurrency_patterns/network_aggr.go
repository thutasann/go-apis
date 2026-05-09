package concurrencypatterns

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// NetworkData represents incoming data from a source
type NetworkData struct {
	SourceID int
	Value    int
}

// networkSource simulates a TCP/peer data stream
func networkSource(
	sourceID int,
	out chan<- NetworkData,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(rand.Intn(400)+200) * time.Millisecond)

		out <- NetworkData{
			SourceID: sourceID,
			Value:    rand.Intn(100),
		}
	}
}

// fanIn merges multiple producers into one channel
func fanIn(
	wg *sync.WaitGroup,
	out chan NetworkData,
) {
	wg.Wait()
	close(out)
}

func Network_Aggr() {
	numSources := 3
	dataChan := make(chan NetworkData)

	var wg sync.WaitGroup

	// fan-out: start sources
	for i := 1; i <= numSources; i++ {
		wg.Add(1)
		go networkSource(i, dataChan, &wg)
	}

	// fan-in: close channel when all resources are done
	go fanIn(&wg, dataChan)

	// single consumer
	for data := range dataChan {
		fmt.Printf(
			"Received from source: %d: value: %d\n",
			data.SourceID,
			data.Value,
		)
	}
}
