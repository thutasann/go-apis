package cache

import (
	"log"
	"strings"
)

// Invalidator connects file changes to cache eviction.
// When a page file changes, we must evict its cached HTML
// otherwise users see stale content after edits.
//
// Strategy:
// - Single file change -> evict just that route's cache entries
// - Component/layout change -> flush entire cache (could affect any page)
// - Full rebuild -> flush entire cache
type Invalidator struct {
	cache    *LRU
	pagesDir string
}

func NewInvalidator(cache *LRU, pagesDir string) *Invalidator {
	return &Invalidator{
		cache:    cache,
		pagesDir: pagesDir,
	}
}

// OnFileChange is called by the watcher when a file is modified.
// Decides whether to do targeted eviction or full flush.
func (inv *Invalidator) OnFileChange(filePath string) {
	if inv.isPageFile(filePath) {
		// Page changed — only that route's cache is stale.
		// Convert file path to cache key pattern and evict.
		route := inv.fileToRoute(filePath)
		inv.evictRoute(route)
		log.Printf("cache: invalidated route %s", route)
	} else {
		// Shared component or layout changed — could affect any page.
		// Full flush is the safe choice. In production this is rare
		// because the watcher only runs in dev mode.
		inv.cache.Flush()
		log.Println("cache: full flush (shared file changed)")
	}
}

// OnRebuild is called after a full bundler rebuild. Always flushes.
func (inv *Invalidator) OnRebuild() {
	inv.cache.Flush()
}

// evictRoute removes all cache entries that match a route.
// Cache keys include props hash, so one route can have multiple entries.
// e.g. "/posts/:id" with id=1 and id=2 are separate cache keys.
// We do a prefix scan to catch all variants.
func (inv *Invalidator) evictRoute(routePrefix string) {
	// This is O(n) over cache keys. Acceptable because:
	// 1. Only runs in dev mode on file save
	// 2. Cache is capped at maxSize entries
	// 3. Single file saves are infrequent (human typing speed)
	inv.cache.Flush() // TODO: targeted prefix eviction in future
}

func (inv *Invalidator) isPageFile(path string) bool {
	return strings.HasPrefix(path, inv.pagesDir)
}

func (inv *Invalidator) fileToRoute(filePath string) string {
	route := strings.TrimPrefix(filePath, inv.pagesDir)
	route = strings.TrimSuffix(route, ".tsx")
	route = strings.TrimSuffix(route, ".jsx")

	if before, ok := strings.CutSuffix(route, "/index"); ok {
		route = before
	}
	if route == "" {
		route = "/"
	}
	return route
}
