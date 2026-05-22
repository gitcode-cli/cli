package cmdutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadTextFileStripsUTF8BOM(t *testing.T) {
	path := filepath.Join(t.TempDir(), "body.md")
	if err := os.WriteFile(path, []byte{0xef, 0xbb, 0xbf, 'b', 'o', 'd', 'y'}, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := ReadTextFile(path)
	if err != nil {
		t.Fatalf("ReadTextFile() error = %v", err)
	}
	if got != "body" {
		t.Fatalf("ReadTextFile() = %q, want %q", got, "body")
	}
}

func TestReadTextFileKeepsNonBOMContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "body.md")
	if err := os.WriteFile(path, []byte("plain body"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := ReadTextFile(path)
	if err != nil {
		t.Fatalf("ReadTextFile() error = %v", err)
	}
	if got != "plain body" {
		t.Fatalf("ReadTextFile() = %q, want %q", got, "plain body")
	}
}
