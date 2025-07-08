package types

import (
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

// ColumnValueWithTTL is a struct designed for storing Column data assocoated with its expiration.
type ColumnValueWithTTL struct {
	Column     ColumnValue
	Expiration time.Time
}

// Integer represents a column value of integer type.
type Integer struct {
	Val int
}

// Value returns the integer value.
func (v Integer) Value() any {
	return v.Val
}

// ToInt returns the integer value directly.
func (v Integer) ToInt() (int, error) {
	return v.Val, nil
}

// ToString converts the integer value to a string.
func (v Integer) ToString() string {
	return fmt.Sprint(v.Val)
}

// ToFloat converts an integer to a float
func (v Integer) ToFloat() (float64, error) {
	return float64(v.Val), nil
}

// Type returns the ColumnType for integer, which is IntType.
func (v Integer) Type() ColumnType {
	return IntType
}

// String represents a column value of string type.
type String struct {
	Val string
}

// Value returns the string value.
func (v String) Value() any {
	return v.Val
}

// ToString returns the string value directly.
func (v String) ToString() string {
	return v.Val
}

// ToInt attempts to convert the string value to an integer.
// Returns the integer and nil if successful, otherwise returns 0 and ErrNoneCastable.
func (v String) ToInt() (int, error) {
	i, err := strconv.Atoi(v.Val)
	if err != nil {
		return 0, ErrNoneCastable
	}
	return i, nil
}

// ToFloat attempts to convert the string value to a float64.
// Returns the float64 and nil if successful, otherwise returns 0.0 and ErrNonCastable.
func (v String) ToFloat() (float64, error) {
	f, err := strconv.ParseFloat(v.Val, 64)
	if err != nil {
		return 0.0, ErrNoneCastable
	}
	return f, nil
}

// Type returns the ColumnType for string, which is StringType.
func (v String) Type() ColumnType {
	return StringType
}

// Float represents a column value of float type.
type Float struct {
	Val float64
}

// Value returns the float value.
func (v Float) Value() any {
	return v.Val
}

// ToInt converts the float value to an integer.
func (v Float) ToInt() (int, error) {
	return int(v.Val), nil
}

// ToString converts the float value to a string.
func (v Float) ToString() string {
	return fmt.Sprint(v.Val)
}

func (v Float) ToFloat() (float64, error) {
	return v.Val, nil
}

// Type returns the ColumnType for float, which is FloatType.
func (v Float) Type() ColumnType {
	return FloatType
}

// DetectColumnType takes a string input and determines its appropriate ColumnType.
func DetectColumnType(input string) (ColumnType, ColumnValue) {
	// Attempt to parse as an integer
	if i, err := strconv.Atoi(input); err == nil {
		return IntType, Integer{Val: i}
	}

	// Attempt to parse as a float
	if f, err := strconv.ParseFloat(input, 64); err == nil {
		return FloatType, Float{Val: f}
	}

	// If neither, it's considered a string type
	return StringType, String{Val: input}
}
