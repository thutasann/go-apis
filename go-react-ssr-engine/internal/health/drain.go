package health

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Drainer manages graceful shutdown.
// When shutdown is triggered:
// 1. Stop accepting new connections (handled by fasthttp.Shutdown)
// 2. Wait for in-flight requests to complete (tracked by activeConns)
// 3. Timeout if requests take too long
//
// This prevents dropped requests during rolling deploys.
type Drainer struct {
	checker  *Checker
	draining atomic.Bool
	timeout  time.Duration
}

func NewDrainer(checker *Checker, timeout time.Duration) *Drainer {
	return &Drainer{
		checker: checker,
		timeout: timeout,
	}
}

// IsDraining returns true during graceful shutdown.
// New requests should get 503 Service Unavailable.
func (d *Drainer) IsDraining() bool {
	return d.draining.Load()
}

// Drain blocks until all in-flight requests complete or timeout.
// Returns the number of connections that were still active at timeout.
func (d *Drainer) Drain() int {
	d.draining.Store(true)
	fmt.Println("drain: stopping new requests...")

	deadline := time.After(d.timeout)
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-deadline:
			remaining := d.checker.activeConns.Load()
			if remaining > 0 {
				fmt.Printf("drain: timeout with %d active connections\n", remaining)
			}
			return int(remaining)

		case <-tick.C:
			active := d.checker.activeConns.Load()
			if active == 0 {
				fmt.Println("drain: all connections drained")
				return 0
			}
			fmt.Printf("drain: waiting for %d connections...\n", active)
		}
	}
}
