package core

import (
	"context"
	"testing"

	"github.com/mateenbagheri/memorabilia/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryCommandRepository_Get_ExistingInteger(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValue),
	}

	key1 := "intKey"
	value1 := "123"
	imc.Set(ctx, key1, value1)

	result, err := imc.Get(ctx, key1)
	assert.NoError(t, err, "Get should not return an error for existing key")
	assert.Equal(t, "123", result, "The returned value should match the stored value as a string")
}

func TestInMemoryCommandRepository_Get_ExistingFloat(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValue),
	}

	key2 := "floatKey"
	value2 := "123.45"
	imc.Set(ctx, key2, value2)

	result, err := imc.Get(ctx, key2)
	assert.NoError(t, err, "Get should not return an error for existing key")
	assert.Equal(t, "123.45", result, "The returned value should match the stored value as a string")
}

func TestInMemoryCommandRepository_Get_ExistingString(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValue),
	}

	key3 := "stringKey"
	value3 := "hello"
	imc.Set(ctx, key3, value3)

	result, err := imc.Get(ctx, key3)
	assert.NoError(t, err, "Get should not return an error for existing key")
	assert.Equal(t, "hello", result, "The returned value should match the stored value as a string")
}

func TestInMemoryCommandRepository_Get_NonExistantValue_ErrNotFound(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValue),
	}

	missingKey := "missingKey"
	result, err := imc.Get(ctx, missingKey)
	assert.Error(t, err, "Get should return an error for a non-existing key")
	assert.Equal(t, ErrNotFound, err, "The error should be ErrNotFound")
	assert.Equal(t, "", result, "The returned value should be an empty string when the key is not found")
}
