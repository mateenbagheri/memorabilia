package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// ErrNoneCastable is an error indicating that a value cannot be cast to the desired type.
var ErrNoneCastable = errors.New("cannot cast error")

// ColumnType is an enum type for different column data types.
type ColumnType int

const (
	// IntType represents integer data type.
	IntType ColumnType = iota
	// StringType represents string data type.
	StringType
	// FloatType represents floating-point data type.
	FloatType
)

// ColumnValue is an interface that defines methods for working with column values.
type ColumnValue interface {
	// Value returns the underlying value of the column.
	Value() any
	// Type returns the type of the column as a ColumnType.
	Type() ColumnType
	// ToInt converts the column value to an integer, if possible.
	ToInt() (int, error)
	// ToString converts the column value to a string.
	ToString() string
	// ToFloat converts the column value to a float64
	ToFloat() (float64, error)
}

// ColumnValueWithTTL is a struct designed for storing Column data associated with its expiration.
type ColumnValueWithTTL struct {
	Column     ColumnValue
	Expiration time.Time
}

// columnValueWithTTLJSON is the wire format used when marshaling ColumnValueWithTTL.
// It stores a string type tag alongside the raw value bytes so that UnmarshalJSON
// can reconstruct the correct concrete ColumnValue implementation.
//
// Example JSON produced for Integer{Val: 42} with no expiry:
//
//	{"type":"int","value":{"Val":42},"expiration":"0001-01-01T00:00:00Z"}
type columnValueWithTTLJSON struct {
	Type       string          `json:"type"`
	Value      json.RawMessage `json:"value"`
	Expiration time.Time       `json:"expiration"`
}

// MarshalJSON implements json.Marshaler for ColumnValueWithTTL.
//
// The problem this solves: json.Marshal on an interface field serializes the
// concrete value (e.g. {"Val":42}) but loses the type information. When
// json.Unmarshal tries to reconstruct the interface field, it has no idea
// which concrete struct to allocate — it only knows the target is ColumnValue,
// which is an interface, not a type. The result is the error:
//
//	cannot unmarshal object into Go struct field of type types.ColumnValue
//
// By writing a "type" tag ("int", "string", "float") alongside the value bytes,
// UnmarshalJSON can read the tag first, allocate the right concrete type, then
// unmarshal the value bytes into it.
// I came across this issue while implementing snapshot in FSM for raft replication.
func (c ColumnValueWithTTL) MarshalJSON() ([]byte, error) {
	valueBytes, err := json.Marshal(c.Column)
	if err != nil {
		return nil, fmt.Errorf("ColumnValueWithTTL marshal value: %w", err)
	}

	typeTag, err := columnTypeToTag(c.Column.Type())
	if err != nil {
		return nil, fmt.Errorf("ColumnValueWithTTL marshal type tag: %w", err)
	}

	return json.Marshal(columnValueWithTTLJSON{
		Type:       typeTag,
		Value:      json.RawMessage(valueBytes),
		Expiration: c.Expiration,
	})
}

// UnmarshalJSON implements json.Unmarshaler for ColumnValueWithTTL.
//
// It reads the "type" tag to decide which concrete struct to allocate,
// then unmarshals the "value" bytes into that struct.
func (c *ColumnValueWithTTL) UnmarshalJSON(data []byte) error {
	var envelope columnValueWithTTLJSON
	if err := json.Unmarshal(data, &envelope); err != nil {
		return fmt.Errorf("ColumnValueWithTTL unmarshal envelope: %w", err)
	}

	col, err := unmarshalColumnValue(envelope.Type, envelope.Value)
	if err != nil {
		return err
	}

	c.Column = col
	c.Expiration = envelope.Expiration
	return nil
}

// columnTypeToTag maps a ColumnType to its string tag used in JSON.
func columnTypeToTag(ct ColumnType) (string, error) {
	switch ct {
	case IntType:
		return "int", nil
	case StringType:
		return "string", nil
	case FloatType:
		return "float", nil
	default:
		return "", fmt.Errorf("unknown ColumnType %d", ct)
	}
}

// unmarshalColumnValue allocates the right concrete type for typeTag and
// deserializes valueBytes into it.
func unmarshalColumnValue(typeTag string, valueBytes json.RawMessage) (ColumnValue, error) {
	switch typeTag {
	case "int":
		var v Integer
		if err := json.Unmarshal(valueBytes, &v); err != nil {
			return nil, fmt.Errorf("ColumnValueWithTTL unmarshal int value: %w", err)
		}
		return v, nil
	case "string":
		var v String
		if err := json.Unmarshal(valueBytes, &v); err != nil {
			return nil, fmt.Errorf("ColumnValueWithTTL unmarshal string value: %w", err)
		}
		return v, nil
	case "float":
		var v Float
		if err := json.Unmarshal(valueBytes, &v); err != nil {
			return nil, fmt.Errorf("ColumnValueWithTTL unmarshal float value: %w", err)
		}
		return v, nil
	default:
		return nil, fmt.Errorf("ColumnValueWithTTL unmarshal: unknown type tag %q", typeTag)
	}
}

// ── Concrete implementations (unchanged) ─────────────────────────────────────

// Integer represents a column value of integer type.
type Integer struct {
	Val int
}

func (v Integer) Value() any                { return v.Val }
func (v Integer) ToInt() (int, error)       { return v.Val, nil }
func (v Integer) ToString() string          { return fmt.Sprint(v.Val) }
func (v Integer) ToFloat() (float64, error) { return float64(v.Val), nil }
func (v Integer) Type() ColumnType          { return IntType }

// String represents a column value of string type.
type String struct {
	Val string
}

func (v String) Value() any       { return v.Val }
func (v String) ToString() string { return v.Val }
func (v String) ToInt() (int, error) {
	i, err := strconv.Atoi(v.Val)
	if err != nil {
		return 0, ErrNoneCastable
	}
	return i, nil
}

func (v String) ToFloat() (float64, error) {
	f, err := strconv.ParseFloat(v.Val, 64)
	if err != nil {
		return 0.0, ErrNoneCastable
	}
	return f, nil
}
func (v String) Type() ColumnType { return StringType }

// Float represents a column value of float type.
type Float struct {
	Val float64
}

func (v Float) Value() any                { return v.Val }
func (v Float) ToInt() (int, error)       { return int(v.Val), nil }
func (v Float) ToString() string          { return fmt.Sprint(v.Val) }
func (v Float) ToFloat() (float64, error) { return v.Val, nil }
func (v Float) Type() ColumnType          { return FloatType }

// DetectColumnType takes a string input and determines its appropriate ColumnType.
func DetectColumnType(input string) (ColumnType, ColumnValue) {
	if i, err := strconv.Atoi(input); err == nil {
		return IntType, Integer{Val: i}
	}
	if f, err := strconv.ParseFloat(input, 64); err == nil {
		return FloatType, Float{Val: f}
	}
	return StringType, String{Val: input}
}
