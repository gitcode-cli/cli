package git

import (
	"strings"
	"testing"
)

func TestValidateRef(t *testing.T) {
	tests := []struct {
		name    string
		ref     string
		wantErr bool
		errMsg  string
	}{
		{"valid branch name", "feature/issue-123", false, ""},
		{"valid simple branch", "main", false, ""},
		{"valid ref with dots", "release/v1.0.0", false, ""},
		{"valid ref with underscore", "feature_branch", false, ""},
		{"valid ref with hyphen inside", "feature-branch", false, ""},
		{"empty ref", "", true, "must not be empty"},
		{"dash prefix", "-f", true, "must not start with '-'"},
		{"double dash prefix", "--force", true, "must not start with '-'"},
		{"space in ref", "feature branch", true, "invalid characters"},
		{"newline in ref", "feature\nbranch", true, "invalid characters"},
		{"semicolon in ref", "branch;rm", true, "invalid characters"},
		{"backtick in ref", "branch`id`", true, "invalid characters"},
		{"dollar in ref", "branch$HOME", true, "invalid characters"},
		{"pipe in ref", "branch|cat", true, "invalid characters"},
		{"leading slash", "/refs/heads/main", true, "invalid characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRef(tt.ref)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateRef(%q) expected error, got nil", tt.ref)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateRef(%q) error = %v, want containing %q", tt.ref, err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateRef(%q) unexpected error: %v", tt.ref, err)
				}
			}
		})
	}
}

func TestValidateFetchURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr bool
		errMsg  string
	}{
		{"valid HTTPS URL", "https://gitcode.com/owner/repo.git", false, ""},
		{"valid HTTP URL", "http://gitcode.com/owner/repo.git", false, ""},
		{"valid SSH URL", "git@gitcode.com:owner/repo.git", false, ""},
		{"empty URL", "", true, "must not be empty"},
		{"dash prefix", "-h", true, "must not start with '-'"},
		{"invalid scheme FTP", "ftp://evil.com/repo.git", true, "unsupported URL scheme"},
		{"invalid scheme file", "file:///etc/passwd", true, "unsupported URL scheme"},
		{"SSH with dash host", "git@-evil.com:repo.git", true, "host must not start with '-'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFetchURL(tt.rawURL)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateFetchURL(%q) expected error, got nil", tt.rawURL)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateFetchURL(%q) error = %v, want containing %q", tt.rawURL, err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateFetchURL(%q) unexpected error: %v", tt.rawURL, err)
				}
			}
		})
	}
}

func TestSafeFetchArgs_ValidateRef(t *testing.T) {
	// These tests verify argument validation without executing git.
	// SafeFetch validates args and would fail on the git execution
	// (no git repo in test), but we only care about validation here.

	err := SafeFetch("origin", "-f", "mybranch")
	if err == nil {
		t.Error("SafeFetch with dash-prefixed ref should fail")
	}
	if !strings.Contains(err.Error(), "invalid remote ref") {
		t.Errorf("error = %v, want 'invalid remote ref'", err)
	}

	err = SafeFetch("--upload-pack=evil", "ref", "branch")
	if err == nil {
		t.Error("SafeFetch with dash-prefixed remote should fail")
	}
	if !strings.Contains(err.Error(), "invalid remote name") {
		t.Errorf("error = %v, want 'invalid remote name'", err)
	}

	err = SafeFetch("origin", "feature/x", "--force")
	if err == nil {
		t.Error("SafeFetch with dash-prefixed branch should fail")
	}
	if !strings.Contains(err.Error(), "invalid local branch") {
		t.Errorf("error = %v, want 'invalid local branch'", err)
	}
}

func TestSafeCheckoutArgs_ValidateRef(t *testing.T) {
	err := SafeCheckout("")
	if err == nil {
		t.Error("SafeCheckout with empty branch should fail")
	}

	err = SafeCheckout("--force")
	if err == nil {
		t.Error("SafeCheckout with dash-prefixed branch should fail")
	}

	err = SafeCheckout("branch;rm -rf /")
	if err == nil {
		t.Error("SafeCheckout with semicolon should fail")
	}
}

func TestSafeFetchFromURL_ValidateURL(t *testing.T) {
	err := SafeFetchFromURL("", "ref", "branch")
	if err == nil {
		t.Error("SafeFetchFromURL with empty URL should fail")
	}
	if !strings.Contains(err.Error(), "invalid fetch URL") {
		t.Errorf("error = %v, want 'invalid fetch URL'", err)
	}

	err = SafeFetchFromURL("-h", "ref", "branch")
	if err == nil {
		t.Error("SafeFetchFromURL with dash-prefixed URL should fail")
	}

	err = SafeFetchFromURL("file:///etc/passwd", "ref", "branch")
	if err == nil {
		t.Error("SafeFetchFromURL with file:// scheme should fail")
	}
}

func TestValidateRef_RejectsGitOptions(t *testing.T) {
	// Various git option injection attempts
	injections := []string{
		"--upload-pack=evil",
		"-c core.gitProxy=evil",
		"--config=protocol.ext.allow=always",
		"-u evil",
	}

	for _, inj := range injections {
		t.Run(inj, func(t *testing.T) {
			err := ValidateRef(inj)
			if err == nil {
				t.Errorf("ValidateRef(%q) should reject git option injection", inj)
			}
		})
	}
}

func TestRemoteURLRejectsOptionInjection(t *testing.T) {
	tests := []struct {
		name   string
		remote string
	}{
		{name: "option injection", remote: "--upload-pack=/tmp/evil"},
		{name: "dash prefix", remote: "-bogus"},
		{name: "empty", remote: ""},
		{name: "shell metacharacter", remote: "origin;rm -rf /"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RemoteURL(tt.remote)
			if err == nil {
				t.Fatalf("RemoteURL(%q) error = nil, want rejection", tt.remote)
			}
			if !strings.Contains(err.Error(), "invalid remote name") {
				t.Fatalf("RemoteURL(%q) error = %v, want 'invalid remote name'", tt.remote, err)
			}
		})
	}
}
