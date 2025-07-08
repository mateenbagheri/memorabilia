package core

import (
	"context"
)

// Delete removes a single key-value pair from the in-memory store if it exists.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - key: The key of the key-value pair to delete.
func (imc *InMemoryCommandRepository) Delete(ctx context.Context, key string) (deleteCount int64) {
	if val, exists := imc.store[key]; exists {
		delete(imc.store, key)
		_ = val.Column.Type()
		return 1
	}
	return 0
}
