package engine

import (
	"fmt"
	"sync"

	"github.com/thutasann/go-react-ssr-engine/internal/config"
)

// Engine owns the V8 worker pool and coordinates rendering.
// One Engine per process. Created at startup, lives until shutdown.
type Engine struct {
	pool *Pool
	cfg  *config.Config

	// mu guards hot-reload bundle swaps. RLock for renders (many concurrent),
	// Lock for bundle replacement (rare, blocks briefly)
	mu sync.RWMutex

	// serverBundle holds the compiled JS that V8 executes.
	// Swapped atomically on rebuild. Every render reads this
	serverBundle string
}

// New creates the engine and spins up the V8 worker pool.
// Call shutdown() when done to release all V8 memory.
func New(cfg *config.Config) (*Engine, error) {
	e := &Engine{
		cfg: cfg,
	}

	pool, err := NewPool(cfg.WorkerPoolSize)
	if err != nil {
		return nil, fmt.Errorf("engine: pool init failed: %w", err)
	}
	e.pool = pool
	return e, nil
}

// LoadBundle sets the server-side JS bundle that workers execute.
// Called after esbuild compiles pages. Safe to call during live requests -
// in-flight renders finish with old bundle, new renders get the new one.
func (e *Engine) LoadBundle(js string) {
	e.mu.Lock()
	e.serverBundle = js
	e.mu.Unlock()
}

// Render takes a route path and JSON props, returns HTML string.
// Grabes a worker from the pool, executes React renderToString, returns worker.
// Blocks only if all workers are busy - backpressure is automatic.
func (e *Engine) Render(route string, propsJSON string) (string, error) {
	e.mu.RLock()
	bundle := e.serverBundle
	e.mu.RUnlock()

	if bundle == "" {
		return "", fmt.Errorf("engine: no bundle loaded")
	}

	// Acquire a worker - blocks if pool is exhausted.
	// This is the backpressure mechanism: under load, requets
	// queue here instead of spawning unbounded goroutines.
	worker := e.pool.Acquire()
	defer e.pool.Release(worker)

	html, err := worker.Execute(bundle, route, propsJSON)
	if err != nil {
		return "", fmt.Errorf("engine: render %s failed: %w", route, err)
	}

	return html, nil
}

// Shutdown drains the pool and destroys all V8 isolates.
// Call this on SIGTERM. After Shutdown, Render calls with panic.
func (e *Engine) Shutdown() {
	e.pool.Shutdown()
}
