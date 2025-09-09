package storage

import (
	"context"
	"time"
)

type CacheStorage interface {
	AddSet(ctx context.Context, key string, expire time.Duration, members ...string) error
	IsSetMember(ctx context.Context, key, member string) (bool, error)
}
