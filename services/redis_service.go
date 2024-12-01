// services/redis_service.go
package services

import (
	"context"
	"fmt"
	"time"

	"crypto-exchange/config"

	"github.com/go-redis/redis/v8"
)

// RedisService encapsulates the Redis client.
type RedisService struct {
	Client *redis.Client
}

// NewRedisService initializes the RedisService.
func NewRedisService(cfg config.RedisConfig) *RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Unable to connect to Redis: %v", err))
	}

	return &RedisService{
		Client: rdb,
	}
}

// Set stores a key-value pair in Redis with an expiration.
func (r *RedisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves the value for a given key from Redis.
func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}