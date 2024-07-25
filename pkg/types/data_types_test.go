package types

import (
	"testing"
)

func TestInteger_ToInt(t *testing.T) {
	i := Integer{Val: 42}
	if got := i.ToInt(); got != 42 {
		t.Errorf("ToInt() = %v, want %v", got, 42)
	}
}

func TestInteger_ToString(t *testing.T) {
	i := Integer{Val: 42}
	if got := i.ToString(); got != "42" {
		t.Errorf("ToString() = %v, want %v", got, "42")
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
