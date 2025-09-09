package storage

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() CacheStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("failed to ping redis server: %s", err.Error())
	}

	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) AddSet(ctx context.Context, key string, member string, expire time.Duration) error {
	return nil
}

func (r *RedisCache) IsSetMember(ctx context.Context, key, member string) (bool, error) {
	return false, nil
}
