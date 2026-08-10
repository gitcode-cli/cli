package api

import (
	"strings"
	"testing"
)

func TestSanitizeLogMessage(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no sensitive", "retry: attempt 1 failed", "retry: attempt 1 failed"},
		{"bearer token", "Authorization: Bearer gho_abc123", "Authorization: Bearer [REDACTED]"},
		{"basic auth", "Authorization: Basic dXNlcjpwYXNz", "Authorization: Basic [REDACTED]"},
		{"access_token query", "GET /repos?access_token=gho_abc123", "GET /repos?access_token=[REDACTED]"},
		{"token query", "token=secret123&foo=bar", "token=[REDACTED]&foo=bar"},
		{"token colon", "token: mysecret", "token: [REDACTED]"},
		{"multiple", "Authorization: Bearer gho_x token=y", "Authorization: Bearer [REDACTED] token=[REDACTED]"},
		{"uppercase bearer", "AUTHORIZATION: BEARER gho_abc123", "Authorization: Bearer [REDACTED]"},
		{"lowercase basic", "authorization: basic dXNlcjpwYXNz", "Authorization: Basic [REDACTED]"},
		{"pattern overlap access_token", "access_token=secret123", "access_token=[REDACTED]"},
		{"empty bearer no match", "Authorization: Bearer ", "Authorization: Bearer "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeLogMessage(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeLogMessage(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeLogMessage_RedactsActualTokens(t *testing.T) {
	tokens := []string{"gho_abc123", "secret123", "dXNlcjpwYXNz"}
	for _, token := range tokens {
		msg := "Authorization: Bearer " + token
		got := SanitizeLogMessage(msg)
		if strings.Contains(got, token) {
			t.Errorf("SanitizeLogMessage() = %q, still contains token %q", got, token)
		}
	}
}
