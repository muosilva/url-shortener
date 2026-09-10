package redis

import (
	"context"
	"errors"
	"time"

	"github.com/muosilva/url-shortener/metrics"
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
	start := time.Now()
	result := metrics.ResultHit
	defer func() {
		metrics.ObserveCacheOperation("get", result, time.Since(start))
	}()

	longURL, err := c.client.Get(ctx, code).Result()
	if errors.Is(err, redis.Nil) {
		result = metrics.ResultMiss
		metrics.CacheMissesTotal.Inc()
		return "", nil
	}
	if err != nil {
		result = metrics.ResultError
		return "", err
	}

	metrics.CacheHitsTotal.Inc()
	return longURL, nil
}

func (c *RedisURLCache) Set(ctx context.Context, code string, longURL string) error {
	start := time.Now()
	result := metrics.ResultSuccess
	defer func() {
		metrics.ObserveCacheOperation("set", result, time.Since(start))
	}()

	err := c.client.Set(ctx, code, longURL, cacheTTL).Err()
	if err != nil {
		result = metrics.ResultError
		return err
	}

	return nil
}
