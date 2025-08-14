package core

import (
	"context"
	"testing"
	"time"

	"github.com/mateenbagheri/memorabilia/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryCommandRepository_Cleanup(t *testing.T) {
	tests := []struct {
		name                string
		store               *InMemoryCommandRepository
		expectedDeleteCount int64
		expectedErr         error
	}{
		{
			name: "one expired key",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key_with_ttl_1": {Column: types.Integer{Val: 1}, Expiration: time.Now().Add(-1 * time.Second)},
					"key_with_ttl_2": {Column: types.Integer{Val: 2}, Expiration: time.Now().Add(1 * time.Hour)},
				},
			),
			expectedDeleteCount: 1,
			expectedErr:         nil,
		},
		{
			name: "two expired keys",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key_with_ttl_1": {Column: types.Integer{Val: 1}, Expiration: time.Now().Add(-1 * time.Second)},
					"key_with_ttl_2": {Column: types.Integer{Val: 2}, Expiration: time.Now().Add(-10 * time.Hour)},
				},
			),
			expectedDeleteCount: 2,
			expectedErr:         nil,
		},
		{
			name: "no expired key",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"key_with_ttl_1": {Column: types.Integer{Val: 1}, Expiration: time.Now().Add(1 * time.Hour)},
					"key_with_ttl_2": {Column: types.Integer{Val: 2}, Expiration: time.Now().Add(1 * time.Hour)},
				},
			),
			expectedDeleteCount: 0,
			expectedErr:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			deleteCount, err := tt.store.Cleanup(ctx)

			assert.Equal(t, tt.expectedDeleteCount, deleteCount, "Cleanup() deleteCount mismatch")
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
