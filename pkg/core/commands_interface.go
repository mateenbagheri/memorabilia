package core

import (
	"context"
	"time"
)

// CommandsRepository is the interface each client needs to implement for
// interacting with memorabilia.
type CommandsRepository interface {
	Get(ctx context.Context, key string) (value string, err error)
	Set(ctx context.Context, key, value string, expiration time.Time) (err error)
	BatchDelete(ctx context.Context, keys []string) (deleteCount int64)
	Delete(ctx context.Context, key string) (deleteCount int64)
	GetExpiredKeys(ctx context.Context) (keys []string, err error)
}
