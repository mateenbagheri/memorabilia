package core

import (
	"context"
)

// CommandRepository is the interface each client needs to implement for
// interacting with memorabilia.
type CommandRepository interface {
	Get(ctx context.Context, key string) (value string, err error)
	Set(ctx context.Context, key, value string) (err error)
}
