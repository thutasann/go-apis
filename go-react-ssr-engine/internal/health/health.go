package health

import (
	"encoding/json"
	"runtime"
	"sync/atomic"
	"time"
)

// Checker tracks server health metrics.
// Embedded in the HTTP handler, exposed at /_health.
// All counters are atomic — no locks in the request path.
type Checker struct {
	startTime     time.Time
	totalRequests atomic.Uint64
	totalErrors   atomic.Uint64
	activeConns   atomic.Int64

	// ready is false during startup (bundling, pool init).
	// Flipped to true once the server is accepting requests.
	// Used by load balancers: /_health returns 503 until ready.
	ready atomic.Bool
}

func NewChecker() *Checker {
	return &Checker{
		startTime: time.Now(),
	}
}

// MarkReady signals the server is fully initialized and accepting traffic.
func (c *Checker) MarkReady() {
	c.ready.Store(true)
}

// RecordRequest increments the total request counter and active connections.
// Call at the start of each request.
func (c *Checker) RecordRequest() {
	c.totalRequests.Add(1)
	c.activeConns.Add(1)
}

// RecordComplete decrements active connections.
// Call at the end of each request (in defer).
func (c *Checker) RecordComplete() {
	c.activeConns.Add(-1)
}

// RecordError increments the error counter.
func (c *Checker) RecordError() {
	c.totalErrors.Add(1)
}

// Status is the JSON response for /_health.
type Status struct {
	Status        string  `json:"status"` // "ok" or "starting"
	Uptime        string  `json:"uptime"`
	UptimeSeconds float64 `json:"uptime_seconds"`
	TotalRequests uint64  `json:"total_requests"`
	TotalErrors   uint64  `json:"total_errors"`
	ActiveConns   int64   `json:"active_connections"`
	ErrorRate     float64 `json:"error_rate_pct"` // errors / total * 100
	GoRoutines    int     `json:"goroutines"`
	MemAllocMB    float64 `json:"mem_alloc_mb"`
	MemSysMB      float64 `json:"mem_sys_mb"`
	NumGC         uint32  `json:"num_gc"`
}

// Check returns current health status as JSON bytes.
func (c *Checker) Check() ([]byte, int) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	total := c.totalRequests.Load()
	errors := c.totalErrors.Load()
	uptime := time.Since(c.startTime)

	var errorRate float64
	if total > 0 {
		errorRate = float64(errors) / float64(total) * 100
	}

	status := "ok"
	httpCode := 200
	if !c.ready.Load() {
		status = "starting"
		httpCode = 503
	}

	s := Status{
		Status:        status,
		Uptime:        uptime.Round(time.Second).String(),
		UptimeSeconds: uptime.Seconds(),
		TotalRequests: total,
		TotalErrors:   errors,
		ActiveConns:   c.activeConns.Load(),
		ErrorRate:     errorRate,
		GoRoutines:    runtime.NumGoroutine(),
		MemAllocMB:    float64(mem.Alloc) / 1024 / 1024,
		MemSysMB:      float64(mem.Sys) / 1024 / 1024,
		NumGC:         mem.NumGC,
	}

	data, _ := json.MarshalIndent(s, "", "  ")
	return data, httpCode
}
