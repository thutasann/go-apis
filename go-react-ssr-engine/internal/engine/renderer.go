package engine

import (
	"fmt"

	v8 "rogchap.com/v8go"
)

// Worker wraps a single V8 isolate + context.
// Each worker is single-threaded (V8 requirement) but multiple workers
// run on different OS threads concurrently via the pool.
//
// V8 isolate = isolated heap. No shared state between workers.
// This means no locks inside the render path - pure parallel execution.
type Worker struct {
	id  int
	iso *v8.Isolate
	ctx *v8.Context
}

// NewWorker creates one V8 isolate with a fresh context.
// The isolate is heavyweight (~10MB) so we reuse it across requets.
func NewWorker(id int) (*Worker, error) {
	iso := v8.NewIsolate()

	// Global object is where we inject the render bridge.
	// React's renderToString will be called via a global function.
	global := v8.NewObjectTemplate(iso)

	ctx := v8.NewContext(iso, global)

	return &Worker{
		id:  id,
		iso: iso,
		ctx: ctx,
	}, nil
}

// Execute runs the server bundle in this isolate and calls the render function.
// bundle = compiled JS from esbuild (contains React + all page components)
// route = matched page path like "/posts/123"
// props = JSON string of server-side props
//
// The bundle must expose a global __renderToString(route, props) function
// that returns a HTML string. We'll set this up in the bundler phase.
func (w *Worker) Execute(bundle, route, propsJSON string) (string, error) {
	// Run the bundle to define all functions adnd components
	// On subsequent calls with the same bundle, V8's code cache
	// makes this near-instant - only first run is expensive.
	_, err := w.ctx.RunScript(bundle, "server_bundle.js")
	if err != nil {
		return "", fmt.Errorf("worker %d: bundle exec failed: %w", w.id, err)
	}

	// Call the render bridge. This invokes ReactDOMServer.renderToString
	// inside V8 and returns the HTML string.
	renderCall := fmt.Sprintf(`__renderToString(%q, %s)`, route, propsJSON)

	val, err := w.ctx.RunScript(renderCall, "render.js")
	if err != nil {
		return "", fmt.Errorf("worker %d: render failed for %s: %w", w.id, route, err)
	}

	return val.String(), nil
}

// Dispose releases the V8 isolate memory.
// Called only during shutdown. After this, the worker is dead.
func (w *Worker) Dispose() {
	if w.ctx != nil {
		w.ctx.Close()
	}
	if w.iso != nil {
		w.iso.Dispose()
	}
}
