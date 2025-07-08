package core

import (
	"context"
	"testing"
	"time"

	"github.com/mateenbagheri/memorabilia/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryCommandRepository_GetExpiredKeys(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		store        *InMemoryCommandRepository
		expectedKeys []string
	}{
		{
			name: "expired and non-expired keys",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"expired1":     {Column: types.String{Val: "a"}, Expiration: now.Add(-2 * time.Minute)},
					"expired2":     {Column: types.String{Val: "b"}, Expiration: now.Add(-1 * time.Second)},
					"notExpired":   {Column: types.String{Val: "c"}, Expiration: now.Add(5 * time.Minute)},
					"noExpiration": {Column: types.String{Val: "d"}, Expiration: time.Time{}},
				},
			),
			expectedKeys: []string{"expired1", "expired2"},
		},
		{
			name: "no expired keys",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"active1": {Column: types.String{Val: "x"}, Expiration: now.Add(1 * time.Hour)},
					"active2": {Column: types.String{Val: "y"}, Expiration: time.Time{}},
				},
			),
			expectedKeys: []string{},
		},
		{
			name: "all keys expired",
			store: NewInMemoryCommandRepositoryWithInitialStore(
				map[string]types.ColumnValueWithTTL{
					"expiredA": {Column: types.String{Val: "1"}, Expiration: now.Add(-10 * time.Minute)},
					"expiredB": {Column: types.String{Val: "2"}, Expiration: now.Add(-1 * time.Hour)},
				},
			),
			expectedKeys: []string{"expiredA", "expiredB"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			keys, err := tt.store.GetExpiredKeys(ctx)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.expectedKeys, keys)
		})
	}
}
