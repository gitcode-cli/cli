// Package larkcli provides a subprocess bridge to the official Feishu/Lark
// CLI tool (lark-cli). GitCode CLI delegates Feishu operations (messaging,
// auth status) to lark-cli instead of reimplementing the Feishu OpenAPI and
// OAuth flow in Go. lark-cli keeps its own credentials in the OS keychain.
package larkcli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// EnvBin overrides the resolved lark-cli binary path. When set and pointing
// to an executable file, it takes precedence over PATH lookup.
const EnvBin = "GC_LARK_CLI_BIN"

// EnvDefaultChat overrides the default Feishu chat id used by notifications.
const EnvDefaultChat = "GC_LARK_DEFAULT_CHAT_ID"

// ErrNotInstalled is returned when lark-cli cannot be found on PATH or via
// GC_LARK_CLI_BIN.
var ErrNotInstalled = errors.New("lark-cli is not installed; run: gc lark install")

// Result holds the captured output of a lark-cli invocation.
type Result struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
}

// RunFunc is the injectable signature for executing lark-cli subcommands.
// Command packages accept a RunFunc so tests can stub out the subprocess.
type RunFunc func(args []string) (*Result, error)

// FindLarkCLI resolves the lark-cli binary path. It returns an empty string
// when no usable binary is found.
func FindLarkCLI() string {
	if p := strings.TrimSpace(os.Getenv(EnvBin)); p != "" {
		if isExecutable(p) {
			return p
		}
	}
	if p, err := exec.LookPath("lark-cli"); err == nil {
		return p
	}
	if p, err := exec.LookPath(larkCLIBinaryName()); err == nil {
		return p
	}
	return ""
}

// larkCLIBinaryName returns the Windows-aware binary name (lark-cli.exe).
func larkCLIBinaryName() string {
	if runtime.GOOS == "windows" {
		return "lark-cli.exe"
	}
	return "lark-cli"
}

func isExecutable(path string) bool {
	if path == "" {
		return false
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	info, err := os.Stat(abs)
	if err != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return info.Mode()&0o111 != 0
}

// DefaultRun is the production RunFunc. It resolves lark-cli via FindLarkCLI
// and executes it with the given args. stdin is NOT inherited (gc passes all
// content as explicit args; lark-cli never needs gc's stdin) and GitCode
// tokens (GC_TOKEN/GITCODE_TOKEN) are stripped from the child environment as
// defense in depth — lark-cli only needs its own keychain credentials.
var DefaultRun RunFunc = func(args []string) (*Result, error) {
	bin := FindLarkCLI()
	if bin == "" {
		return nil, ErrNotInstalled
	}
	return runExec(bin, args)
}

// runExec invokes name with args, capturing stdout/stderr. stdin is left
// closed and GitCode tokens are stripped from the inherited environment.
func runExec(name string, args []string) (*Result, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = nil
	cmd.Env = sanitizedEnv(os.Environ())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("failed to run lark-cli: %w", err)
		}
	}
	return &Result{Stdout: stdout.Bytes(), Stderr: stderr.Bytes(), ExitCode: exitCode}, nil
}

// sanitizedEnv returns env entries with GitCode token values removed so the
// lark-cli subprocess cannot observe GC_TOKEN/GITCODE_TOKEN.
func sanitizedEnv(env []string) []string {
	const tokenPrefix = "GC_TOKEN="
	const tokenPrefix2 = "GITCODE_TOKEN="
	out := make([]string, 0, len(env))
	for _, e := range env {
		if strings.HasPrefix(e, tokenPrefix) || strings.HasPrefix(e, tokenPrefix2) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// JSONResult runs lark-cli with "--json" appended and returns raw stdout.
// Callers parse the JSON payload according to lark-cli's success/error
// envelope contract.
func JSONResult(run RunFunc, args []string) (*Result, error) {
	if run == nil {
		run = DefaultRun
	}
	full := append([]string{}, args...)
	full = append(full, "--json")
	return run(full)
}

// EnsureInstalled returns ErrNotInstalled when lark-cli is not available.
func EnsureInstalled() error {
	if FindLarkCLI() == "" {
		return ErrNotInstalled
	}
	return nil
}
