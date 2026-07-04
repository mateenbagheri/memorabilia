package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mateenbagheri/memorabilia/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// roundTrip marshals v to JSON then unmarshals back into a new ColumnValueWithTTL.
func roundTrip(t *testing.T, v types.ColumnValueWithTTL) types.ColumnValueWithTTL {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err, "marshal failed")

	var result types.ColumnValueWithTTL
	require.NoError(t, json.Unmarshal(b, &result), "unmarshal failed")
	return result
}

func TestColumnValueWithTTL_JSON_Integer(t *testing.T) {
	original := types.ColumnValueWithTTL{
		Column:     types.Integer{Val: 42},
		Expiration: time.Time{},
	}
	got := roundTrip(t, original)

	assert.Equal(t, types.IntType, got.Column.Type())
	assert.Equal(t, "42", got.Column.ToString())
	assert.True(t, got.Expiration.IsZero())
}

func TestColumnValueWithTTL_JSON_String(t *testing.T) {
	original := types.ColumnValueWithTTL{
		Column:     types.String{Val: "hello world"},
		Expiration: time.Time{},
	}
	got := roundTrip(t, original)

	assert.Equal(t, types.StringType, got.Column.Type())
	assert.Equal(t, "hello world", got.Column.ToString())
}

func TestColumnValueWithTTL_JSON_Float(t *testing.T) {
	original := types.ColumnValueWithTTL{
		Column:     types.Float{Val: 3.14},
		Expiration: time.Time{},
	}
	got := roundTrip(t, original)

	assert.Equal(t, types.FloatType, got.Column.Type())
	f, err := got.Column.ToFloat()
	require.NoError(t, err)
	assert.InDelta(t, 3.14, f, 0.0001)
}

func TestColumnValueWithTTL_JSON_WithExpiration(t *testing.T) {
	expiry := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	original := types.ColumnValueWithTTL{
		Column:     types.String{Val: "expires"},
		Expiration: expiry,
	}
	got := roundTrip(t, original)

	assert.Equal(t, "expires", got.Column.ToString())
	assert.True(t, got.Expiration.Equal(expiry),
		"expiration mismatch: want %v got %v", expiry, got.Expiration)
}

func TestColumnValueWithTTL_JSON_MapRoundTrip(t *testing.T) {
	// This is the exact scenario the Raft snapshot uses:
	// marshal a whole map, unmarshal it back, assert every entry survives.
	store := map[string]types.ColumnValueWithTTL{
		"int_key":   {Column: types.Integer{Val: 99}},
		"str_key":   {Column: types.String{Val: "memorabilia"}},
		"float_key": {Column: types.Float{Val: 2.718}},
	}

	b, err := json.Marshal(store)
	require.NoError(t, err)

	var restored map[string]types.ColumnValueWithTTL
	require.NoError(t, json.Unmarshal(b, &restored))

	require.Len(t, restored, 3)
	assert.Equal(t, "99", restored["int_key"].Column.ToString())
	assert.Equal(t, "memorabilia", restored["str_key"].Column.ToString())
	f, _ := restored["float_key"].Column.ToFloat()
	assert.InDelta(t, 2.718, f, 0.0001)
}
