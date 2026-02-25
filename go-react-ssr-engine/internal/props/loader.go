package props

import (
	"encoding/json"
	"fmt"
	"sync"

	v8 "rogchap.com/v8go"
)

// Loader executes getServerSideProps functions defined in page bundles.
// Each page can optionally export a getServerSideProps function.
// If present, it runs server-side before render to fetch data.
//
// The loader runs in the same V8 isolate pool as the renderer.
// This keeps props fetching fast and avoids extra process overhead.
type Loader struct {
	// registry maps route pattern -> whether the page has getServerSideProps.
	// Built during bundler analysis. Checked before attempting to call
	// the function in V8 — avoids wasted V8 calls for static pages.
	registry map[string]bool
	mu       sync.RWMutex
}

func NewLoader() *Loader {
	return &Loader{
		registry: make(map[string]bool),
	}
}

// Register marks a route as having getServerSideProps.
// Called during bundle analysis.
func (l *Loader) Register(route string) {
	l.mu.Lock()
	l.registry[route] = true
	l.mu.Unlock()
}

// HasServerProps checks if a route needs server-side props loading.
// Fast path: if false, skip V8 entirely and render with empty props.
func (l *Loader) HasServerProps(route string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.registry[route]
}

// ClearRegistry wipes registration on rebuild.
func (l *Loader) ClearRegistry() {
	l.mu.Lock()
	l.registry = make(map[string]bool)
	l.mu.Unlock()
}

// LoadProps executes getServerSideProps in a V8 context.
// The bundle must expose __getServerSideProps(route, context) globally.
// Returns PageProps which may contain data, redirect, or notFound.
//
// ctx contains route params, query string, etc.
// bundle is the server JS (same one used for rendering).
func LoadProps(ctx *PageContext, iso *v8.Isolate, v8ctx *v8.Context, bundle string) (*PageProps, error) {
	// Serialize the page context to JSON so V8 can read it
	ctxJSON, err := json.Marshal(ctx)
	if err != nil {
		return nil, fmt.Errorf("props: context serialize failed: %w", err)
	}

	// Call the props bridge function in V8.
	// This runs the page's getServerSideProps and returns JSON.
	script := fmt.Sprintf(`__getServerSideProps(%q, %s)`, ctx.Route, string(ctxJSON))

	val, err := v8ctx.RunScript(script, "props.js")
	if err != nil {
		return nil, fmt.Errorf("props: execution failed for %s: %w", ctx.Route, err)
	}

	// Parse the returned JSON into PageProps
	var props PageProps
	if err := json.Unmarshal([]byte(val.String()), &props); err != nil {
		return nil, fmt.Errorf("props: parse failed for %s: %w", ctx.Route, err)
	}

	return &props, nil
}

// CacheKey generates a unique cache key from route + params.
// Two requests to /posts/:id with different ids get different cache keys.
// Format: "/posts/:id|id=123" — simple, deterministic, no hash collisions.
func CacheKey(ctx *PageContext) string {
	key := ctx.Route
	if len(ctx.Params) > 0 {
		key += "|"
		for k, v := range ctx.Params {
			key += k + "=" + v + "&"
		}
	}
	return key
}
