package storage

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	noNewSetMembers = 0
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

func (r *RedisCache) AddSet(ctx context.Context, key string, expire time.Duration, members ...string) error {
	count, err := r.client.SAdd(ctx, key, members).Result()
	if err != nil {
		return err
	}

	if count == noNewSetMembers {
		return nil
	}

	_, err = r.client.ExpireNX(ctx, key, expire).Result()

	return err
}

func (r *RedisCache) IsSetMember(ctx context.Context, key, member string) (bool, error) {
	isMember, err := r.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		return false, err
	}

	return isMember, nil
}
