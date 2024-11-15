package validation

import (
	"testing"
)

func TestValidateJobTimeFormat(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		// Valid cases
		{"10h", true},    // Valid case with no 'm' and 's'
		{"5m", true},     // Valid case with no 'h' and 's'
		{"12s", true},    // Valid case with no 'h' and 'm'
		{"3h20m", true},  // Valid case with no 's'
		{"2h45s", true},  // Valid case with no 'm'
		{"15m30s", true}, // Valid case with no 'h'
		{"2h45m10s", true},

		// Invalid cases
		{"00h", false},     // Leading zero
		{"01m", false},     // Leading zero
		{"5m50", false},    // Missing 's' after seconds
		{"10h20", false},   // Missing 'm' after minutes
		{"h20m10s", false}, // Missing number before 'h'
		{"2h45s", true},    // Valid case with no 'm'
		{"", false},        // Empty input
	}

	for _, tc := range testCases {
		err := ValidateJobTimeFormat(tc.input)
		if (err == nil) != tc.expected {
			t.Errorf("ValidateTimeFormat(%q): expected %v, got %v", tc.input, tc.expected, err == nil)
		}
	}
}
