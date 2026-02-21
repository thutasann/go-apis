package http

import (
	"net/http"
	"time"

	"github.com/thutasann/go-ssr-engine/internal/engine"
	"github.com/thutasann/go-ssr-engine/internal/worker"
)

// Handler wires HTTP to worker pool with dynamic variable injection
type Handler struct {
	Pool           *worker.WorkerPool
	Tpl            *engine.Template
	VarNameToIndex map[string]uint16 // mapping from compiler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Read query parameters as template variables
	varMap := map[string]string{}
	for key := range h.VarNameToIndex {
		val := r.URL.Query().Get(key)
		varMap[key] = val
	}

	ctx := engine.NewRenderContext(varMap, h.Tpl, h.VarNameToIndex)

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

	// Wait for worker to finish
	<-done
}
