package core

import (
	"context"

	"github.com/mateenbagheri/memorabilia/pkg/types"
)

// Set adds or updates a key-value pair in the in-memory store.
// It locks the repository to ensure thread safety while performing the operation.
// The function takes a context, a key, and a value as input, and returns an error if any issues arise.
func (imc *InMemoryCommandRepository) Set(ctx context.Context, key, value string) (err error) {
	imc.mu.Lock()
	defer imc.mu.Unlock()
	_, columnValue := types.DetectColumnType(value)
	imc.store[key] = columnValue
	return nil
}
