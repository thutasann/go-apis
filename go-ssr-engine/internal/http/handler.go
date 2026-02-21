package http

import (
	"net/http"
	"time"

	"github.com/thutasann/go-ssr-engine/internal/engine"
	"github.com/thutasann/go-ssr-engine/internal/worker"
)

// Handler wires HTTP to worker pool.
//
// No rendering logic here.
// This boundary only.
type Handler struct {
	Pool *worker.WorkerPool
	Tpl  *engine.Template
}

// ServeHTTP implements http.Handler.
//
// IMPORTANT:
// - No heavy logic here.
// - Build RenderContxt
// - Submit to pool
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := engine.RenderContext{
		Values: [][]byte{
			[]byte("World"),
		},
	}

	done := make(chan struct{})

	job := worker.Job{
		Tpl:  h.Tpl,
		Ctx:  ctx,
		Res:  w,
		Done: done,
	}

	err := h.Pool.Submit(job, 50*time.Millisecond)
	if err != nil {
		http.Error(w, "Server Busy", http.StatusServiceUnavailable)
		return
	}

	// Wait for worker to finish rendering
	<-done
}
