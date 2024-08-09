package types

import (
	"testing"
)

func TestInteger_ToInt(t *testing.T) {
	i := Integer{Val: 42}
	if got, err := i.ToInt(); err != nil || got != 42 {
		if err != nil {
			t.Errorf("Expected no errors but got following error = %v", err)
		}
		t.Errorf("ToInt() = %v, want %v", got, 42)
	}
}

func TestInteger_ToString(t *testing.T) {
	i := Integer{Val: 42}
	if got := i.ToString(); got != "42" {
		t.Errorf("ToString() = %v, want %v", got, "42")
	}
}

func TestInteger_ToFloat(t *testing.T) {
	f := Float{Val: 40.24}
	got, err := f.ToFloat()
	if err != nil {
		t.Errorf("ToFloat() expected no errors, got %v", err)
	}

	if got != 40.24 {
		t.Errorf("ToFloat() = %v, want %v", got, 42.2)
	}
}

func TestString_ToString(t *testing.T) {
	s := String{Val: "hello"}
	if got := s.ToString(); got != "hello" {
		t.Errorf("ToString() = %v, want %v", got, "hello")
	}
}

func TestString_ToInt(t *testing.T) {
	tests := []struct {
		name    string
		s       String
		want    int
		wantErr bool
	}{
		{"valid int", String{Val: "123"}, 123, false},
		{"invalid int", String{Val: "abc"}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ToInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_ToFloat(t *testing.T) {
	tests := []struct {
		name    string
		s       String
		want    float64
		wantErr bool
	}{
		{"valid integer value", String{Val: "42"}, 42.0, false},
		{"valid floating value", String{Val: "23.2"}, 23.2, false},
		{"valid zero value", String{Val: "0"}, 0.0, false},
		{"invalid value", String{Val: "test"}, 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ToFloat()
			if err != nil && !tt.wantErr {
				t.Errorf("ToFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Type(t *testing.T) {
	s := String{Val: "test"}
	if got := s.Type(); got != StringType {
		t.Errorf("Type() = %v, want %v", got, StringType)
	}
}

func TestInteger_Type(t *testing.T) {
	i := Integer{Val: 42}
	if got := i.Type(); got != IntType {
		t.Errorf("Type() = %v, want %v", got, IntType)
	}
}

func TestFloat_Type(t *testing.T) {
	f := Float{Val: 12.5}
	if got := f.Type(); got != FloatType {
		t.Errorf("Type() = %v, want %v", got, FloatType)
	}
}

func TestInteger_Value(t *testing.T) {
	i := Integer{Val: 42}
	if got := i.Value(); got != 42 {
		t.Errorf("Value() = %v, want %v", got, 42)
	}
}

func TestString_Value(t *testing.T) {
	s := String{Val: "hello"}
	if got := s.Value(); got != "hello" {
		t.Errorf("Value() = %v, want %v", got, "hello")
	}
}

func TestFloat_Value(t *testing.T) {
	f := Float{Val: 23.6}
	if got := f.Value(); got != 23.6 {
		t.Errorf("Value() = %v, want %v", got, 23.6)
	}
}

func TestDetectColumnType(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedType  ColumnType
		expectedValue ColumnValue
	}{
		{"Integer Value", "123", IntType, Integer{Val: 123}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columnType, value := DetectColumnType(tt.input)
			if columnType != tt.expectedType {
				t.Errorf("DetectColumnType(%v) = %v, want %v",
					tt.input, columnType, tt.expectedType,
				)
			}

			if value != tt.expectedValue {
				t.Errorf("DetectColumnType(%v)::Value = %v, want %v",
					tt.input, value, tt.expectedValue,
				)
			}
		})
	}
}
