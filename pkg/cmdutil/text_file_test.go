package cmdutil

import (
	"os"
	"path/filepath"
	"strings"
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

func TestDecodeUserTextStripsUTF16LEBOM(t *testing.T) {
	got := DecodeUserText([]byte{0xff, 0xfe, 0x60, 0x4f, 0x7d, 0x59})
	if got != "你好" {
		t.Fatalf("DecodeUserText() = %q, want %q", got, "你好")
	}
}

func TestReadTextUsesDecodeUserText(t *testing.T) {
	got, err := ReadText(strings.NewReader("\xef\xbb\xbfbody"))
	if err != nil {
		t.Fatalf("ReadText() error = %v", err)
	}
	if got != "body" {
		t.Fatalf("ReadText() = %q, want %q", got, "body")
	}
}

func TestDecodeUserTextFallsBackToGB18030(t *testing.T) {
	got := DecodeUserText([]byte{0xc4, 0xe3, 0xba, 0xc3})
	if got != "你好" {
		t.Fatalf("DecodeUserText() = %q, want %q", got, "你好")
	}
}
