package hydration

import (
	"fmt"
	"strings"

	"github.com/thutasann/go-react-ssr-engine/pkg/html"
)

// Hydrator prepares the HTML document with everything the browser needs
// to "wake up" the SSR HTML into a live React app.
//
// Hydration flow:
// 1. Browser receives full HTML (user sees content immediately)
// 2. Browser downloads client JS bundle (non-blocking, type=module)
// 3. Client JS calls hydrateRoot() on #__reactgo with embedded props
// 4. React attaches event listeners without re-rendering DOM
//
// Result: instant content + fast interactivity. Best of both worlds.
type Hydrator struct {
	manifest *Manifest
}

func NewHydrator(manifest *Manifest) *Hydrator {
	return &Hydrator{manifest: manifest}
}

// Prepare configures the HTML document with client scripts and preload hints.
// Called after SSR render, before final HTML serialization.
func (h *Hydrator) Prepare(doc *html.Document, route string, propsJSON string) {
	doc.PropsJSON = propsJSON

	entry, ok := h.manifest.Get(route)
	if !ok {
		// No client bundle for this route — SSR only, no interactivity.
		// This is valid for purely static pages.
		return
	}

	doc.ClientScript = entry.JSPath

	// Preload shared chunks so browser fetches them in parallel with page JS.
	// Without preload, browser discovers chunks only after parsing the page bundle,
	// adding a full round-trip of latency.
	for _, dep := range entry.Deps {
		doc.Head.LinkTags = append(doc.Head.LinkTags, html.LinkTag{
			Rel:  "modulepreload",
			Href: dep,
			As:   "script",
		})
	}

	// Preload the page bundle itself
	doc.Head.LinkTags = append(doc.Head.LinkTags, html.LinkTag{
		Rel:  "modulepreload",
		Href: entry.JSPath,
		As:   "script",
	})

	// Optional CSS
	if entry.CSSPath != "" {
		doc.Head.Styles = append(doc.Head.Styles, html.Style{
			Href: entry.CSSPath,
		})
	}

}

// GenerateClientEntry produces the inline hydration bootstrap script.
// This is the tiny JS snippet that calls hydrateRoot() with the right
// component and props. It's the bridge between SSR HTML and live React.
//
// Output looks like:
//
//	import Page from '/_reactgo/pages/index.js';
//	import { hydrateRoot } from 'react-dom/client';
//	hydrateRoot(document.getElementById('__reactgo'),
//	  React.createElement(Page, window.__REACTGO_DATA__));
func GenerateClientEntry(route string, entry ManifestEntry) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("import Page from '%s';\n", entry.JSPath))
	sb.WriteString("import { hydrateRoot } from 'react-dom/client';\n")
	sb.WriteString("import { createElement } from 'react';\n\n")
	sb.WriteString("hydrateRoot(\n")
	sb.WriteString("  document.getElementById('__reactgo'),\n")
	sb.WriteString("  createElement(Page, window.__REACTGO_DATA__)\n")
	sb.WriteString(");\n")

	return sb.String()
}
