package cmdutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadInputFile(t *testing.T) {
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origCwd) })

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := os.WriteFile("body.md", []byte("hello"), 0o644); err != nil {
		t.Fatalf("write body.md: %v", err)
	}
	if err := os.Mkdir("docs", 0o755); err != nil {
		t.Fatalf("mkdir docs: %v", err)
	}
	if err := os.WriteFile("docs/body.md", []byte("nested"), 0o644); err != nil {
		t.Fatalf("write docs/body.md: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{"valid file in cwd", "body.md", false, ""},
		{"valid file in subdir", "docs/body.md", false, ""},
		{"valid absolute path in cwd", filepath.Join(tmpDir, "body.md"), false, ""},
		{"empty path", "", true, "must not be empty"},
		{"nonexistent file", "missing.md", true, ""},
		{"directory", "docs", true, "directory"},
		{"traversal outside cwd", "../../../etc/passwd", true, "within the current directory"},
		{"absolute path outside cwd", "/etc/passwd", true, "within the current directory"},
		{"traversal to home", "../../.bashrc", true, "within the current directory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ReadInputFile(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ReadInputFile(%q) expected error, got nil (data=%q)", tt.path, string(data))
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ReadInputFile(%q) error = %v, want containing %q", tt.path, err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ReadInputFile(%q) unexpected error: %v", tt.path, err)
				}
			}
		})
	}
}

func TestReadInputFileRejectsOversize(t *testing.T) {
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origCwd) })

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	big := make([]byte, MaxInputFileSize+1)
	if err := os.WriteFile("big.bin", big, 0o644); err != nil {
		t.Fatalf("write big.bin: %v", err)
	}

	_, err = ReadInputFile("big.bin")
	if err == nil {
		t.Fatal("expected error for oversize file, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds size limit") {
		t.Errorf("error = %v, want containing 'exceeds size limit'", err)
	}
}

func TestReadInputFileAcceptsExactLimit(t *testing.T) {
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origCwd) })

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	exact := make([]byte, MaxInputFileSize)
	if err := os.WriteFile("exact.bin", exact, 0o644); err != nil {
		t.Fatalf("write exact.bin: %v", err)
	}

	data, err := ReadInputFile("exact.bin")
	if err != nil {
		t.Fatalf("ReadInputFile(exact.bin) unexpected error: %v", err)
	}
	if int64(len(data)) != MaxInputFileSize {
		t.Errorf("data length = %d, want %d", len(data), MaxInputFileSize)
	}
}

func TestReadInputFileReturnsContent(t *testing.T) {
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origCwd) })

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	want := "request body content"
	if err := os.WriteFile("input.txt", []byte(want), 0o644); err != nil {
		t.Fatalf("write input.txt: %v", err)
	}

	data, err := ReadInputFile("input.txt")
	if err != nil {
		t.Fatalf("ReadInputFile unexpected error: %v", err)
	}
	if string(data) != want {
		t.Errorf("data = %q, want %q", string(data), want)
	}
}
