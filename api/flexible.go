package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexibleNumber is a JSON number field that can decode both string and int
// representations from GitCode API responses. It always stores the value as
// a string for display consistency.
type FlexibleNumber string

// UnmarshalJSON implements json.Unmarshaler. It accepts both "123" (string)
// and 123 (integer).
func (fn *FlexibleNumber) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*fn = ""
		return nil
	}
	// Try as integer first
	var v int
	if err := json.Unmarshal(data, &v); err == nil {
		*fn = FlexibleNumber(strconv.Itoa(v))
		return nil
	}
	// Fall back to string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("FlexibleNumber: expected string or int, got %s", string(data))
	}
	*fn = FlexibleNumber(s)
	return nil
}

// MarshalJSON implements json.Marshaler. It outputs a string representation.
func (fn FlexibleNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(fn))
}

// String returns the string representation (used by fmt %s and %v).
func (fn FlexibleNumber) String() string {
	return string(fn)
}

// Int returns the integer representation, or an error on parse failure.
func (fn FlexibleNumber) Int() (int, error) {
	if fn == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(string(fn))
	if err != nil {
		return 0, fmt.Errorf("FlexibleNumber.Int: cannot parse %q as int: %w", string(fn), err)
	}
	return n, nil
}

// MustInt returns the integer representation, or 0 on parse failure.
// Prefer Int() for new code; MustInt() is a convenience for callers
// that have already validated the value or tolerate zero-on-error.
func (fn FlexibleNumber) MustInt() int {
	n, err := strconv.Atoi(string(fn))
	if err != nil {
		return 0
	}
	return n
}
