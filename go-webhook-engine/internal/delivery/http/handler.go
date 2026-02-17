package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/thutasann/go-webhook-engine/internal/domain"
	"github.com/thutasann/go-webhook-engine/internal/queue"
	"github.com/thutasann/go-webhook-engine/internal/repository"
)

type Handler struct {
	repo  repository.EventRepository
	queue queue.Queue
}

func NewHandler(repo repository.EventRepository, queue queue.Queue) *Handler {
	return &Handler{
		repo:  repo,
		queue: queue,
	}
}

type WebhookRequest struct {
	IdempotencyKey string          `json:"idempotency_key"`
	Type           string          `json:"type"`
	Payload        json.RawMessage `json:"payload"`
}

func (h *Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	var req WebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.IdempotencyKey == "" || req.Type == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	event := &domain.Event{
		IdempotencyKey: req.IdempotencyKey,
		Type:           req.Type,
		Payload:        req.Payload,
		MaxRetries:     3,
		Status:         domain.StatusPending,
		RetryCount:     0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.repo.Create(ctx, event); err != nil {
		http.Error(w, "failed to persist event", http.StatusInternalServerError)
		return
	}

	if err := h.queue.Enqueue(ctx, event.ID.Hex()); err != nil {
		http.Error(w, "failed to enqueue event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}
