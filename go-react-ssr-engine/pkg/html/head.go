package html

import (
	"fmt"
	"strings"
)

// Head manages <head> tag contents per-request.
// Each page can set its own title, meta tags, scripts, styles.
// Accumulated during render, flushed once into the HTML document.
type Head struct {
	Title    string
	MetaTags []MetaTag
	Scripts  []Script
	Styles   []Style
	LinkTags []LinkTag
}

type MetaTag struct {
	Name     string // "description", "viewport", etc.
	Property string // "og:title", "og:image" — used instead of Name for OpenGraph
	Content  string
}

type Script struct {
	Src   string
	Async bool
	Defer bool
	Type  string // "module" for ES modules
}

type Style struct {
	Href string
}

type LinkTag struct {
	Rel  string // "stylesheet", "preload", "icon"
	Href string
	As   string // "script", "style" — for preload hints
}

// NewHead returns a Head with sensible defaults every page needs.
func NewHead() *Head {
	return &Head{
		Title: "reactgo",
		MetaTags: []MetaTag{
			{Name: "viewport", Content: "width=device-width, initial-scale=1"},
		},
	}
}

// Render serializes the Head into HTML string for injection into <head>.
// No template engine — string builder is faster and has zero allocations
// after the initial grow.
func (h *Head) Render() string {
	var sb strings.Builder
	sb.Grow(512) // typical <head> is 200-400 bytes

	// Title
	sb.WriteString(fmt.Sprintf("  <title>%s</title>\n", h.Title))

	// Meta tags
	for _, m := range h.MetaTags {
		if m.Property != "" {
			sb.WriteString(fmt.Sprintf("  <meta property=\"%s\" content=\"%s\">\n", m.Property, m.Content))
		} else {
			sb.WriteString(fmt.Sprintf("  <meta name=\"%s\" content=\"%s\">\n", m.Name, m.Content))
		}
	}

	// Preload hints — tell browser to start downloading before parser reaches the tag.
	// Critical for hydration JS — shaves 50-100ms on slow connections.
	for _, l := range h.LinkTags {
		attrs := fmt.Sprintf("rel=\"%s\" href=\"%s\"", l.Rel, l.Href)
		if l.As != "" {
			attrs += fmt.Sprintf(" as=\"%s\"", l.As)
		}
		sb.WriteString(fmt.Sprintf("  <link %s>\n", attrs))
	}

	// Stylesheets
	for _, s := range h.Styles {
		sb.WriteString(fmt.Sprintf("  <link rel=\"stylesheet\" href=\"%s\">\n", s.Href))
	}

	// Scripts in head (async/defer only — blocking scripts go at body end)
	for _, s := range h.Scripts {
		attrs := fmt.Sprintf("src=\"%s\"", s.Src)
		if s.Type != "" {
			attrs += fmt.Sprintf(" type=\"%s\"", s.Type)
		}
		if s.Async {
			attrs += " async"
		}
		if s.Defer {
			attrs += " defer"
		}
		sb.WriteString(fmt.Sprintf("  <script %s></script>\n", attrs))
	}

	return sb.String()
}
