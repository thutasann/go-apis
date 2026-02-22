package router

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/thutasann/go-react-ssr-engine/internal/config"
)

// Router holds the active route tree and provides thread-safe matching.
// Uses atomic.Pointer so the tree can be swapped on hot reload
// without any locking in the request path.
type Router struct {
	cfg  *config.Config
	tree atomic.Pointer[Tree]
}

// New creates a router and builds the initial route tree from pagesDir.
func New(cfg *config.Config) (*Router, error) {
	r := &Router{cfg: cfg}

	tree, err := r.buildTree()
	if err != nil {
		return nil, err
	}
	r.tree.Store(tree)

	return r, nil
}

// Rebuild rescans pagesDir and atomically swaps the route tree.
// In-flight requests keep using the old tree. New requests get the new one.
// No locks, no blocking, no race conditions.
func (r *Router) Rebuild() error {
	tree, err := r.buildTree()
	if err != nil {
		return err
	}
	r.tree.Store(tree)
	return nil
}

// Match finds a route for the given URL path.
// Returns the route, extracted params, and whether a match was found.
// This is called on every HTTP request — must be fast.
func (r *Router) Match(path string) (*Route, map[string]string, bool) {
	tree := r.tree.Load()
	route, params := tree.Match(path)
	if route == nil {
		return nil, nil, false
	}
	return route, params, true
}

// Routes returns all registered routes. Used for debug/logging only.
func (r *Router) Routes() []string {
	return r.collectRoutes()
}

// buildTree scans pagesDir and constructs a new Tree.
// Same file-to-route logic as the bundler so routes always match bundles.
func (r *Router) buildTree() (*Tree, error) {
	tree := NewTree()
	count := 0

	err := filepath.Walk(r.cfg.PagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".tsx" && ext != ".jsx" {
			return nil
		}

		base := filepath.Base(path)
		if strings.HasPrefix(base, "_") {
			return nil
		}

		route := filePathToRoute(path, r.cfg.PagesDir)
		tree.Insert(route, path)
		count++

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("router: scan failed: %w", err)
	}

	fmt.Printf("router: %d routes registered\n", count)
	return tree, nil
}

// collectRoutes does an in-order traversal for debug output.
func (r *Router) collectRoutes() []string {
	var routes []string
	tree := r.tree.Load()
	collectNode(tree.root, "", &routes)
	return routes
}

func collectNode(n *node, prefix string, routes *[]string) {
	if n.handler != nil {
		*routes = append(*routes, n.handler.Pattern)
	}
	for seg, child := range n.children {
		collectNode(child, prefix+"/"+seg, routes)
	}
	if n.paramChild != nil {
		collectNode(n.paramChild, prefix+"/"+n.paramChild.segment, routes)
	}
}

// filePathToRoute mirrors bundler's logic exactly.
// Duplicated intentionally — router and bundler must agree on route format
// but should not import each other (circular dependency).
func filePathToRoute(filePath string, pagesDir string) string {
	route := strings.TrimPrefix(filePath, pagesDir)
	route = strings.TrimSuffix(route, filepath.Ext(route))
	route = filepath.ToSlash(route)
	route = strings.ReplaceAll(route, "[", ":")
	route = strings.ReplaceAll(route, "]", "")

	if strings.HasSuffix(route, "/index") {
		route = strings.TrimSuffix(route, "/index")
	}
	if route == "" {
		route = "/"
	}
	return route
}
