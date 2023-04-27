package redis

import (
	"context"
	"time"
)

type DB interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, val string, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
}
