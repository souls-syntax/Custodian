package cache

import (
	"context"
	"time"
	"github.com/souls-syntax/Templates/internal/models"
	"github.com/redis/go-redis/v9"
	"encoding/json"
)

type RedisCache struct {
	client				*redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client:client}
}

func (r *RedisCache) Get(ctx context.Context, key string) (*models.CacheDecision, bool) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false
	}
	if err != nil {
		return nil, false // infra failure treated as miss
	}

	var res models.CacheDecision
	if err := json.Unmarshal([]byte(val), &res); err != nil {
		return nil, false
	}

	return &res, true
}

func (r *RedisCache) Set(ctx context.Context, key string, val models.CacheDecision, ttl time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}
