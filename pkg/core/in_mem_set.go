package core

import (
	"context"
	"time"

	"github.com/mateenbagheri/memorabilia/pkg/types"
)

// Set adds or updates a key-value pair in the in-memory store.
// It ensures thread safety by locking the repository during the operation.
//
// Parameters:
//   - ctx: Context for request-scoped values, cancellation, and deadlines.
//   - key: The key to associate with the value.
//   - value: The value to store.
//   - expiration: The expiration time for the key-value pair. If set to time.Time{},
//     the key-value pair will not expire.
//
// Returns:
//   - error: An error if any issues occur during the operation.
func (imc *InMemoryCommandRepository) Set(
	ctx context.Context,
	key, value string,
	expiration time.Time,
) (err error) {
	imc.mu.Lock()
	defer imc.mu.Unlock()
	_, columnValue := types.DetectColumnType(value)
	imc.store[key] = types.ColumnValueWithTTL{
		Column:     columnValue,
		Expiration: expiration,
	}
	return nil
}
