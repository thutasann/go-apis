package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/thutasann/go-react-ssr-engine/internal/bundler"
	"github.com/thutasann/go-react-ssr-engine/internal/cache"
	"github.com/thutasann/go-react-ssr-engine/internal/config"
	"github.com/thutasann/go-react-ssr-engine/internal/engine"
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

	// --- 1. Bundle React pages ---
	b := bundler.New(cfg)
	buildResult, err := b.Build()
	if err != nil {
		log.Fatalf("bundler: %v", err)
	}

	// --- 2. Start V8 engine pool ---
	eng, err := engine.New(cfg)
	if err != nil {
		log.Fatalf("engine: %v", err)
	}
	eng.LoadBundle(buildResult.ServerBundle)

	// --- 3. Build route tree ---
	rt, err := router.New(cfg)
	if err != nil {
		log.Fatalf("router: %v", err)
	}
	for _, route := range rt.Routes() {
		fmt.Printf("  route: %s\n", route)
	}

	// --- 4. Setup cache ---
	lru := cache.NewLRU(cfg.CacheMaxEntries)
	if cfg.Dev {
		lru = cache.NewLRU(0) // no cache in dev — always fresh renders
	}

	// --- 5. Setup hydration manifest ---
	manifest := hydration.NewManifest()
	clientDir := filepath.Join(cfg.BuildDir, "client")
	manifest.Build(clientDir, buildResult.ClientEntries)
	hydrator := hydration.NewHydrator(manifest)

	// --- 6. Props loader ---
	propsLoader := props.NewLoader()

	// --- 7. Build request handler ---
	handler := buildHandler(cfg, eng, rt, lru, hydrator, propsLoader)

	// --- 8. File watcher for dev mode ---
	if cfg.Dev {
		inv := cache.NewInvalidator(lru, cfg.PagesDir)
		w := bundler.NewWatcher(cfg, b, func(result *bundler.BuildResult) {
			eng.LoadBundle(result.ServerBundle)
			rt.Rebuild()
			manifest.Build(clientDir, result.ClientEntries)
			inv.OnRebuild()
		})
		if err := w.Start(); err != nil {
			log.Printf("watcher: %v (continuing without hot reload)", err)
		}
	}

	// --- 9. Start server ---
	server := &fasthttp.Server{
		Handler:                       handler,
		Name:                          "reactgo",
		Concurrency:                   256 * 1024, // max concurrent connections
		DisableHeaderNamesNormalizing: true,       // skip unnecessary work
	}

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		fmt.Printf("reactgo listening on %s (dev=%v, workers=%d)\n", addr, cfg.Dev, cfg.WorkerPoolSize)
		if err := server.ListenAndServe(addr); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	// --- 10. Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("shutting down...")
	server.Shutdown()
	eng.Shutdown()
	fmt.Println("done")
}

// buildHandler creates the fasthttp request handler.
// This is the hot path — every HTTP request goes through here.
// No allocations except for cache misses (where we must render).
func buildHandler(
	cfg *config.Config,
	eng *engine.Engine,
	rt *router.Router,
	lru *cache.LRU,
	hydrator *hydration.Hydrator,
	propsLoader *props.Loader,
) fasthttp.RequestHandler {

	// Static file handler for public/ directory
	publicFS := &fasthttp.FS{
		Root:               cfg.PublicDir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: false,
		AcceptByteRange:    true, // supports Range requests for large files
	}
	publicHandler := publicFS.NewRequestHandler()

	// Static file handler for client bundles (/_reactgo/*)
	clientFS := &fasthttp.FS{
		Root:               filepath.Join(cfg.BuildDir, "client"),
		IndexNames:         nil,
		GenerateIndexPages: false,
		AcceptByteRange:    true,
		PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
			// Strip /_reactgo/ prefix so FS finds the right file
			path := string(ctx.Path())
			return []byte(strings.TrimPrefix(path, "/_reactgo"))
		},
	}
	clientHandler := clientFS.NewRequestHandler()

	// Middleware chain — applied to every SSR request
	chain := router.Chain(
		router.Recovery(),
		router.Logger(),
	)

	// The core render handler that the middleware wraps
	renderFn := func(rctx *router.RequestContext) (string, error) {
		// --- Cache check ---
		cacheKey := props.CacheKey(&props.PageContext{
			Route:  rctx.Route.Pattern,
			Params: rctx.Params,
		})

		if cached, ok := lru.Get(cacheKey); ok {
			return cached, nil
		}

		// --- Load props if page has getServerSideProps ---
		// pageCtx := &props.PageContext{
		// 	Route:  rctx.Route.Pattern,
		// 	Params: rctx.Params,
		// 	Path:   rctx.Path,
		// }

		var pageProps *props.PageProps
		if propsLoader.HasServerProps(rctx.Route.Pattern) {
			// TODO: props loading via V8 — wired in Phase 8
			pageProps = &props.PageProps{Props: make(map[string]interface{})}
		} else {
			pageProps = &props.PageProps{Props: make(map[string]interface{})}
		}

		// Handle redirects from props
		if pageProps.Redirect != nil {
			rctx.StatusCode = 302
			if pageProps.Redirect.Permanent {
				rctx.StatusCode = 301
			}
			rctx.Headers["Location"] = pageProps.Redirect.Destination
			return "", nil
		}

		// Handle 404 from props
		if pageProps.NotFound {
			rctx.StatusCode = 404
			return "<h1>404 - Not Found</h1>", nil
		}

		// --- SSR Render ---
		propsJSON, err := pageProps.ToJSON()
		if err != nil {
			return "", err
		}

		bodyHTML, err := eng.Render(rctx.Route.Pattern, propsJSON)
		if err != nil {
			return "", err
		}

		// --- Assemble full HTML document ---
		doc := html.NewDocument()
		doc.BodyHTML = bodyHTML
		hydrator.Prepare(doc, rctx.Route.Pattern, propsJSON)

		fullHTML := doc.Render()

		// --- Cache the result ---
		lru.Set(cacheKey, fullHTML)

		return fullHTML, nil
	}

	// Wrap render function with middleware
	handlerWithMiddleware := chain(renderFn)

	// Top-level request dispatcher
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		// Static files from public/
		if fileExists(filepath.Join(cfg.PublicDir, path)) {
			publicHandler(ctx)
			return
		}

		// Client bundles
		if strings.HasPrefix(path, "/_reactgo/") {
			clientHandler(ctx)
			return
		}

		// Route matching
		route, params, found := rt.Match(path)
		if !found {
			ctx.SetStatusCode(404)
			ctx.SetContentType("text/html; charset=utf-8")
			ctx.WriteString("<h1>404 - Not Found</h1>")
			return
		}

		// Build request context
		rctx := router.NewRequestContext(path)
		rctx.Route = route
		rctx.Params = params

		// Execute render pipeline
		htmlResult, err := handlerWithMiddleware(rctx)
		if err != nil {
			ctx.SetStatusCode(500)
			ctx.SetContentType("text/html; charset=utf-8")
			ctx.WriteString("<h1>500 - Internal Server Error</h1>")
			log.Printf("render error: %v", err)
			return
		}

		// Handle redirects
		if location, ok := rctx.Headers["Location"]; ok {
			ctx.Redirect(location, rctx.StatusCode)
			return
		}

		// Send HTML response
		ctx.SetStatusCode(rctx.StatusCode)
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.WriteString(htmlResult)
	}
}

// fileExists is a fast check used by the static file dispatcher.
// os.Stat is cached by the OS so repeated calls are cheap.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
