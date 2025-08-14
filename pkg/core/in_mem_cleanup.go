package core

import (
	"context"
)

func (r *InMemoryCommandRepository) Cleanup(ctx context.Context) (deleteCount int64, err error) {
	keys, err := r.GetExpiredKeys(ctx)
	if err != nil {
		return 0, err
	}

	deleteCount = r.BatchDelete(ctx, keys)

	return deleteCount, nil
}
