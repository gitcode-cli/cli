package browser

import (
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
		{"file scheme", "file:///etc/passwd", true},
		{"ftp scheme", "ftp://evil.com", true},
		{"javascript scheme", "javascript:alert(1)", true},
		{"data scheme", "data:text/html,<script>", true},
		{"empty string", "", true},
		{"no scheme", "example.com", true},
		{"ssh scheme", "ssh://evil.com", true},
		{"custom scheme", "custom://handler", true},
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
	if !contains(err.Error(), "unsupported URL scheme") {
		t.Errorf("error message should contain 'unsupported URL scheme', got: %s", err.Error())
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
