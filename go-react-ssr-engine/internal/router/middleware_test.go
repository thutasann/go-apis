package router

import (
	"testing"
)

func TestChainOrder(t *testing.T) {
	// Verify middleware executes in correct order: A -> B -> handler
	var order []string

	mwA := func(next HandlerFunc) HandlerFunc {
		return func(ctx *RequestContext) (string, error) {
			order = append(order, "A-before")
			html, err := next(ctx)
			order = append(order, "A-after")
			return html, err
		}
	}

	mwB := func(next HandlerFunc) HandlerFunc {
		return func(ctx *RequestContext) (string, error) {
			order = append(order, "B-before")
			html, err := next(ctx)
			order = append(order, "B-after")
			return html, err
		}
	}

	handler := func(ctx *RequestContext) (string, error) {
		order = append(order, "handler")
		return "ok", nil
	}

	chain := Chain(mwA, mwB)
	wrapped := chain(handler)

	ctx := NewRequestContext("/test")
	wrapped(ctx)

	expected := []string{"A-before", "B-before", "handler", "B-after", "A-after"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("step %d: expected %s, got %s", i, v, order[i])
		}
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	// Recovery should catch panics and return 500 instead of crashing
	handler := func(ctx *RequestContext) (string, error) {
		panic("render exploded")
	}

	recovered := Recovery()(handler)
	ctx := NewRequestContext("/boom")

	html, err := recovered(ctx)
	if err != nil {
		t.Errorf("expected nil error after recovery, got %v", err)
	}
	if ctx.StatusCode != 500 {
		t.Errorf("expected status 500, got %d", ctx.StatusCode)
	}
	if html == "" {
		t.Error("expected error HTML, got empty")
	}
}
