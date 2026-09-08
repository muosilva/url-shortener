package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const cacheTTL = 5 * time.Minute

type RedisURLCache struct {
	client *redis.Client
}

func NewRedisURLCache(client *redis.Client) *RedisURLCache {
	return &RedisURLCache{
		client: client,
	}
}

func (c *RedisURLCache) Get(ctx context.Context, code string) (string, error) {
	return c.client.Get(ctx, code).Result()
}

func (c *RedisURLCache) Set(ctx context.Context, code string, longURL string) error {
	return c.client.Set(ctx, code, longURL, cacheTTL).Err()
}
