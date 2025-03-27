package core

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("value for given key was not found")
var ErrKeyExpired = errors.New("the given key has expired")

// Get retrieves the value associated with the provided key from the in-memory store.
// It ensures thread safety by locking the repository for reading during the operation.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - key: The key to look up in the store.
//
// Returns:
//   - value: The value associated with the key, as a string.
//   - error: An error if the key is not found or if any other issue occurs.
func (imc *InMemoryCommandRepository) Get(ctx context.Context, key string) (value string, err error) {
	imc.mu.RLock()
	defer imc.mu.RUnlock()
	valueWithTTL, ok := imc.store[key]
	if !ok {
		return "", ErrNotFound
	}

	if !valueWithTTL.Epiration.IsZero() && time.Now().After(valueWithTTL.Epiration) {
		return "", ErrKeyExpired
	}

	return valueWithTTL.Column.ToString(), nil
}
