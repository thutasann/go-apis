package repository

import (
	"context"
	"errors"
	"time"

	"github.com/thutasann/go-webhook-engine/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoEventRepository struct {
	collection *mongo.Collection
}

func NewMongoEventRepository(db *mongo.Database, collectionName string) *MongoEventRepository {
	return &MongoEventRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *MongoEventRepository) Create(ctx context.Context, event *domain.Event) error {
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now
	event.Status = domain.StatusPending

	_, err := r.collection.InsertOne(ctx, event)
	return err
}

func (r *MongoEventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var event domain.Event
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *MongoEventRepository) UpdateStatus(ctx context.Context, id string, status domain.Status) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateByID(ctx, objID, update)
	return err
}

func (r *MongoEventRepository) IncrementRetry(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$inc": bson.M{
			"retry_count": 1,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("event not found")
	}

	return nil
}
