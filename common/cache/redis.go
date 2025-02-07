package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	IRedisCache interface {
		ZIncrBy(key string, increment float64, member string) (float64, error)
		ZRevRangeWithScores(key string, start, end int64) ([]redis.Z, error)
	}
	RedisCache struct {
		client *redis.Client
	}
)

var RCache IRedisCache

func NewRedisCache(client *redis.Client) IRedisCache {
	return &RedisCache{
		client: client,
	}
}

const (
	REDIS_KEEP_TTL = redis.KeepTTL
)

func (r *RedisCache) ZIncrBy(key string, increment float64, member string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	score := r.client.ZIncrBy(ctx, key, increment, member)
	return score.Result()
}

func (r *RedisCache) ZRevRangeWithScores(key string, start, end int64) ([]redis.Z, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return r.client.ZRevRangeWithScores(ctx, key, start, end).Result()
}
