package storage

import (
	"context"
	"time"
)

type CacheStorage interface {
	AddSet(ctx context.Context, key string, member string, expire time.Duration) error
	IsSetMember(ctx context.Context, key, member string) (bool, error)
}
