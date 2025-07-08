package core

import (
	"context"
	"testing"
	"time"

	"github.com/mateenbagheri/memorabilia/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryCommandRepository_Get_ExistingInteger(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValueWithTTL),
	}

	key1 := "intKey"
	value1 := "123"
	imc.Set(ctx, key1, value1, time.Time{})

	result, err := imc.Get(ctx, key1)
	assert.NoError(t, err, "Get should not return an error for existing key")
	assert.Equal(t, "123", result, "The returned value should match the stored value as a string")
}

func TestInMemoryCommandRepository_Get_ExistingFloat(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValueWithTTL),
	}

	key2 := "floatKey"
	value2 := "123.45"
	imc.Set(ctx, key2, value2, time.Time{})

	result, err := imc.Get(ctx, key2)
	assert.NoError(t, err, "Get should not return an error for existing key")
	assert.Equal(t, "123.45", result, "The returned value should match the stored value as a string")
}

func TestInMemoryCommandRepository_Get_ExistingString(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValueWithTTL),
	}

	key3 := "stringKey"
	value3 := "hello"
	imc.Set(ctx, key3, value3, time.Time{})

	result, err := imc.Get(ctx, key3)
	assert.NoError(t, err, "Get should not return an error for existing key")
	assert.Equal(t, "hello", result, "The returned value should match the stored value as a string")
}

func TestInMemoryCommandRepository_Get_NonExistantValue_ErrKeyNotFoundForDeleteOp(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValueWithTTL),
	}

	missingKey := "missingKey"
	result, err := imc.Get(ctx, missingKey)
	assert.Error(t, err, "Get should return an error for a non-existing key")
	assert.Equal(t, ErrNotFoundForGetOp, err, "The error should be ErrNotFoundForGetOp")
	assert.Equal(t, "", result, "The returned value should be an empty string when the key is not found")
}

func TestInMemoryCommandRepository_Get_UnexpiredValue(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory repository instance
	imc := &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValueWithTTL),
	}

	key := "expiration_unexpired_test"
	value := "hello-kitty"
	expiration := time.Now().Add(3000 * time.Millisecond) // Set expiration to 3 seconds in the future

	// Set the key-value pair with an expiration time
	err := imc.Set(ctx, key, value, expiration)
	if err != nil {
		t.Fatalf("Failed to set key-value pair: %v", err)
	}

	// Test case 1: Retrieve the value immediately after setting it (no wait)
	t.Run("Retrieve immediately", func(t *testing.T) {
		retrievedValue, err := imc.Get(ctx, key)
		if err != nil {
			t.Fatalf("Failed to get key: %v", err)
		}

		// Verify that the retrieved value matches the original value
		if retrievedValue != value {
			t.Errorf("Expected value %q, got %q", value, retrievedValue)
		}
	})

	// Test case 2: Retrieve the value after waiting for 2 seconds (less than the expiration time)
	t.Run("Retrieve before expiration", func(t *testing.T) {
		waitChan := time.After(2000 * time.Millisecond) // Wait for 2 seconds
		<-waitChan

		retrievedValue, err := imc.Get(ctx, key)
		if err != nil {
			t.Fatalf("Failed to get key: %v", err)
		}

		// Verify that the retrieved value matches the original value
		if retrievedValue != value {
			t.Errorf("Expected value %q, got %q", value, retrievedValue)
		}
	})

	// Test case 3: Retrieve the value after waiting for 4 seconds (longer than the expiration time)
	t.Run("Retrieve after expiration", func(t *testing.T) {
		waitChan := time.After(4000 * time.Millisecond) // Wait for 4 seconds
		<-waitChan

		retrievedValue, err := imc.Get(ctx, key)
		if err == nil || err != ErrKeyExpiredForGetOp {
			t.Fatalf("Expected ErrKeyExpired, got: %v", err)
		}

		// Verify that no value is returned
		if retrievedValue != "" {
			t.Errorf("Expected empty value, got %q", retrievedValue)
		}
	})
}
