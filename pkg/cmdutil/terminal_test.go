package cmdutil

import "testing"

func TestSanitizeTerminalText(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "plain", value: "work laptop", want: "work laptop"},
		{name: "C0 controls", value: "work\n\t\x1b[31m", want: "work[31m"},
		{name: "C1 controls", value: "work\u0085key", want: "workkey"},
		{name: "unicode text", value: "开发密钥", want: "开发密钥"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeTerminalText(tt.value); got != tt.want {
				t.Fatalf("SanitizeTerminalText() = %q, want %q", got, tt.want)
			}
		})
	}
}
