package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a Redis client from the REDIS_URL environment variable.
// Falls back to localhost:6379 if the variable is not set.
func NewRedisClient() (*redis.Client, error) {
	rawURL := os.Getenv("REDIS_URL")

	if rawURL == "" {
		rawURL = "redis://localhost:6379"
	}
	opts, err := redis.ParseURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_URL: %w", err)
	}
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return client, nil
}
