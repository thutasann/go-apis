package html

import (
	"fmt"
	"strings"
)

// Document is the full HTML shell wrapping the React app.
// Equivalent to Next.js _document.tsx but in Go — no JS overhead.
//
// Layout:
// <!DOCTYPE html>
// <html>
//
//	<head> ... managed by Head ... </head>
//	<body>
//	  <div id="__reactgo"> ... SSR HTML injected here ... </div>
//	  <script> ... hydration data + client bundle ... </script>
//	</body>
//
// </html>
type Document struct {
	Head         *Head
	BodyHTML     string // SSR Rendered HTML from v8
	PropsJSON    string // Serialized props for client hydration
	ClientScript string // path to client JS bundle
	Lang         string
}

// NewDocument creats a document with defaults.
func NewDocument() *Document {
	return &Document{
		Head: NewHead(),
		Lang: "en",
	}
}

// Render produces the final HTML string sent to the browser.
// String builder — no html/template parsing overhead.
// Benchmarks at ~2μs for a typical page.
func (d *Document) Render() string {
	var sb strings.Builder
	sb.Grow(4096) // most pages fit in 4KB

	sb.WriteString("<!DOCTYPE html>\n")
	sb.WriteString(fmt.Sprintf("<html lang=\"%s\">\n", d.Lang))

	// <head>
	sb.WriteString("<head>\n")
	sb.WriteString("  <meta charset=\"utf-8\">\n")
	sb.WriteString(d.Head.Render())
	sb.WriteString("</head>\n")

	// <body>
	sb.WriteString("<body>\n")

	// Root div — React hydrate targets this element.
	// ID must match the client hydration script below.
	sb.WriteString("  <div id=\"__reactgo\">")
	sb.WriteString(d.BodyHTML)
	sb.WriteString("</div>\n")

	// Hydration data — embedded as JSON in a script tag.
	// The client bundle reads this instead of making another API call.
	// This is how SSR data transfers to client without double-fetching.
	if d.PropsJSON != "" {
		sb.WriteString("  <script>window.__REACTGO_DATA__=")
		sb.WriteString(d.PropsJSON)
		sb.WriteString("</script>\n")
	}

	// Client bundle — type=module for ES module format (tree-shaking, modern syntax).
	// Loaded after SSR HTML is already painted — user sees content instantly,
	// interactivity activates once this script runs.
	if d.ClientScript != "" {
		sb.WriteString(fmt.Sprintf("  <script type=\"module\" src=\"%s\"></script>\n", d.ClientScript))
	}

	sb.WriteString("</body>\n")
	sb.WriteString("</html>")

	return sb.String()
}
