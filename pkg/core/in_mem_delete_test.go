package core

import (
	"context"
	"testing"
	"time"

	"github.com/mateenbagheri/memorabilia/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryCommandRepository_Delete(t *testing.T) {
	tests := []struct {
		name                string
		initialStore        *InMemoryCommandRepository
		keyToDelete         string
		expectedDeleteCount int64
		expectedStoreSize   int
		expectKeyToBeAbsent bool
	}{
		{
			name: "Delete existing key",
			initialStore: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key1": {Column: types.String{Val: "value1"}},
					"key2": {Column: types.Integer{Val: 123}},
				},
			),
			keyToDelete:         "key1",
			expectedDeleteCount: 1,
			expectedStoreSize:   1,
			expectKeyToBeAbsent: true,
		},
		{
			name: "Delete non-existing key",
			initialStore: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key1": {Column: types.String{Val: "value1"}},
				},
			),
			keyToDelete:         "nonexistent_key",
			expectedDeleteCount: 0,
			expectedStoreSize:   1,
			expectKeyToBeAbsent: true, // Key should still be absent if it never existed XD
		},
		{
			name: "Delete from empty store",
			initialStore: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{},
			),
			keyToDelete:         "any_key",
			expectedDeleteCount: 0,
			expectedStoreSize:   0,
			expectKeyToBeAbsent: true,
		},
		{
			name: "Delete one of multiple keys",
			initialStore: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"keyA": {Column: types.String{Val: "keyAValue"}},
					"keyB": {Column: types.String{Val: "keyBValue"}},
					"keyC": {Column: types.String{Val: "KeyCValue"}},
				},
			),
			keyToDelete:         "keyB",
			expectedDeleteCount: 1,
			expectedStoreSize:   2,
			expectKeyToBeAbsent: true,
		},
		{
			name: "Delete key with zero TTL (already expired, if logic existed)",
			initialStore: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"expired_key": {
						Column:     types.String{Val: "expired_val"},
						Expiration: time.Time{},
					},
				},
			),
			keyToDelete:         "expired_key",
			expectedDeleteCount: 1,
			expectedStoreSize:   0,
			expectKeyToBeAbsent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			deleteCount := tt.initialStore.Delete(ctx, tt.keyToDelete)

			assert.Equal(t, tt.expectedDeleteCount, deleteCount, "unexpected delete count")

			assert.Len(t, tt.initialStore.store, tt.expectedStoreSize, "unexpected store size after deletion")

			// Assert whether the key exists or is absent as expected
			_, exists := tt.initialStore.store[tt.keyToDelete]
			assert.Equal(t, tt.expectKeyToBeAbsent, !exists, "key presence in store after deletion is incorrect")
		})
	}
}
