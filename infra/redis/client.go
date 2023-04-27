package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go-rest-api/logger"
	"log"
	"time"
)

type Redis struct {
	*redis.Client
	databaseId int
	lgr        logger.StructLogger
}

func New(ctx context.Context, uri string, databaseId int, lgr logger.StructLogger, timeout time.Duration) (*Redis, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	connectionOption := &redis.Options{
		Addr:     uri,
		DB:       databaseId,
		Password: "",
	}

	log.Println("hitting redis connect...")

	client := redis.NewClient(connectionOption)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("redis connected...")

	rds := &Redis{
		client,
		databaseId,
		lgr,
	}

	return rds, nil
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *Redis) Set(ctx context.Context, key string, val string, exp time.Duration) error {
	return r.Client.Set(ctx, key, val, exp).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	cmd := r.Client.Get(ctx, key)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Val(), nil
}

func (r *Redis) Del(ctx context.Context, keys ...string) error {
	cmd := r.Client.Del(ctx, keys...)
	return cmd.Err()
}
