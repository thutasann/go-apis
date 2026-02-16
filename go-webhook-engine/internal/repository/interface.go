package repository

import (
	"context"

	"github.com/thutasann/go-webhook-engine/internal/domain"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error

	GetByID(ctx context.Context, id string) (*domain.Event, error)

	UpdateStatus(ctx context.Context, id string, status domain.Status) error

	IncrementRetry(ctx context.Context, id string) error
}
