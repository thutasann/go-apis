package control

import (
	"sync"

	"github.com/thutasann/ctp/internal/telemetry"
)

// Aggregator maintains concurrent-safe metrics
type Aggregator struct {
	mu            sync.Mutex
	TotalMessages int
	PerStation    map[string]int
}

// Initialize a new Aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{
		PerStation: make(map[string]int),
	}
}

// Add ingests one telemetry message
func (a *Aggregator) Add(t telemetry.Telemetry) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.TotalMessages++
	a.PerStation[t.Station]++
}
