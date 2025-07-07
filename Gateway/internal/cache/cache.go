package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, token string) (uuid string, ok bool, err error)
	Save(ctx context.Context, token string, uuid string, expirationTime time.Duration) error
}
