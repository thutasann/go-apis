package router

import (
	"testing"
)

func TestStaticRoutes(t *testing.T) {
	tree := NewTree()
	tree.Insert("/", "pages/index.tsx")
	tree.Insert("/about", "pages/about.tsx")
	tree.Insert("/blog/archive", "pages/blog/archive.tsx")

	tests := []struct {
		path      string
		wantMatch bool
		wantPage  string
	}{
		{"/", true, "pages/index.tsx"},
		{"/about", true, "pages/about.tsx"},
		{"/blog/archive", true, "pages/blog/archive.tsx"},
		{"/nonexistent", false, ""},
		{"/about/extra", false, ""},
	}

	for _, tt := range tests {
		route, _ := tree.Match(tt.path)
		if tt.wantMatch && route == nil {
			t.Errorf("path %s: expected match, got nil", tt.path)
		}
		if !tt.wantMatch && route != nil {
			t.Errorf("path %s: expected no match, got %s", tt.path, route.PagePath)
		}
		if tt.wantMatch && route != nil && route.PagePath != tt.wantPage {
			t.Errorf("path %s: expected page %s, got %s", tt.path, tt.wantPage, route.PagePath)
		}
	}
}

func TestDynamicRoutes(t *testing.T) {
	tree := NewTree()
	tree.Insert("/posts/:id", "pages/posts/[id].tsx")
	tree.Insert("/users/:uid/posts/:pid", "pages/users/[uid]/posts/[pid].tsx")

	// Single param
	route, params := tree.Match("/posts/123")
	if route == nil {
		t.Fatal("/posts/123: expected match")
	}
	if params["id"] != "123" {
		t.Errorf("expected id=123, got %s", params["id"])
	}

	// Nested params
	route, params = tree.Match("/users/alice/posts/456")
	if route == nil {
		t.Fatal("/users/alice/posts/456: expected match")
	}
	if params["uid"] != "alice" {
		t.Errorf("expected uid=alice, got %s", params["uid"])
	}
	if params["pid"] != "456" {
		t.Errorf("expected pid=456, got %s", params["pid"])
	}
}

func TestStaticOverParam(t *testing.T) {
	// Static routes should match before param routes at the same level.
	// "/posts/featured" is static, "/posts/:id" is dynamic.
	tree := NewTree()
	tree.Insert("/posts/:id", "pages/posts/[id].tsx")
	tree.Insert("/posts/featured", "pages/posts/featured.tsx")

	route, _ := tree.Match("/posts/featured")
	if route == nil {
		t.Fatal("expected match for /posts/featured")
	}
	if route.PagePath != "pages/posts/featured.tsx" {
		t.Errorf("expected static match, got %s", route.PagePath)
	}

	route, params := tree.Match("/posts/999")
	if route == nil {
		t.Fatal("expected match for /posts/999")
	}
	if params["id"] != "999" {
		t.Errorf("expected id=999, got %s", params["id"])
	}
}

func BenchmarkTreeMatch(b *testing.B) {
	tree := NewTree()
	// Simulate a medium-size app with 50 routes
	for i := 0; i < 50; i++ {
		tree.Insert(
			"/section"+string(rune('a'+i%26))+"/:id",
			"pages/section.tsx",
		)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tree.Match("/sectionm/12345")
		}
	})
}
