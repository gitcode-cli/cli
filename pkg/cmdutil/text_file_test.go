package cmdutil

import (
	"errors"
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

func TestReadTextDoesNotApplyLossyPowerShellGuard(t *testing.T) {
	got, err := ReadText(strings.NewReader("Windows????-20260526"))
	if err != nil {
		t.Fatalf("ReadText() error = %v", err)
	}
	if got != "Windows????-20260526" {
		t.Fatalf("ReadText() = %q", got)
	}
}

func TestDecodeUserTextFallsBackToGB18030(t *testing.T) {
	got := DecodeUserText([]byte{0xc4, 0xe3, 0xba, 0xc3})
	if got != "你好" {
		t.Fatalf("DecodeUserText() = %q, want %q", got, "你好")
	}
}

func TestIsLikelyLossyPowerShellStdin(t *testing.T) {
	if !isLikelyLossyPowerShellStdin([]byte("Windows????-20260526"), "Windows????-20260526", "windows", false) {
		t.Fatal("expected lossy PowerShell stdin to be detected")
	}
}

func TestIsLikelyLossyPowerShellStdinCanBeBypassed(t *testing.T) {
	if isLikelyLossyPowerShellStdin([]byte("Windows????-20260526"), "Windows????-20260526", "windows", true) {
		t.Fatal("expected explicit lossy stdin bypass to be honored")
	}
}

func TestIsLikelyLossyPowerShellStdinIgnoresUTF8Text(t *testing.T) {
	raw := []byte("Windows中文正文")
	if isLikelyLossyPowerShellStdin(raw, string(raw), "windows", false) {
		t.Fatal("expected valid UTF-8 stdin text to be accepted")
	}
}

func TestIsLikelyLossyPowerShellStdinNonWindows(t *testing.T) {
	if isLikelyLossyPowerShellStdin([]byte("Linux????-20260526"), "Linux????-20260526", "linux", false) {
		t.Fatal("expected non-Windows stdin to be accepted")
	}
}

func TestNewLossyPowerShellStdinErrorIncludesFlagAndExamples(t *testing.T) {
	err := newLossyPowerShellStdinError("--comment-file")
	if !errors.Is(err, ErrLossyPowerShellStdin) {
		t.Fatalf("error does not wrap ErrLossyPowerShellStdin: %v", err)
	}
	text := err.Error()
	for _, want := range []string{
		"--comment-file -",
		"Set-Content -Path comment.md",
		"gitcode pr review 1 -R owner/repo --comment-file",
		"GITCODE_CLI_ALLOW_LOSSY_STDIN=1",
		"正确用法",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("error %q does not contain %q", text, want)
		}
	}
}
