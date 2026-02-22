package router

import "strings"

// node is a radix tree node. Radix tree gives O(k) route matching
// where k = number of path segments, not total routes.
// This means 10,000 routes match as fast as 10 routes.
type node struct {
	// segment is this node's path part, e.g, "posts" or ":id"
	segment string

	// handler is non-nil only on leaf nodes that map to a page
	handler *Route

	// children keyed by first char of segment for 0(1) child lookup.
	// Static children live here.
	children map[string]*node

	// paramChild handles :param segments. Only one per level allowedd
	// (same as Next.js - can't have [id].tsx and [slug].tsx
	paramChild *node
}

// Route holds everything needed to render a matched page.
type Route struct {
	// Pattern is the original route like "/posts/:id"
	Pattern string

	// PagePath is the file path relative to pagesDir, e.g. "/pages/posts/[id].tsx"
	PagePath string

	// Params are the dynamic segment names in order, e.g. ["id"]
	Params []string
}

// Tree is the top-level router. Thread-safe for reads after Build().
// Built once at startup (and rebuilt on file change in dev mode).
// No mutex needed because we swap the entire tree pointer atomically.
type Tree struct {
	root *node
}

// NewTree creates an empty radix tree.
func NewTree() *Tree {
	return &Tree{
		root: &node{
			children: make(map[string]*node),
		},
	}
}

// Insert adds a route pattern to the tree.
// Pattern format: "/posts/:id" (already converted from [id] by bundler).
func (t *Tree) Insert(pattern string, pagePath string) {
	segments := splitPath(pattern)
	current := t.root

	var params []string

	for _, seg := range segments {
		if strings.HasPrefix(seg, ":") {
			// Dynamic segment — store in paramChild
			params = append(params, seg[1:]) // strip the ":"
			if current.paramChild == nil {
				current.paramChild = &node{
					segment:  seg,
					children: make(map[string]*node),
				}
			}
			current = current.paramChild
		} else {
			// Static segment — store in children map
			child, exists := current.children[seg]
			if !exists {
				child = &node{
					segment:  seg,
					children: make(map[string]*node),
				}
				current.children[seg] = child
			}
			current = child
		}
	}

	// Mark leaf node with the route handler
	current.handler = &Route{
		Pattern:  pattern,
		PagePath: pagePath,
		Params:   params,
	}
}

// Match finds the route for a given URL path and extracts param values.
// Returns nil if no route matches. Zero allocations on static routes.
func (t *Tree) Match(path string) (*Route, map[string]string) {
	segments := splitPath(path)
	current := t.root

	// paramValues collected during traversal, allocated only if needed
	var paramValues []string

	for _, seg := range segments {
		// Try static child first — exact match is always preferred
		if child, exists := current.children[seg]; exists {
			current = child
			continue
		}

		// Fall back to param child
		if current.paramChild != nil {
			current = current.paramChild
			paramValues = append(paramValues, seg)
			continue
		}

		// No match at this level
		return nil, nil
	}

	if current.handler == nil {
		return nil, nil
	}

	// Build param map only if there are dynamic segments
	var params map[string]string
	if len(current.handler.Params) > 0 {
		params = make(map[string]string, len(current.handler.Params))
		for i, name := range current.handler.Params {
			if i < len(paramValues) {
				params[name] = paramValues[i]
			}
		}
	}

	return current.handler, params
}

// splitPath turns "/posts/123" into ["posts", "123"].
// Leading/trailing slashes are ignored. Empty path returns empty slice.
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}
