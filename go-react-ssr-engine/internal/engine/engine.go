package engine

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/thutasann/go-react-ssr-engine/internal/config"
)

// Engine coordinates V8 rendering across the worker pool.
// Thread-safe for concurrent Render/RenderProps calls.
type Engine struct {
	pool *Pool
	cfg  *config.Config

	mu           sync.RWMutex
	serverBundle string
}

func New(cfg *config.Config) (*Engine, error) {
	e := &Engine{cfg: cfg}

	pool, err := NewPool(cfg.WorkerPoolSize)
	if err != nil {
		return nil, fmt.Errorf("engine: pool init: %w", err)
	}
	e.pool = pool

	return e, nil
}

// LoadBundle swaps the server JS bundle atomically.
// In-flight renders complete with old bundle. New renders get new bundle.
func (e *Engine) LoadBundle(js string) {
	e.mu.Lock()
	e.serverBundle = js
	e.mu.Unlock()
}

// Render executes React SSR for a route with given props JSON.
// Acquires a worker, renders, releases. Blocks if pool exhausted.
func (e *Engine) Render(route string, propsJSON string) (string, error) {
	e.mu.RLock()
	bundle := e.serverBundle
	e.mu.RUnlock()

	if bundle == "" {
		return "", fmt.Errorf("engine: no bundle loaded")
	}

	worker := e.pool.Acquire()
	defer e.pool.Release(worker)

	return worker.Execute(bundle, route, propsJSON)
}

// RenderProps executes getServerSideProps for a route.
// context is serialized PageContext JSON.
// Returns raw JSON string from V8.
func (e *Engine) RenderProps(route string, context interface{}) (string, error) {
	e.mu.RLock()
	bundle := e.serverBundle
	e.mu.RUnlock()

	if bundle == "" {
		return "", fmt.Errorf("engine: no bundle loaded")
	}

	ctxJSON, err := json.Marshal(context)
	if err != nil {
		return "", fmt.Errorf("engine: context marshal: %w", err)
	}

	worker := e.pool.Acquire()
	defer e.pool.Release(worker)

	return worker.ExecuteProps(bundle, route, string(ctxJSON))
}

func (e *Engine) Shutdown() {
	e.pool.Shutdown()
}
