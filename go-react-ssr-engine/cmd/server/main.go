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
	handler := buildHandler(cfg, eng, rt, lru, hydrator)

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
	}

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		fmt.Printf("\n  reactgo ready at http://localhost%s\n\n", addr)
		if err := server.ListenAndServe(addr); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nshutting down...")
	server.Shutdown()
	eng.Shutdown()
	fmt.Println("done")
}

func buildHandler(
	cfg *config.Config,
	eng *engine.Engine,
	rt *router.Router,
	lru *cache.LRU,
	hydrator *hydration.Hydrator,
) fasthttp.RequestHandler {

	// Static files from public/
	publicFS := &fasthttp.FS{
		Root:               cfg.PublicDir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: false,
		AcceptByteRange:    true,
	}
	publicHandler := publicFS.NewRequestHandler()

	// Client bundles at /_reactgo/*
	clientDir := filepath.Join(cfg.BuildDir, "client")
	clientFS := &fasthttp.FS{
		Root:               clientDir,
		IndexNames:         nil,
		GenerateIndexPages: false,
		AcceptByteRange:    true,
	}
	clientHandler := clientFS.NewRequestHandler()

	chain := router.Chain(
		router.Recovery(),
		router.Logger(),
	)

	renderFn := func(rctx *router.RequestContext) (string, error) {
		// --- Cache check ---
		pageCtx := &props.PageContext{
			Route:  rctx.Route.Pattern,
			Params: rctx.Params,
			Path:   rctx.Path,
		}
		cacheKey := props.CacheKey(pageCtx)

		if cached, ok := lru.Get(cacheKey); ok {
			return cached, nil
		}

		// --- Props loading via V8 ---
		propsJSON := "{}"
		propsResult, err := eng.RenderProps(rctx.Route.Pattern, pageCtx)
		if err != nil {
			log.Printf("props error: %v", err)
			// Non-fatal — render with empty props
		} else {
			// Parse to check for redirect/notFound
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
			// Extract just the props field for rendering
			var wrapper struct {
				Props json.RawMessage `json:"props"`
			}
			if err := json.Unmarshal([]byte(propsResult), &wrapper); err == nil && wrapper.Props != nil {
				propsJSON = string(wrapper.Props)
			}
		}

		// --- SSR Render ---
		bodyHTML, err := eng.Render(rctx.Route.Pattern, propsJSON)
		if err != nil {
			return "", err
		}

		// --- Assemble document ---
		doc := html.NewDocument()
		doc.BodyHTML = bodyHTML
		hydrator.Prepare(doc, rctx.Route.Pattern, propsJSON)

		fullHTML := doc.Render()
		lru.Set(cacheKey, fullHTML)

		return fullHTML, nil
	}

	handlerWithMiddleware := chain(renderFn)

	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		// Client bundles — must check before public to avoid collision
		if strings.HasPrefix(path, "/_reactgo/") {
			// Rewrite path: strip prefix for FS handler
			ctx.URI().SetPath(strings.TrimPrefix(path, "/_reactgo"))
			clientHandler(ctx)
			return
		}

		// Static public files
		if fileExists(filepath.Join(cfg.PublicDir, path)) {
			publicHandler(ctx)
			return
		}

		// Route matching
		route, params, found := rt.Match(path)
		if !found {
			ctx.SetStatusCode(404)
			ctx.SetContentType("text/html; charset=utf-8")
			fmt.Fprintf(ctx, "<h1>404 - Not Found</h1><p>No page matches %s</p>", path)
			return
		}

		// Build request context
		rctx := router.NewRequestContext(path)
		rctx.Route = route
		rctx.Params = params

		// Parse query string into params
		ctx.QueryArgs().VisitAll(func(key, value []byte) {
			if rctx.Params == nil {
				rctx.Params = make(map[string]string)
			}
			rctx.Params["query_"+string(key)] = string(value)
		})

		// Render
		htmlResult, err := handlerWithMiddleware(rctx)
		if err != nil {
			ctx.SetStatusCode(500)
			ctx.SetContentType("text/html; charset=utf-8")
			fmt.Fprintf(ctx, "<h1>500 - Server Error</h1>")
			if cfg.Dev {
				fmt.Fprintf(ctx, "<pre>%v</pre>", err)
			}
			log.Printf("render error: %v", err)
			return
		}

		// Redirects
		if location, ok := rctx.Headers["Location"]; ok {
			ctx.Redirect(location, rctx.StatusCode)
			return
		}

		ctx.SetStatusCode(rctx.StatusCode)
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.WriteString(htmlResult)
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
