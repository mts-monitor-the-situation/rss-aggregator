package redisdb

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Connect establishes a connection to the Redis server.
func Connect(addr string) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr: addr,
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
