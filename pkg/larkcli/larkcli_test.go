package larkcli

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// writeStubBinary writes a tiny executable that echoes its args to stdout,
// writes "stderr-marker" to stderr, and exits with the given code.
func writeStubBinary(t *testing.T, exitCode int) string {
	t.Helper()
	dir := t.TempDir()
	name := "stub-lark"
	if runtime.GOOS == "windows" {
		name = "stub-lark.bat"
	}
	path := filepath.Join(dir, name)
	script := "#!/bin/sh\n" +
		"echo \"args: $*\"\n" +
		"echo stderr-marker 1>&2\n" +
		"exit " + itoa(exitCode) + "\n"
	if runtime.GOOS == "windows" {
		script = "@echo args: %*\r\n@echo stderr-marker 1>&2\r\n@exit " + itoa(exitCode) + "\r\n"
	}
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write stub: %v", err)
	}
	return path
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b []byte
	if i < 0 {
		b = append(b, '-')
		i = -i
	}
	var digs []byte
	for i > 0 {
		digs = append([]byte{byte('0' + i%10)}, digs...)
		i /= 10
	}
	return string(append(b, digs...))
}

func TestEnvConstants(t *testing.T) {
	if EnvBin != "GC_LARK_CLI_BIN" {
		t.Errorf("EnvBin = %q, want GC_LARK_CLI_BIN", EnvBin)
	}
	if EnvDefaultChat != "GC_LARK_DEFAULT_CHAT_ID" {
		t.Errorf("EnvDefaultChat = %q, want GC_LARK_DEFAULT_CHAT_ID", EnvDefaultChat)
	}
}

func TestFindLarkCLI_EnvOverride(t *testing.T) {
	stub := writeStubBinary(t, 0)
	t.Setenv(EnvBin, stub)
	if got := FindLarkCLI(); got != stub {
		t.Errorf("FindLarkCLI() = %q, want %q", got, stub)
	}
}

// clearPath sets PATH to an empty temp dir so exec.LookPath cannot resolve
// lark-cli from the real PATH. Used to construct the "not installed" scenario
// on machines where lark-cli happens to be installed globally.
func clearPath(t *testing.T) {
	t.Helper()
	t.Setenv("PATH", t.TempDir())
}

func TestFindLarkCLI_NotInstalled(t *testing.T) {
	clearPath(t)
	t.Setenv(EnvBin, "")
	if got := FindLarkCLI(); got != "" {
		t.Errorf("FindLarkCLI() = %q, want empty when not installed", got)
	}
}

func TestFindLarkCLI_MissingOverrideFallsBackToPath(t *testing.T) {
	t.Setenv(EnvBin, "/nonexistent/path/does-not-exist")
	// With a missing override, FindLarkCLI falls back to PATH lookup. The
	// result is platform/environment dependent, so we only assert that the
	// function does not return the nonexistent override path.
	got := FindLarkCLI()
	if got == "/nonexistent/path/does-not-exist" {
		t.Errorf("FindLarkCLI() returned the missing override path %q", got)
	}
}

func TestEnsureInstalled(t *testing.T) {
	clearPath(t)
	t.Setenv(EnvBin, "/nonexistent/missing")
	if err := EnsureInstalled(); err != ErrNotInstalled {
		t.Errorf("EnsureInstalled() err = %v, want ErrNotInstalled", err)
	}
	t.Setenv(EnvBin, writeStubBinary(t, 0))
	if err := EnsureInstalled(); err != nil {
		t.Errorf("EnsureInstalled() err = %v, want nil", err)
	}
}

func TestDefaultRun_Success(t *testing.T) {
	stub := writeStubBinary(t, 0)
	t.Setenv(EnvBin, stub)

	res, err := DefaultRun([]string{"im", "+messages-send", "--chat-id", "oc_x"})
	if err != nil {
		t.Fatalf("DefaultRun err = %v", err)
	}
	if res.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", res.ExitCode)
	}
	if !strings.Contains(string(res.Stdout), "args:") {
		t.Errorf("stdout = %q, want to contain args:", res.Stdout)
	}
	if !strings.Contains(string(res.Stderr), "stderr-marker") {
		t.Errorf("stderr = %q, want stderr-marker", res.Stderr)
	}
}

func TestDefaultRun_NonZeroExit(t *testing.T) {
	stub := writeStubBinary(t, 3)
	t.Setenv(EnvBin, stub)

	res, err := DefaultRun([]string{"x"})
	if err != nil {
		t.Fatalf("DefaultRun err = %v (exit-error should be nil)", err)
	}
	if res.ExitCode != 3 {
		t.Errorf("ExitCode = %d, want 3", res.ExitCode)
	}
}

func TestDefaultRun_NotInstalled(t *testing.T) {
	clearPath(t)
	t.Setenv(EnvBin, "")
	if _, err := DefaultRun([]string{"x"}); err != ErrNotInstalled {
		t.Errorf("DefaultRun err = %v, want ErrNotInstalled", err)
	}
}

func TestJSONResult_AppendsJSONFlag(t *testing.T) {
	stub := writeStubBinary(t, 0)
	t.Setenv(EnvBin, stub)

	res, err := JSONResult(func(args []string) (*Result, error) {
		want := "--json"
		if args[len(args)-1] != want {
			t.Errorf("last arg = %q, want %q", args[len(args)-1], want)
		}
		return &Result{Stdout: []byte("{}")}, nil
	}, []string{"auth", "status"})
	if err != nil {
		t.Fatalf("JSONResult err = %v", err)
	}
	if string(res.Stdout) != "{}" {
		t.Errorf("stdout = %q, want {}", res.Stdout)
	}
}

func TestJSONResult_NilRunFallsBackToDefault(t *testing.T) {
	clearPath(t)
	t.Setenv(EnvBin, "")
	if _, err := JSONResult(nil, []string{"x"}); err != ErrNotInstalled {
		t.Errorf("JSONResult(nil) err = %v, want ErrNotInstalled", err)
	}
}
