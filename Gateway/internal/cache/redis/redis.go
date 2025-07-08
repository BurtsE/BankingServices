package redis

import (
	"context"
	"errors"
	"fmt"
	"gateway/internal/cache"
	"gateway/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

var _ cache.Cache = (*Redis)(nil)

type Redis struct {
	client *redis.Client
}

func NewRedisCache(cfg *config.Config) (*Redis, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		//Password: cfg.Redis.Password,
		//DB:       cfg.Redis.Database,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Redis{client: rdb}, nil
}

func (r *Redis) Get(ctx context.Context, token string) (string, bool, error) {
	uuid, err := r.client.Get(ctx, token).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	
	return uuid, true, nil
}

func (r *Redis) Save(ctx context.Context, token string, uuid string, expirationTime time.Duration) error {
	return r.client.Set(ctx, token, uuid, expirationTime).Err()
}

func (r *Redis) Close(ctx context.Context) error {
	return r.client.Close()
}
