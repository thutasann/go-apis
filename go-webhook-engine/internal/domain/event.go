package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	IdempotencyKey string             `bson:"idempotency_key"`
	Type           string             `bson:"type"`
	Payload        []byte             `bson:"payload"`

	Status     Status `bson:"status"`
	RetryCount int    `bson:"retry_count"`
	MaxRetries int    `bson:"max_retries"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
