package output

import (
	"fmt"
	"strings"
)

// Format controls high-level output formatting.
type Format string

const (
	FormatJSON   Format = "json"
	FormatSimple Format = "simple"
	FormatTable  Format = "table"
)

// ParseFormat validates a string format value.
func ParseFormat(raw string) (Format, error) {
	switch strings.TrimSpace(raw) {
	case "", string(FormatSimple):
		return FormatSimple, nil
	case string(FormatJSON):
		return FormatJSON, nil
	case string(FormatTable):
		return FormatTable, nil
	default:
		return "", fmt.Errorf("invalid format %q: expected json, simple, or table", raw)
	}
}

// TimeFormat controls how timestamps are rendered in text output.
type TimeFormat string

const (
	TimeFormatAbsolute TimeFormat = "absolute"
	TimeFormatRelative TimeFormat = "relative"
)

// ParseTimeFormat validates a string time-format value.
func ParseTimeFormat(raw string) (TimeFormat, error) {
	switch strings.TrimSpace(raw) {
	case "", string(TimeFormatAbsolute):
		return TimeFormatAbsolute, nil
	case string(TimeFormatRelative):
		return TimeFormatRelative, nil
	default:
		return "", fmt.Errorf("invalid time format %q: expected absolute or relative", raw)
	}
}
