package core

import (
	"context"
	"testing"

	"github.com/mateenbagheri/memorabilia/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryCommandRepository_Set_Integer(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValue),
	}

	// Adding a new key with an integer value
	key1 := "intKey"
	value1 := "123"

	err := imc.Set(ctx, key1, value1)
	assert.NoError(t, err, "Set should not return an error")

	// Check if the key-value pair was correctly stored and its type is integer
	storedValue1, exists := imc.store[key1]
	assert.True(t, exists, "The key should exist in the store")
	assert.Equal(t, types.IntType, storedValue1.Type(), "The stored value should be of IntType")
	assert.Equal(t, 123, storedValue1.Value(), "The stored value should match the input value as an integer")
}

func TestInMemoryCommandRepository_Set_Float(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValue),
	}

	// Adding a new key with a float value
	key2 := "floatKey"
	value2 := "123.45"

	err := imc.Set(ctx, key2, value2)
	assert.NoError(t, err, "Set should not return an error")

	// Check if the key-value pair was correctly stored and its type is float
	storedValue2, exists := imc.store[key2]
	assert.True(t, exists, "The key should exist in the store")
	assert.Equal(t, types.FloatType, storedValue2.Type(), "The stored value should be of FloatType")
	assert.Equal(t, 123.45, storedValue2.Value(), "The stored value should match the input value as a float")

}

func TestInMemoryCommandRepository_Set_String(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValue),
	}

	// Adding a new key with a string value
	key3 := "stringKey"
	value3 := "hello"

	err := imc.Set(ctx, key3, value3)
	assert.NoError(t, err, "Set should not return an error")

	// Check if the key-value pair was correctly stored and its type is string
	storedValue3, exists := imc.store[key3]
	assert.True(t, exists, "The key should exist in the store")
	assert.Equal(t, types.StringType, storedValue3.Type(), "The stored value should be of StringType")
	assert.Equal(t, "hello", storedValue3.Value(), "The stored value should match the input value as a string")
}

func TestInMemoryCommandRepository_Set_ExistingValue(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValue),
	}

	// Updating an existing key with a new value
	newValue := "456"

	key1 := "intKey"
	value1 := "123"

	err := imc.Set(ctx, key1, value1)
	assert.NoError(t, err, "Set should not return an error")

	err = imc.Set(ctx, key1, newValue)
	assert.NoError(t, err, "Set should not return an error")

	// Check if the key's value was updated and its type is integer
	storedValue1Updated, exists := imc.store[key1]
	assert.True(t, exists, "The key should exist in the store")

	assert.Equal(t, types.IntType, storedValue1Updated.Type(),
		"The updated value should be of IntType",
	)
	assert.Equal(t, 456, storedValue1Updated.Value(),
		"The updated value should match the new input value as an integer",
	)
}
