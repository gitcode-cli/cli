package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// FlexibleTime is a custom time type that can parse multiple time formats
// returned by GitCode API.
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler interface.
// It tries multiple time formats to parse the input.
func (ft *FlexibleTime) UnmarshalJSON(data []byte) error {
	// Handle null
	if string(data) == "null" {
		ft.Time = time.Time{}
		return nil
	}

	// Remove quotes
	str := strings.Trim(string(data), "\"")
	if str == "" {
		ft.Time = time.Time{}
		return nil
	}

	// List of time formats to try (most common first)
	formats := []string{
		"2006-01-02T15:04:05Z07:00",                // RFC3339 with timezone
		"2006-01-02T15:04:05Z",                      // RFC3339 UTC
		"2006-01-02T15:04:05",                       // ISO 8601 without timezone
		"2006-01-02 15:04:05",                       // Common datetime format
		time.RFC3339,                                 // Standard RFC3339
		"2006-01-02",                                 // Date only
		"2006-01-02T15:04:05.999999999Z07:00",      // RFC3339 with nanoseconds
		"2006-01-02T15:04:05.999999999Z",           // ISO 8601 with nanoseconds UTC
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, str)
		if err == nil {
			ft.Time = t
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("failed to parse time %q: %w", str, lastErr)
}

// MarshalJSON implements json.Marshaler interface.
func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	if ft.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ft.Time.Format(time.RFC3339))
}

// IsZero returns true if the time is zero.
func (ft FlexibleTime) IsZero() bool {
	return ft.Time.IsZero()
}

// Format returns a formatted string representation of the time.
func (ft FlexibleTime) Format(layout string) string {
	return ft.Time.Format(layout)
}

// String returns a string representation of the time.
func (ft FlexibleTime) String() string {
	return ft.Time.String()
}