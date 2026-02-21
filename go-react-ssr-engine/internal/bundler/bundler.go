package bundler

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/thutasann/go-react-ssr-engine/internal/config"

	"github.com/evanw/esbuild/pkg/api"
)

// Bundler compiles React pages into two outputs:
// 1. Server bundle - single file, runs in V8 for SSR
// 2. Client bundles - per-page chunks, shipped to browser for hydration.
type Bundler struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Bundler {
	return &Bundler{cfg: cfg}
}

// BuildResult holds path and content from a successful build.
type BuildResult struct {
	// ServerBundle is the JS string to load into V8 via engine.LoadBundle()
	ServerBundle string

	// ClientEntries map route path -> client JS file path
	// Used by hydration to inject the right <script> per page
	ClientEntries map[string]string
}

// build scans pagesDir for all .jsx/.tsx files and compiles them.
// Called once at startup and again on every file change in dev mode.
func (b *Bundler) Build() (*BuildResult, error) {
	// Discover all page entry points
	entries, err := b.discoverPages()
	if err != nil {
		return nil, fmt.Errorf("bundler: page discovery failed: %w", err)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("bundler: no pages found in %s", b.cfg.PagesDir)
	}

	// Build server bundle — all pages in one file, no code splitting.
	// V8 loads one blob, fast cold start.
	serverJS, err := b.buildServer(entries)
	if err != nil {
		return nil, fmt.Errorf("bundler: server build failed: %w", err)
	}

	// Build client bundles — per-page splitting so browser only
	// downloads JS for the current page, not the entire app.
	clientEntries, err := b.buildClient(entries)
	if err != nil {
		return nil, fmt.Errorf("bundler: client build failed: %w", err)
	}

	return &BuildResult{
		ServerBundle:  serverJS,
		ClientEntries: clientEntries,
	}, nil
}

// discoverPages walks pagesDir and returns .tsx/.jsx files as entry points.
// Skips files starting with _ (like _app.tsx, _document.tsx - special files)
func (b *Bundler) discoverPages() ([]string, error) {
	var entries []string

	err := filepath.Walk(b.cfg.PagesDir, func(path string, info fs.FileInfo, err error) error {
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

		// Skip special files — these are handled separately
		base := filepath.Base(path)
		if strings.HasPrefix(base, "_") {
			return nil
		}

		entries = append(entries, path)
		return nil
	})

	return entries, err
}

// buildServer creates a single JS bundle for V8.
// Wraps all page components in a route map and exposes __renderToString
func (b *Bundler) buildServer(entries []string) (string, error) {
	// Generate a virtual entry that imports all pages and builds a route map.
	// This becomes the single entry point esbuild compiles.
	virtualEntry := b.generateServerEntry(entries)

	serverDir := filepath.Join(b.cfg.BuildDir, "server")
	os.MkdirAll(serverDir, 0755)

	// Write virtual entry to disk — esbuild needs a real file
	entryPath := filepath.Join(serverDir, "_entry.jsx")
	if err := os.WriteFile(entryPath, []byte(virtualEntry), 0644); err != nil {
		return "", err
	}

	result := api.Build(api.BuildOptions{
		EntryPoints:      []string{entryPath},
		Bundle:           true,                // Resolve all imports into one file
		Write:            false,               // Keep in memory — we pass string to V8
		Platform:         api.PlatformNeutral, // Not node, not browser — V8 standalone
		Format:           api.FormatIIFE,      // Immediately invoked — defines globals on exec
		Target:           api.ES2020,
		JSX:              api.JSXAutomatic, // React 17+ JSX transform, no manual import React
		Sourcemap:        api.SourceMapNone,
		MinifySyntax:     !b.cfg.Dev,
		MinifyWhitespace: !b.cfg.Dev,
	})

	if len(result.Errors) > 0 {
		return "", fmt.Errorf("esbuild server: %s", result.Errors[0].Text)
	}

	return string(result.OutputFiles[0].Contents), nil
}

// buildClient creates per-page bundles for browser hydration.
// Each page gets its own chunk so the browser only loads what it needs.
func (b *Bundler) buildClient(entries []string) (map[string]string, error) {
	clientDir := filepath.Join(b.cfg.BuildDir, "client")
	os.MkdirAll(clientDir, 0755)

	result := api.Build(api.BuildOptions{
		EntryPoints:      entries,
		Bundle:           true,
		Write:            true, // Write to disk — served as static files
		Outdir:           clientDir,
		Platform:         api.PlatformBrowser,
		Format:           api.FormatESModule, // ES modules for modern browsers
		Target:           api.ES2020,
		JSX:              api.JSXAutomatic,
		Splitting:        true, // Shared chunks for common dependencies (React, etc)
		ChunkNames:       "chunks/[name]-[hash]",
		Sourcemap:        api.SourceMapLinked, // Separate .map files, not inlined
		MinifySyntax:     !b.cfg.Dev,
		MinifyWhitespace: !b.cfg.Dev,
	})

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("esbuild client: %s", result.Errors[0].Text)
	}

	// Map route paths to output file paths
	clientEntries := make(map[string]string)
	for _, entry := range entries {
		route := b.filePathToRoute(entry)
		// esbuild mirrors input structure in outdir
		outFile := filepath.Join(clientDir, strings.TrimPrefix(entry, b.cfg.PagesDir))
		outFile = strings.TrimSuffix(outFile, filepath.Ext(outFile)) + ".js"
		clientEntries[route] = outFile
	}

	return clientEntries, nil
}

// generateServerEntry creates JS that imports all pages into a route map
// and exposes a global __renderToString function for V8 to call.
func (b *Bundler) generateServerEntry(entries []string) string {
	var sb strings.Builder

	sb.WriteString("import { renderToString } from 'react-dom/server';\n")
	sb.WriteString("import { createElement } from 'react';\n\n")

	// Import each page component with a safe variable name
	sb.WriteString("const routes = {};\n\n")
	for i, entry := range entries {
		route := b.filePathToRoute(entry)
		// Use absolute path for reliable resolution
		absPath, _ := filepath.Abs(entry)
		sb.WriteString(fmt.Sprintf("import Page%d from '%s';\n", i, absPath))
		sb.WriteString(fmt.Sprintf("routes['%s'] = Page%d;\n\n", route, i))
	}

	// The global render bridge — called by Worker.Execute()
	sb.WriteString(`
globalThis.__renderToString = function(route, props) {
  const Component = routes[route];
  if (!Component) {
    return '<div>404 - Page not found</div>';
  }
  return renderToString(createElement(Component, props));
};
`)

	return sb.String()
}

// filePathToRoute converts "pages/posts/[id].tsx" -> "/posts/:id"
// This is the Next.js-style file system routing convention.
func (b *Bundler) filePathToRoute(filePath string) string {
	// Strip pagesDir prefix and extension
	route := strings.TrimPrefix(filePath, b.cfg.PagesDir)
	route = strings.TrimSuffix(route, filepath.Ext(route))

	// Normalize separators
	route = filepath.ToSlash(route)

	// Convert [param] -> :param for radix tree matching
	route = strings.ReplaceAll(route, "[", ":")
	route = strings.ReplaceAll(route, "]", "")

	// /index -> /
	if strings.HasSuffix(route, "/index") {
		route = strings.TrimSuffix(route, "/index")
	}
	if route == "" {
		route = "/"
	}

	return route
}
