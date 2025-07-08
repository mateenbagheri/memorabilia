package core

import "context"

// BatchDelete removes multiple key-value pairs from the in-memory store.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - keys: A slice of keys to be deleted from the store.
//
// Returns:
//   - deleteCount: The number of keys that were successfully deleted.
func (imc *InMemoryCommandRepository) BatchDelete(ctx context.Context, keys []string) (deleteCount int64) {
	for _, key := range keys {
		deleteCount += imc.Delete(ctx, key)
	}
	return deleteCount
}
