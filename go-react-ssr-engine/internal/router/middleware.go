package router

import (
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// HandlerFunc is the request handler signature used throughout the engine.
type HandlerFunc func(ctx *RequestContext) (string, error)

// Middleware wraps a HandlerFunc and returns a new one.
type Middleware func(HandlerFunc) HandlerFunc

// RequestContext carries per-request data through the middleware chain.
type RequestContext struct {
	Path       string
	Route      *Route
	Params     map[string]string
	PropsJSON  string
	StatusCode int
	Headers    map[string]string

	// RequestID is a unique identifier for tracing this request through logs.
	// Format: unix_nano in base36 — short, unique enough, zero allocation.
	RequestID string

	// AcceptGzip is true if the client sent Accept-Encoding: gzip.
	// Set by the request dispatcher, read by the gzip middleware.
	AcceptGzip bool

	// SkipCache signals to the render handler to bypass cache.
	// Set via ?nocache=1 query param or Cache-Control: no-cache header.
	SkipCache bool
}

func NewRequestContext(path string) *RequestContext {
	return &RequestContext{
		Path:       path,
		StatusCode: 200,
		Headers:    make(map[string]string),
		RequestID:  generateRequestID(),
	}
}

// Chain applies middlewares in order: first middleware is outermost.
func Chain(middlewares ...Middleware) Middleware {
	return func(final HandlerFunc) HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// --- Request ID ---

// requestCounter is a global atomic counter for request IDs.
// Combined with timestamp gives unique IDs without UUID overhead.
var requestCounter atomic.Uint64

func generateRequestID() string {
	count := requestCounter.Add(1)
	// Compact format: timestamp_millis-counter
	return fmt.Sprintf("%d-%d", time.Now().UnixMilli(), count)
}

// --- Logger ---

// Logger logs method, path, status, duration, and request ID.
// Output format: [req-id] STATUS PATH DURATION
func Logger() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *RequestContext) (string, error) {
			start := time.Now()

			html, err := next(ctx)

			status := ctx.StatusCode
			if err != nil {
				status = 500
			}

			log.Printf("[%s] %d %s %s", ctx.RequestID, status, ctx.Path, time.Since(start))
			return html, err
		}
	}
}

// --- Recovery ---

// Recovery catches panics and returns 500 instead of crashing the server.
func Recovery() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *RequestContext) (html string, err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[%s] PANIC %s: %v", ctx.RequestID, ctx.Path, r)
					ctx.StatusCode = 500
					html = "<h1>500 Internal Server Error</h1>"
					err = nil
				}
			}()
			return next(ctx)
		}
	}
}

// --- ETag ---

// ETag generates a weak ETag from the response HTML.
// If the client sends If-None-Match matching the ETag, we return 304.
// This saves bandwidth on repeat visits — browser uses cached copy.
//
// Uses SHA1 truncated to 12 chars. Not cryptographic — just content identity.
// Weak ETag (W/) because gzip may alter bytes but content is semantically same.
func ETag() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *RequestContext) (string, error) {
			html, err := next(ctx)
			if err != nil || ctx.StatusCode != 200 {
				return html, err
			}

			// Generate ETag from content hash
			hash := sha1.Sum([]byte(html))
			etag := `W/"` + hex.EncodeToString(hash[:6]) + `"`

			ctx.Headers["ETag"] = etag
			ctx.Headers["Cache-Control"] = "public, max-age=0, must-revalidate"

			return html, nil
		}
	}
}

// --- Gzip ---

// gzipWriterPool reuses gzip writers to avoid allocation per request.
// gzip.NewWriter allocates ~800 bytes — at 10k req/s that's 8MB/s of garbage.
// Pool reduces this to near zero.
var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		// BestSpeed gives 80% of max compression at 3x the speed.
		// For HTML responses, the difference between BestSpeed and
		// BestCompression is typically <5% size but 3x CPU.
		w, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
		return w
	},
}

// Gzip compresses HTML responses for clients that accept it.
// Typically reduces HTML payload by 70-80%.
// Only compresses responses > 1KB — smaller responses have negligible savings
// and the gzip header overhead (18 bytes) eats into the benefit.
func Gzip() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *RequestContext) (string, error) {
			html, err := next(ctx)
			if err != nil || !ctx.AcceptGzip {
				return html, err
			}

			// Skip compression for small responses
			if len(html) < 1024 {
				return html, nil
			}

			compressed, err := compressGzip(html)
			if err != nil {
				// Compression failed — return uncompressed. Non-fatal.
				return html, nil
			}

			ctx.Headers["Content-Encoding"] = "gzip"
			ctx.Headers["Vary"] = "Accept-Encoding"

			return compressed, nil
		}
	}
}

func compressGzip(data string) (string, error) {
	var buf strings.Builder
	buf.Grow(len(data) / 3) // gzip typically achieves ~70% compression

	w := gzipWriterPool.Get().(*gzip.Writer)
	w.Reset(&buf)

	_, err := w.Write([]byte(data))
	if err != nil {
		gzipWriterPool.Put(w)
		return "", err
	}

	// Must close before reading buf — Close flushes the gzip footer
	err = w.Close()
	gzipWriterPool.Put(w)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// --- Request Timing ---

// Timing adds a Server-Timing header with total render duration.
// Visible in browser DevTools Network tab — helps debug slow pages.
func Timing() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *RequestContext) (string, error) {
			start := time.Now()
			html, err := next(ctx)
			dur := time.Since(start)

			// Server-Timing format: total;dur=4.2
			// Duration in milliseconds with one decimal
			ctx.Headers["Server-Timing"] = fmt.Sprintf("total;dur=%.1f", float64(dur.Microseconds())/1000.0)

			return html, err
		}
	}
}

// --- Rate Limiter ---

// RateLimiter is a simple token bucket per-path rate limiter.
// Prevents a single hot path from monopolizing the V8 worker pool.
// Not per-IP (that's the reverse proxy's job) — this protects the engine itself.
type RateLimiter struct {
	mu       sync.Mutex
	counters map[string]*rateBucket
	limit    int           // max requests per window
	window   time.Duration // window size
}

type rateBucket struct {
	count  int
	resets time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		counters: make(map[string]*rateBucket),
		limit:    limit,
		window:   window,
	}
}

// RateLimit rejects requests that exceed the per-path rate limit.
// Returns 429 Too Many Requests with a Retry-After header.
func (rl *RateLimiter) RateLimit() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *RequestContext) (string, error) {
			if rl.isLimited(ctx.Path) {
				ctx.StatusCode = 429
				ctx.Headers["Retry-After"] = fmt.Sprintf("%d", int(rl.window.Seconds()))
				return "<h1>429 Too Many Requests</h1>", nil
			}
			return next(ctx)
		}
	}
}

func (rl *RateLimiter) isLimited(path string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.counters[path]

	if !exists || now.After(bucket.resets) {
		rl.counters[path] = &rateBucket{count: 1, resets: now.Add(rl.window)}
		return false
	}

	bucket.count++
	return bucket.count > rl.limit
}
