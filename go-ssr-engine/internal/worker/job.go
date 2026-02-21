package worker

import (
	"net/http"

	"github.com/thutasann/go-ssr-engine/internal/engine"
)

// Job represents one rendering task.
//
// Design goals:
// - Keep struct small
// - Avoid interface{} usage
// - Avoid maps
// - Template is immutable -> safe for concurrent use
//
// IMPORTANT:
// http.ResponseWriter is an interface and will escape.
// This is unavoidable at HTTP boundary.
// Keep it here, not inside engine
type Job struct {
	Tpl *engine.Template
	Ctx engine.RenderContext
	Res http.ResponseWriter

	Done chan struct{} // signal completion
}
