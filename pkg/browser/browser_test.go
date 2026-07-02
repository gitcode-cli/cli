package browser

import (
	"strings"
	"testing"
)

func TestValidateURLScheme(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"http", "http://example.com", false},
		{"https", "https://example.com", false},
		{"uppercase HTTP", "HTTP://example.com", false},
		{"uppercase HTTPS", "HTTPS://example.com", false},
		{"mixed case Https", "Https://example.com", false},
		{"url with userinfo", "http://user:pass@example.com", false},
		{"localhost", "http://127.0.0.1", false},
		{"ipv6", "http://[::1]", false},
		{"file scheme", "file:///etc/passwd", true},
		{"ftp scheme", "ftp://evil.com", true},
		{"javascript scheme", "javascript:alert(1)", true},
		{"data scheme", "data:text/html,<script>", true},
		{"ssh scheme", "ssh://evil.com", true},
		{"custom scheme", "custom://handler", true},
		{"mailto scheme", "mailto:test@example.com", true},
		{"empty string", "", true},
		{"no scheme", "example.com", true},
		{"protocol-relative", "//evil.com", true},
		{"leading space", "  http://example.com", true},
		{"control char in url", "http://example.com\n", true},
		{"null byte", "http://exa\x00mple.com", true},
		{"scheme with space", "ht tp://example.com", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURLScheme(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURLScheme(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestValidateURLSchemeErrMessage(t *testing.T) {
	err := validateURLScheme("file:///etc/passwd")
	if err == nil {
		t.Fatal("expected error for file scheme, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported URL scheme") {
		t.Errorf("error message should contain 'unsupported URL scheme', got: %s", err.Error())
	}
	if !strings.Contains(err.Error(), "file") {
		t.Errorf("error message should contain scheme name 'file', got: %s", err.Error())
	}
}

func TestValidateURLSchemeParseError(t *testing.T) {
	err := validateURLScheme("  http://example.com")
	if err == nil {
		t.Fatal("expected error for leading space, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse URL") {
		t.Errorf("error message should contain 'failed to parse URL', got: %s", err.Error())
	}
}

func TestOpenRejectsBadScheme(t *testing.T) {
	// Verify that Open() enforces scheme validation before reaching exec.Command.
	// This protects against regression if validateURLScheme call is removed.
	err := Open("javascript:alert(1)")
	if err == nil {
		t.Fatal("expected Open to reject javascript scheme, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported URL scheme") {
		t.Errorf("Open should return scheme validation error, got: %s", err.Error())
	}
}
