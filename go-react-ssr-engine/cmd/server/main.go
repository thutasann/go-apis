package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/thutasann/go-react-ssr-engine/internal/bundler"
	"github.com/thutasann/go-react-ssr-engine/internal/cache"
	"github.com/thutasann/go-react-ssr-engine/internal/config"
	"github.com/thutasann/go-react-ssr-engine/internal/engine"
	"github.com/thutasann/go-react-ssr-engine/internal/hydration"
	"github.com/thutasann/go-react-ssr-engine/internal/props"
	"github.com/thutasann/go-react-ssr-engine/internal/router"
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
