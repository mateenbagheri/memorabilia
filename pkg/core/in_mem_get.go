package core

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("value for given key was not found")

// Get retrieves the value associated with the provided key from the in-memory store.
// If the key is found, it returns the value as a string. If the key is not found,
// it returns an error indicating that the key does not exist.
func (imc *InMemoryCommandRepository) Get(ctx context.Context, key string) (value string, err error) {
	imc.mu.RLock()
	defer imc.mu.RUnlock()
	val, ok := imc.store[key]
	if !ok {
		return "", ErrNotFound
	}

	return val.ToString(), nil
}
