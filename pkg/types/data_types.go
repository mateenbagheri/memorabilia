package types

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrNoneCastable = errors.New("can not cast error")

type ColumnType int

const (
	IntType ColumnType = iota
	StringType
	FloatType
)

type ColumnValue interface {
	Value() any
	Type() ColumnType
	ToInt() int
	ToString() string
}

type Integer struct {
	Val int
}

func (v Integer) Value() any {
	return v.Val
}

func (v Integer) ToInt() int {
	return v.Val
}

func (v Integer) ToString() string {
	return fmt.Sprint(v.Val)
}

func (v Integer) Type() ColumnType {
	return IntType
}

type String struct {
	Val string
}

func (v String) Value() any {
	return v.Val
}

func (v String) ToString() string {
	return v.Val
}

func (v String) ToInt() (int, error) {
	i, err := strconv.Atoi(v.Val)
	if err != nil {
		return 0, ErrNoneCastable
	}
	return i, nil
}

func (v String) Type() ColumnType {
	return StringType
}
