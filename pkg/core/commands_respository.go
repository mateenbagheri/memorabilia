package core

import (
	"context"
)

type CommandRepository interface {
	Get(ctx context.Context, key string) (value string, err error)
	Set(ctx context.Context, key, value string) (err error)
}
