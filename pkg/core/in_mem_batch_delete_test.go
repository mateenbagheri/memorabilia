package core

import (
	"context"
	"testing"

	"github.com/mateenbagheri/memorabilia/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryCommandRepository_BatchDelete(t *testing.T) {
	tests := []struct {
		name                string
		store               *InMemoryCommandRepository
		keysToDelete        []string
		expectedDeleteCount int64
		expectedRemaining   map[string]types.ColumnValueWithTTL // Use map for clearer expected state
	}{
		{
			name: "Delete existing keys",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key1": {Column: types.Float{Val: 12.3}},
					"key2": {Column: types.String{Val: "hello"}},
					"key3": {Column: types.Integer{Val: 100}},
				},
			),
			keysToDelete:        []string{"key1", "key3"},
			expectedDeleteCount: 2,
			expectedRemaining:   map[string]types.ColumnValueWithTTL{"key2": {Column: types.String{Val: "hello"}}},
		},
		{
			name: "Delete non-existing keys",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key1": {Column: types.Float{Val: 12.3}},
				},
			),
			keysToDelete:        []string{"nonExistentKey1", "nonExistentKey2"},
			expectedDeleteCount: 0,
			expectedRemaining: map[string]types.ColumnValueWithTTL{
				"key1": {Column: types.Float{Val: 12.3}},
			},
		},
		{
			name: "Delete a mix of existing and non-existing keys",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key1": {Column: types.Float{Val: 12.3}},
					"key2": {Column: types.String{Val: "hello"}},
				},
			),
			keysToDelete:        []string{"key1", "nonExistentKey", "key3"},
			expectedDeleteCount: 1,
			expectedRemaining: map[string]types.ColumnValueWithTTL{
				"key2": {Column: types.String{Val: "hello"}},
			},
		},
		{
			name: "Delete all keys",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key1": {Column: types.Float{Val: 12.3}},
					"key2": {Column: types.String{Val: "hello"}},
				},
			),
			keysToDelete:        []string{"key1", "key2"},
			expectedDeleteCount: 2,
			expectedRemaining:   map[string]types.ColumnValueWithTTL{},
		},
		{
			name: "Delete from empty repository",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{},
			),
			keysToDelete:        []string{"key1", "key2"},
			expectedDeleteCount: 0,
			expectedRemaining:   map[string]types.ColumnValueWithTTL{},
		},
		{
			name: "Empty keys to delete list",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key1": {Column: types.Float{Val: 12.3}},
				},
			),
			keysToDelete:        []string{},
			expectedDeleteCount: 0,
			expectedRemaining: map[string]types.ColumnValueWithTTL{
				"key1": {Column: types.Float{Val: 12.3}},
			},
		},
		{
			name: "Keys with TTLs should be deleted correctly",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key_with_ttl_1": {Column: types.Integer{Val: 1}},
					"key_with_ttl_2": {Column: types.Integer{Val: 2}},
				},
			),
			keysToDelete:        []string{"key_with_ttl_1"},
			expectedDeleteCount: 1,
			expectedRemaining:   map[string]types.ColumnValueWithTTL{"key_with_ttl_2": {Column: types.Integer{Val: 2}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			deleteCount := tt.store.BatchDelete(ctx, tt.keysToDelete)

			assert.Equal(t, tt.expectedDeleteCount, deleteCount, "BatchDelete() deleteCount mismatch")
			// As I checked, Equal function uses reflect.DeepEqual underneath. So we're fine
			assert.Equal(t, tt.expectedRemaining, tt.store.store)
		})
	}
}
