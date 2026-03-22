package iostreams

import (
	"fmt"
	"strings"
)

// ColorScheme provides color output utilities
type ColorScheme struct {
	noColor bool
}

// Color codes
const (
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	bold    = "\033[1m"
	reset   = "\033[0m"
)

// Green returns the text in green
func (c *ColorScheme) Green(text string) string {
	if c.noColor {
		return text
	}
	return fmt.Sprintf("%s%s%s", green, text, reset)
}

// Red returns the text in red
func (c *ColorScheme) Red(text string) string {
	if c.noColor {
		return text
	}
	return fmt.Sprintf("%s%s%s", red, text, reset)
}

// Yellow returns the text in yellow
func (c *ColorScheme) Yellow(text string) string {
	if c.noColor {
		return text
	}
	return fmt.Sprintf("%s%s%s", yellow, text, reset)
}

// Blue returns the text in blue
func (c *ColorScheme) Blue(text string) string {
	if c.noColor {
		return text
	}
	return fmt.Sprintf("%s%s%s", blue, text, reset)
}

// Cyan returns the text in cyan
func (c *ColorScheme) Cyan(text string) string {
	if c.noColor {
		return text
	}
	return fmt.Sprintf("%s%s%s", cyan, text, reset)
}

// Magenta returns the text in magenta
func (c *ColorScheme) Magenta(text string) string {
	if c.noColor {
		return text
	}
	return fmt.Sprintf("%s%s%s", magenta, text, reset)
}

// Bold returns the text in bold
func (c *ColorScheme) Bold(text string) string {
	if c.noColor {
		return text
	}
	return fmt.Sprintf("%s%s%s", bold, text, reset)
}

// SuccessIcon returns a green checkmark
func (c *ColorScheme) SuccessIcon() string {
	if c.noColor {
		return "✓"
	}
	return fmt.Sprintf("%s✓%s", green, reset)
}

// FailureIcon returns a red X
func (c *ColorScheme) FailureIcon() string {
	if c.noColor {
		return "✗"
	}
	return fmt.Sprintf("%s✗%s", red, reset)
}

// WarningIcon returns a yellow warning sign
func (c *ColorScheme) WarningIcon() string {
	if c.noColor {
		return "!"
	}
	return fmt.Sprintf("%s!%s", yellow, reset)
}

// ColorFromString returns a color function based on a string name
func (c *ColorScheme) ColorFromString(colorName string) func(string) string {
	switch strings.ToLower(colorName) {
	case "red":
		return c.Red
	case "green":
		return c.Green
	case "yellow":
		return c.Yellow
	case "blue":
		return c.Blue
	case "cyan":
		return c.Cyan
	case "magenta":
		return c.Magenta
	case "bold":
		return c.Bold
	default:
		return func(s string) string { return s }
	}
}