package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/thutasann/go-react-ssr-engine/internal/bundler"
	"github.com/thutasann/go-react-ssr-engine/internal/cache"
	"github.com/thutasann/go-react-ssr-engine/internal/config"
	"github.com/thutasann/go-react-ssr-engine/internal/engine"
	"github.com/thutasann/go-react-ssr-engine/internal/health"
	"github.com/thutasann/go-react-ssr-engine/internal/hydration"
	"github.com/thutasann/go-react-ssr-engine/internal/props"
	"github.com/thutasann/go-react-ssr-engine/internal/router"
	"github.com/thutasann/go-react-ssr-engine/pkg/html"
	"github.com/valyala/fasthttp"
)

func main() {
	dev := flag.Bool("dev", false, "enable dev mode")
	port := flag.Int("port", 0, "override port")
	flag.Parse()

	cfg, err := config.Load("reactgo.config.json")
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if *dev {
		cfg.Dev = true
	}
	if *port != 0 {
		cfg.Port = *port
	}

	// --- Health checker ---
	checker := health.NewChecker()
	drainer := health.NewDrainer(checker, 30*time.Second)

	// --- 1. Bundle ---
	b := bundler.New(cfg)
	buildResult, err := b.Build()
	if err != nil {
		log.Fatalf("bundler: %v", err)
	}
	fmt.Printf("bundler: server bundle %d bytes, %d client entries\n",
		len(buildResult.ServerBundle), len(buildResult.ClientEntries))

	// --- 2. V8 Engine ---
	eng, err := engine.New(cfg)
	if err != nil {
		log.Fatalf("engine: %v", err)
	}
	eng.LoadBundle(buildResult.ServerBundle)

	// --- 3. Router ---
	rt, err := router.New(cfg)
	if err != nil {
		log.Fatalf("router: %v", err)
	}
	for _, route := range rt.Routes() {
		fmt.Printf("  route: %s\n", route)
	}

	// --- 4. Cache ---
	lru := cache.NewLRU(cfg.CacheMaxEntries)
	if cfg.Dev {
		lru = cache.NewLRU(0)
	}

	// --- 5. Hydration ---
	manifest := hydration.NewManifest()
	clientDir := filepath.Join(cfg.BuildDir, "client")
	manifest.Build(clientDir, buildResult.ClientEntries)
	hydrator := hydration.NewHydrator(manifest)

	// --- 6. Handler ---
	handler := buildHandler(cfg, eng, rt, lru, hydrator, checker, drainer)

	// --- 7. Watcher ---
	if cfg.Dev {
		inv := cache.NewInvalidator(lru, cfg.PagesDir)
		w := bundler.NewWatcher(cfg, b, func(result *bundler.BuildResult) {
			eng.LoadBundle(result.ServerBundle)
			rt.Rebuild()
			manifest.Build(clientDir, result.ClientEntries)
			inv.OnRebuild()
			fmt.Println("hot reload complete")
		})
		if err := w.Start(); err != nil {
			log.Printf("watcher: %v", err)
		}
	}

	// --- 8. Server ---
	server := &fasthttp.Server{
		Handler:                       handler,
		Name:                          "reactgo",
		Concurrency:                   256 * 1024,
		DisableHeaderNamesNormalizing: true,
		ReadTimeout:                   10 * time.Second,
		WriteTimeout:                  15 * time.Second,
		IdleTimeout:                   120 * time.Second,
		MaxRequestBodySize:            4 * 1024 * 1024, // 4MB
	}

	// Mark ready after everything is initialized
	checker.MarkReady()

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		fmt.Printf("\n  reactgo ready at http://localhost%s\n", addr)
		fmt.Printf("  health: http://localhost%s/_health\n\n", addr)
		if err := server.ListenAndServe(addr); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	// --- 9. Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nshutting down...")

	// Stop accepting new connections
	server.Shutdown()

	// Wait for in-flight requests to finish
	remaining := drainer.Drain()
	if remaining > 0 {
		fmt.Printf("warning: %d requests did not complete\n", remaining)
	}

	eng.Shutdown()
	fmt.Println("done")
}

func buildHandler(
	cfg *config.Config,
	eng *engine.Engine,
	rt *router.Router,
	lru *cache.LRU,
	hydrator *hydration.Hydrator,
	checker *health.Checker,
	drainer *health.Drainer,
) fasthttp.RequestHandler {

	// Static file handlers
	publicFS := &fasthttp.FS{
		Root:               cfg.PublicDir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: false,
		AcceptByteRange:    true,
		Compress:           true, // fasthttp built-in compression for static files
	}
	publicHandler := publicFS.NewRequestHandler()

	clientDir := filepath.Join(cfg.BuildDir, "client")
	clientFS := &fasthttp.FS{
		Root:               clientDir,
		IndexNames:         nil,
		GenerateIndexPages: false,
		AcceptByteRange:    true,
		Compress:           true,
		// Client bundles are content-hashed — safe to cache forever.
		// Browser will fetch new URL when hash changes after rebuild.
		CacheDuration: 365 * 24 * time.Hour,
	}
	clientHandler := clientFS.NewRequestHandler()

	// Rate limiter: 1000 requests per path per second.
	// Protects V8 pool from being monopolized by a single hot route.
	limiter := router.NewRateLimiter(1000, 1*time.Second)

	// Middleware chain — order matters:
	// Recovery (outermost, catches everything)
	// -> Timing (measures total including middleware)
	// -> RateLimit (reject before expensive work)
	// -> Logger (log after we know the status)
	// -> ETag (add caching headers)
	// -> Gzip (compress last, after ETag is computed)
	chain := router.Chain(
		router.Recovery(),
		router.Timing(),
		limiter.RateLimit(),
		router.Logger(),
		router.ETag(),
		router.Gzip(),
	)

	renderFn := func(rctx *router.RequestContext) (string, error) {
		// --- Cache ---
		pageCtx := &props.PageContext{
			Route:  rctx.Route.Pattern,
			Params: rctx.Params,
			Path:   rctx.Path,
		}
		cacheKey := props.CacheKey(pageCtx)

		if !rctx.SkipCache {
			if cached, ok := lru.Get(cacheKey); ok {
				return cached, nil
			}
		}

		// --- Props ---
		propsJSON := "{}"
		propsResult, err := eng.RenderProps(rctx.Route.Pattern, pageCtx)
		if err != nil {
			log.Printf("[%s] props error: %v", rctx.RequestID, err)
		} else {
			var pageProps props.PageProps
			if err := json.Unmarshal([]byte(propsResult), &pageProps); err == nil {
				if pageProps.Redirect != nil {
					rctx.StatusCode = 302
					if pageProps.Redirect.Permanent {
						rctx.StatusCode = 301
					}
					rctx.Headers["Location"] = pageProps.Redirect.Destination
					return "", nil
				}
				if pageProps.NotFound {
					rctx.StatusCode = 404
					return "<h1>404 - Not Found</h1>", nil
				}
			}
			propsJSON = propsResult
			var wrapper struct {
				Props json.RawMessage `json:"props"`
			}
			if err := json.Unmarshal([]byte(propsResult), &wrapper); err == nil && wrapper.Props != nil {
				propsJSON = string(wrapper.Props)
			}
		}

		// --- SSR ---
		bodyHTML, err := eng.Render(rctx.Route.Pattern, propsJSON)
		if err != nil {
			return "", err
		}

		// --- Document ---
		doc := html.NewDocument()
		doc.BodyHTML = bodyHTML
		hydrator.Prepare(doc, rctx.Route.Pattern, propsJSON)

		fullHTML := doc.Render()

		if !rctx.SkipCache {
			lru.Set(cacheKey, fullHTML)
		}

		return fullHTML, nil
	}

	handlerWithMiddleware := chain(renderFn)

	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		// --- Health endpoint ---
		if path == "/_health" {
			data, code := checker.Check()
			ctx.SetStatusCode(code)
			ctx.SetContentType("application/json")
			ctx.Write(data)
			return
		}

		// --- Draining check ---
		if drainer.IsDraining() {
			ctx.SetStatusCode(503)
			ctx.SetContentType("text/plain")
			ctx.WriteString("service shutting down")
			return
		}

		// --- Track request ---
		checker.RecordRequest()
		defer checker.RecordComplete()

		// --- Client bundles ---
		if strings.HasPrefix(path, "/_reactgo/") {
			ctx.URI().SetPath(strings.TrimPrefix(path, "/_reactgo"))
			clientHandler(ctx)
			return
		}

		// --- Static files ---
		if fileExists(filepath.Join(cfg.PublicDir, path)) {
			publicHandler(ctx)
			return
		}

		// --- Route matching ---
		route, params, found := rt.Match(path)
		if !found {
			ctx.SetStatusCode(404)
			ctx.SetContentType("text/html; charset=utf-8")
			fmt.Fprintf(ctx, "<h1>404 - Not Found</h1>")
			return
		}

		// --- Build request context ---
		rctx := router.NewRequestContext(path)
		rctx.Route = route
		rctx.Params = params
		rctx.AcceptGzip = strings.Contains(
			string(ctx.Request.Header.Peek("Accept-Encoding")), "gzip",
		)
		rctx.SkipCache = len(ctx.QueryArgs().Peek("nocache")) > 0 ||
			strings.Contains(string(ctx.Request.Header.Peek("Cache-Control")), "no-cache")

		// Parse query params
		ctx.QueryArgs().VisitAll(func(key, value []byte) {
			if rctx.Params == nil {
				rctx.Params = make(map[string]string)
			}
			rctx.Params["query_"+string(key)] = string(value)
		})

		// --- Render ---
		htmlResult, err := handlerWithMiddleware(rctx)
		if err != nil {
			checker.RecordError()
			ctx.SetStatusCode(500)
			ctx.SetContentType("text/html; charset=utf-8")
			fmt.Fprintf(ctx, "<h1>500 - Server Error</h1>")
			if cfg.Dev {
				fmt.Fprintf(ctx, "<pre>%v</pre>", err)
			}
			log.Printf("[%s] render error: %v", rctx.RequestID, err)
			return
		}

		// --- Redirects ---
		if location, ok := rctx.Headers["Location"]; ok {
			ctx.Redirect(location, rctx.StatusCode)
			return
		}

		// --- ETag 304 check ---
		if etag, ok := rctx.Headers["ETag"]; ok {
			clientETag := string(ctx.Request.Header.Peek("If-None-Match"))
			if clientETag == etag {
				ctx.SetStatusCode(304)
				return
			}
		}

		// --- Write response ---
		ctx.SetStatusCode(rctx.StatusCode)
		ctx.SetContentType("text/html; charset=utf-8")

		// Set all headers from middleware chain
		for k, v := range rctx.Headers {
			ctx.Response.Header.Set(k, v)
		}

		// Add security headers
		ctx.Response.Header.Set("X-Content-Type-Options", "nosniff")
		ctx.Response.Header.Set("X-Frame-Options", "SAMEORIGIN")
		ctx.Response.Header.Set("X-Request-ID", rctx.RequestID)

		ctx.WriteString(htmlResult)
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
