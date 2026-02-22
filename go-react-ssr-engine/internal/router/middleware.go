package router

import (
	"log"
	"time"
)

// HandlerFunc is the request handler signature used throughout the engine.
// Receives a context with route info, return HTML or error
type HandlerFunc func(ctx *RequestContext) (string, error)

// Middleware wraps a HandlerFunc and returns a new one.
// Classic onion model - outermost middleware runs first.
type Middleware func(HandlerFunc) HandlerFunc

// RequestContext carries per-request data through the middleware chain.
// Created once per request, never shared across goroutines.
type RequestContext struct {
	// Path is the raw URL path, e.g. "/posts/123"
	Path string

	// Route is the matched route, nil if 404
	Route *Route

	// Params holds extracted URL params, e.g. {"id": "123"}
	Params map[string]string

	// PropsJSON is serialized props from getServerSideProps.
	// Populated by the props loader middleware.
	PropsJSON string

	// StatusCode defaults to 200, middleware can change it.
	StatusCode int

	// Headers from the HTTP response
	Headers map[string]string
}

// NewRequestContext creates a context with sensible defaults.
func NewRequestContext(path string) *RequestContext {
	return &RequestContext{
		Path:       path,
		StatusCode: 200,
		Headers:    make(map[string]string),
	}
}

// Chain applies middlewares in order: first middleware is outermost.
// Chain(A, B, C)(handler) executes as: A -> B -> C -> handler -> C -> B -> A
func Chain(middlewares ...Middleware) Middleware {
	return func(final HandlerFunc) HandlerFunc {
		// Apply in reverse so first middleware wraps outermost
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// Logger middleware logs request duration and status.
// Runs on every request. Overhead: one time.Now() call + one log.PrintF.
func Logger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *RequestContext) (string, error) {
			start := time.Now()

			html, err := next(ctx)

			status := ctx.StatusCode
			if err != nil {
				status = 500
			}

			log.Printf("%d %s %s", status, ctx.Path, time.Since(start))
			return html, err
		}
	}
}

// Recovery middleware catches panics from handlers and converts to 500.
// Without this, a panic in one render kills the entire server.
func Recovery() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *RequestContext) (html string, err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic recovered on %s: %v", ctx.Path, r)
					ctx.StatusCode = 500
					html = "<h1>500 Internal Server Error</h1>"
					err = nil // swallow panic, return error page
				}
			}()
			return next(ctx)
		}
	}
}
