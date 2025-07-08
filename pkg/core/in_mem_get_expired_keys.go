package core

import (
	"context"
	"time"
)

// GetExpiredKeys returns a list of keys whose expiration time has passed.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//
// Returns:
//   - keys: A slice of expired keys.
//   - err: An error if any issue occurs (always nil in current implementation).
func (imc *InMemoryCommandRepository) GetExpiredKeys(ctx context.Context) (keys []string, err error) {
	imc.mu.RLock()
	defer imc.mu.RUnlock()

	now := time.Now()
	for key, val := range imc.store {
		if !val.Expiration.IsZero() && now.After(val.Expiration) {
			keys = append(keys, key)
		}
	}
	return keys, nil
}
