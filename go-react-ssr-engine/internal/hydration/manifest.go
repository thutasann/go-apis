package hydration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Manifest maps route patterns to their client-side JS bundle paths.
// Built after esbuild compiles client bundles. Read on every request
// to inject the correct <script> tag for the matched page.
//
// Thread-safe: atomic swap on rebuild, RLock-free reads via sync.RWMutex.
type Manifest struct {
	mu      sync.RWMutex
	entries map[string]ManifestEntry
}

// ManifestEntry holds the client bundle info for one route.
type ManifestEntry struct {
	// JSPath is the URL path to the client JS file, e.g. "/_reactgo/pages/index.js"
	JSPath string `json:"js"`

	// CSSPath is optional — populated if the page imports CSS
	CSSPath string `json:"css,omitempty"`

	// Deps lists shared chunk paths this page depends on (React, etc.)
	// Preloaded via <link rel="modulepreload"> for faster hydration.
	Deps []string `json:"deps,omitempty"`
}

func NewManifest() *Manifest {
	return &Manifest{
		entries: make(map[string]ManifestEntry),
	}
}

// Build scans esbuild client output dir and maps routes to bundle paths.
// Called after every successful build.
func (m *Manifest) Build(clientDir string, routeMap map[string]string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Reset entries for clean rebuild
	m.entries = make(map[string]ManifestEntry, len(routeMap))

	// Discover shared chunks (React, common utils)
	chunks, err := m.findChunks(clientDir)
	if err != nil {
		return err
	}

	for route, filePath := range routeMap {
		// Convert filesystem path to URL path served by static handler
		urlPath := "/_reactgo/" + strings.TrimPrefix(filePath, clientDir+"/")

		m.entries[route] = ManifestEntry{
			JSPath: urlPath,
			Deps:   chunks,
		}
	}

	return nil
}

// Get returns the manifest entry for a route. Called per-request.
func (m *Manifest) Get(route string) (ManifestEntry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.entries[route]
	return entry, ok
}

// findChunks locates shared JS chunks in the chunks/ subdirectory.
// These are common dependencies esbuild extracted (React, ReactDOM, etc).
func (m *Manifest) findChunks(clientDir string) ([]string, error) {
	chunksDir := filepath.Join(clientDir, "chunks")
	if _, err := os.Stat(chunksDir); os.IsNotExist(err) {
		return nil, nil // no chunks yet
	}

	var chunks []string
	entries, err := os.ReadDir(chunksDir)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".js") {
			chunks = append(chunks, "/_reactgo/chunks/"+e.Name())
		}
	}

	return chunks, nil
}

// WriteJSON saves manifest to disk for debugging and external tooling.
func (m *Manifest) WriteJSON(path string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := json.MarshalIndent(m.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
