package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
	key    string
}

func NewRedisQueue(client *redis.Client, key string) *RedisQueue {
	return &RedisQueue{
		client: client,
		key:    key,
	}
}

func (q *RedisQueue) Enqueue(ctx context.Context, eventID string) error {
	return q.client.LPush(ctx, q.key, eventID).Err()
}

func (q *RedisQueue) Dequeue(ctx context.Context) (string, error) {
	// 0 = block forever
	result, err := q.client.BRPop(ctx, 0, q.key).Result()
	if err != nil {
		return "", err
	}

	// BRPop returns [key, value]
	return result[1], nil
}
