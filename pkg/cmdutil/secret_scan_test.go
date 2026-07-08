package cmdutil

import (
	"strings"
	"testing"
)

func TestScanContentForSecrets_Clean(t *testing.T) {
	t.Setenv("GC_TOKEN", "")
	t.Setenv("GITCODE_TOKEN", "")
	cases := []string{
		"normal issue body without secrets",
		"PR description with code: func main() { fmt.Println(\"hello\") }",
		"",
		"中文内容，无敏感信息",
	}
	for _, c := range cases {
		if err := ScanContentForSecrets(c); err != nil {
			t.Errorf("ScanContentForSecrets(%q) unexpected error: %v", c, err)
		}
	}
}

func TestScanContentForSecrets_GC_TOKEN(t *testing.T) {
	t.Setenv("GC_TOKEN", "my-secret-gc-token-123")
	t.Setenv("GITCODE_TOKEN", "")
	if err := ScanContentForSecrets("body with my-secret-gc-token-123 inside"); err == nil {
		t.Fatal("expected error for GC_TOKEN value in content, got nil")
	}
	if err := ScanContentForSecrets("clean body"); err != nil {
		t.Errorf("unexpected error for clean body: %v", err)
	}
}

func TestScanContentForSecrets_GITCODE_TOKEN(t *testing.T) {
	t.Setenv("GC_TOKEN", "")
	t.Setenv("GITCODE_TOKEN", "gitcode-token-value-xyz")
	err := ScanContentForSecrets("leaked: gitcode-token-value-xyz")
	if err == nil {
		t.Fatal("expected error for GITCODE_TOKEN value in content, got nil")
	}
	if !strings.Contains(err.Error(), "GITCODE_TOKEN") {
		t.Errorf("error = %v, want containing 'GITCODE_TOKEN'", err)
	}
}

func TestScanContentForSecrets_NoTokenSet(t *testing.T) {
	t.Setenv("GC_TOKEN", "")
	t.Setenv("GITCODE_TOKEN", "")
	if err := ScanContentForSecrets("any content"); err != nil {
		t.Errorf("unexpected error when no token set: %v", err)
	}
}
