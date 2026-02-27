package bundler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/thutasann/go-react-ssr-engine/internal/config"
)

type Bundler struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Bundler {
	return &Bundler{cfg: cfg}
}

type BuildResult struct {
	ServerBundle  string
	ClientEntries map[string]string
}

func (b *Bundler) Build() (*BuildResult, error) {
	entries, err := b.discoverPages()
	if err != nil {
		return nil, fmt.Errorf("bundler: page discovery failed: %w", err)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("bundler: no pages found in %s", b.cfg.PagesDir)
	}

	serverJS, err := b.buildServer(entries)
	if err != nil {
		return nil, fmt.Errorf("bundler: server build failed: %w", err)
	}

	clientEntries, err := b.buildClient(entries)
	if err != nil {
		return nil, fmt.Errorf("bundler: client build failed: %w", err)
	}

	return &BuildResult{
		ServerBundle:  serverJS,
		ClientEntries: clientEntries,
	}, nil
}

func (b *Bundler) discoverPages() ([]string, error) {
	var entries []string

	err := filepath.Walk(b.cfg.PagesDir, func(path string, info os.FileInfo, err error) error {
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

		entries = append(entries, path)
		return nil
	})

	return entries, err
}

func (b *Bundler) buildServer(entries []string) (string, error) {
	virtualEntry := b.generateServerEntry(entries)

	serverDir := filepath.Join(b.cfg.BuildDir, "server")
	os.MkdirAll(serverDir, 0755)

	entryPath := filepath.Join(serverDir, "_entry.jsx")
	if err := os.WriteFile(entryPath, []byte(virtualEntry), 0644); err != nil {
		return "", err
	}

	// Resolve node_modules from project root so esbuild finds react/react-dom
	absNodeModules, _ := filepath.Abs("node_modules")

	result := api.Build(api.BuildOptions{
		EntryPoints:      []string{entryPath},
		Bundle:           true,
		Write:            false,
		Platform:         api.PlatformNeutral,
		Format:           api.FormatIIFE,
		Target:           api.ES2020,
		JSX:              api.JSXAutomatic,
		Sourcemap:        api.SourceMapNone,
		MinifySyntax:     !b.cfg.Dev,
		MinifyWhitespace: !b.cfg.Dev,

		// Tell esbuild where to find node_modules.
		// Without this, react/react-dom imports fail.
		NodePaths: []string{absNodeModules},

		// Define process.env.NODE_ENV so React uses production build
		// in prod mode (smaller, faster) and development build in dev
		// mode (better error messages).
		Define: map[string]string{
			"process.env.NODE_ENV": fmt.Sprintf(`"%s"`, b.envMode()),
		},
	})

	if len(result.Errors) > 0 {
		return "", fmt.Errorf("esbuild server: %s", result.Errors[0].Text)
	}

	return string(result.OutputFiles[0].Contents), nil
}

func (b *Bundler) buildClient(entries []string) (map[string]string, error) {
	clientDir := filepath.Join(b.cfg.BuildDir, "client")
	os.MkdirAll(clientDir, 0755)

	// Generate per-page hydration entry files.
	// Each page gets a tiny JS file that imports the component
	// and calls hydrateRoot. esbuild bundles each one separately.
	hydrateEntries, err := b.generateClientEntries(entries, clientDir)
	if err != nil {
		return nil, err
	}

	absNodeModules, _ := filepath.Abs("node_modules")

	result := api.Build(api.BuildOptions{
		EntryPoints:      hydrateEntries,
		Bundle:           true,
		Write:            true,
		Outdir:           clientDir,
		Platform:         api.PlatformBrowser,
		Format:           api.FormatESModule,
		Target:           api.ES2020,
		JSX:              api.JSXAutomatic,
		Splitting:        true,
		ChunkNames:       "chunks/[name]-[hash]",
		Sourcemap:        api.SourceMapLinked,
		MinifySyntax:     !b.cfg.Dev,
		MinifyWhitespace: !b.cfg.Dev,
		NodePaths:        []string{absNodeModules},
		Define: map[string]string{
			"process.env.NODE_ENV": fmt.Sprintf(`"%s"`, b.envMode()),
		},
	})

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("esbuild client: %s", result.Errors[0].Text)
	}

	// Map route -> output JS URL path
	clientMap := make(map[string]string)
	for _, entry := range entries {
		route := b.filePathToRoute(entry)
		// Hydrate entry mirrors page structure: pages/index.tsx -> _hydrate_index.js
		name := b.hydrateEntryName(entry)
		outFile := filepath.Join(clientDir, name+".js")
		if _, err := os.Stat(outFile); err == nil {
			clientMap[route] = outFile
		}
	}

	return clientMap, nil
}

// generateClientEntries creates per-page hydration scripts.
// Each script imports its page component and calls hydrateRoot.
// These become the entrypoints for the client build.
func (b *Bundler) generateClientEntries(entries []string, clientDir string) ([]string, error) {
	var hydrateFiles []string

	for _, entry := range entries {
		absPage, _ := filepath.Abs(entry)
		name := b.hydrateEntryName(entry)

		// Tiny hydration bootstrap — this is all the client-specific JS per page.
		// React, ReactDOM are shared chunks extracted by esbuild splitting.
		script := fmt.Sprintf(`import { hydrateRoot } from 'react-dom/client';
import { createElement } from 'react';
import Page from '%s';

const container = document.getElementById('__reactgo');
const props = window.__REACTGO_DATA__ || {};
hydrateRoot(container, createElement(Page, props));
`, absPage)

		hydratePath := filepath.Join(clientDir, name+".jsx")
		if err := os.WriteFile(hydratePath, []byte(script), 0644); err != nil {
			return nil, err
		}
		hydrateFiles = append(hydrateFiles, hydratePath)
	}

	return hydrateFiles, nil
}

// generateServerEntry creates the V8 entry that registers all pages
// and exposes __renderToString and __getServerSideProps globals.
func (b *Bundler) generateServerEntry(entries []string) string {
	var sb strings.Builder

	sb.WriteString(`var React = require('react');
var ReactDOMServer = require('react-dom/server');

var routes = {};
var propsLoaders = {};

`)

	for i, entry := range entries {
		absPath, _ := filepath.Abs(entry)
		route := b.filePathToRoute(entry)

		sb.WriteString(fmt.Sprintf("var Page%d = require('%s');\n", i, absPath))
		// Support both default and named exports
		sb.WriteString(fmt.Sprintf("var Comp%d = Page%d.default || Page%d;\n", i, i, i))
		sb.WriteString(fmt.Sprintf("routes['%s'] = Comp%d;\n", route, i))

		// Register getServerSideProps if exported
		sb.WriteString(fmt.Sprintf("if (Page%d.getServerSideProps) { propsLoaders['%s'] = Page%d.getServerSideProps; }\n\n", i, route, i))
	}

	// Global render bridge — called by Worker.Execute()
	sb.WriteString(`
globalThis.__renderToString = function(route, props) {
  var Component = routes[route];
  if (!Component) {
    return '<div>404 - Page not found</div>';
  }
  try {
    return ReactDOMServer.renderToString(React.createElement(Component, props));
  } catch(e) {
    return '<div>Render Error: ' + e.message + '</div>';
  }
};

globalThis.__getServerSideProps = function(route, context) {
  var loader = propsLoaders[route];
  if (!loader) {
    return JSON.stringify({ props: {} });
  }
  try {
    var result = loader(context);
    return JSON.stringify(result);
  } catch(e) {
    return JSON.stringify({ props: {}, error: e.message });
  }
};

globalThis.__hasServerProps = function(route) {
  return !!propsLoaders[route];
};
`)

	return sb.String()
}

func (b *Bundler) filePathToRoute(filePath string) string {
	route := strings.TrimPrefix(filePath, b.cfg.PagesDir)
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

// hydrateEntryName creates a unique flat filename for hydration entries.
// pages/posts/[id].tsx -> _hydrate_posts_id
func (b *Bundler) hydrateEntryName(filePath string) string {
	name := strings.TrimPrefix(filePath, b.cfg.PagesDir)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = filepath.ToSlash(name)
	name = strings.Trim(name, "/")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "[", "")
	name = strings.ReplaceAll(name, "]", "")
	return "_hydrate_" + name
}

func (b *Bundler) envMode() string {
	if b.cfg.Dev {
		return "development"
	}
	return "production"
}
