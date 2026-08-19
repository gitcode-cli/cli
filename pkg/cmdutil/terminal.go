package cmdutil

import "strings"

// SanitizeTerminalText removes control characters from untrusted text before terminal output.
func SanitizeTerminalText(value string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || (r >= 0x7f && r <= 0x9f) {
			return -1
		}
		return r
	}, value)
}
