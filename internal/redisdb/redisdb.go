package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Connect establishes a connection to the Redis server.
func Connect(addr string) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	// Test the connection
	ctxPing, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctxPing).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// AddToStream adds a message to a Redis stream with the specified name.
func AddToStream(client *redis.Client, streamName string, message map[string]any) error {
	ctx := context.Background()
	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: message,
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to add message to stream %s: %w", streamName, err)
	}

	return nil
}
