package core

import (
	"context"
	"time"
)

// CommandRepository is the interface each client needs to implement for
// interacting with memorabilia.
type CommandRepository interface {
	Get(ctx context.Context, key string) (value string, err error)
	Set(ctx context.Context, key, value string, expiration time.Time) (err error)
}
