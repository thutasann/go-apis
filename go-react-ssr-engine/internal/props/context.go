package props

import "encoding/json"

// PageContext carries request-scoped data to the props loader
// Created per-request from the router match result.
// This ist he Go equivalent of Next.js getServerSideProps context.
type PageContext struct {
	// Route pattern e.g. "/posts/:id"
	Route string

	// Params extracted from URL e.g. {"id": "123"}
	Params map[string]string

	// Query string params e.g. {"sort": "date"}
	Query map[string]string

	// Path is the full request path e.g. "/posts/123?sort=date"
	Path string
}

// PageProps is what getServerSideProps returns.
// Serialized to JSON and passed to both V8 (for SSR) and browser (for hydration).
type PageProps struct {
	// Props is the actual data the component receives
	Props map[string]interface{} `json:"props"`

	// Redirect tells the server to send a 302 instead of rendering
	Redirect *Redirect `json:"redirect,omitempty"`

	// NotFound triggers a 404 page
	NotFound bool `json:"notFound,omitempty"`
}

// Redirect holds redirect target info.
type Redirect struct {
	Destination string `json:"destination"`
	Permanent   bool   `json:"permanent"` // 301 vs 302
}

// ToJSON serializes props for passing to V8 and embedding in HTML.
// Returns "{}" on empty props — never null, React components expect an object.
func (p *PageProps) ToJSON() (string, error) {
	if p == nil || p.Props == nil {
		return "{}", nil
	}
	data, err := json.Marshal(p.Props)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
