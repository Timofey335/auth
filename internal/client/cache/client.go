package cache

import (
	"context"
)

// RedisClient - интерфейс для использования методов redis
type RedisClient interface {
	HashSet(ctx context.Context, key string, values interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	HGetAll(ctx context.Context, key string) ([]interface{}, error)
	Get(ctx context.Context, key string) (interface{}, error)
	DeleteHashSet(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, expiration int64) error
	Ping(ctx context.Context) error
}
