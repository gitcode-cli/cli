package check

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/precommit"
)

func TestNewCmdCheck(t *testing.T) {
	f := cmdutil.TestFactory()
	cmd := NewCmdCheck(f, nil)
	if cmd == nil {
		t.Fatal("NewCmdCheck returned nil")
	}
	if cmd.Use != "check" {
		t.Errorf("Use = %q, want check", cmd.Use)
	}
}

func TestCheckRunNoConfigSkips(t *testing.T) {
	io, _, out, _ := iostreams.Test()
	r := &stubRunner{}
	opts := &CheckOptions{
		IO:      io,
		GitRoot: func() (string, error) { return t.TempDir(), nil },
		Runner:  r,
	}
	if err := checkRun(opts); err != nil {
		t.Fatalf("checkRun() error = %v", err)
	}
	if !strings.Contains(out.String(), "No pre-commit configuration") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestCheckRunNotReadyReturnsExit1(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir)
	io, _, _, errOut := iostreams.Test()
	r := &stubRunner{}
	opts := &CheckOptions{
		IO:      io,
		GitRoot: func() (string, error) { return dir, nil },
		Runner:  r,
	}
	err := checkRun(opts)
	if err == nil {
		t.Fatal("expected error for not-ready environment")
	}
	if got := cmdutil.ExitCode(err); got != cmdutil.ExitError {
		t.Fatalf("ExitCode = %d, want %d", got, cmdutil.ExitError)
	}
	if !strings.Contains(errOut.String(), "Environment not ready") {
		t.Fatalf("expected stderr hint, got %q", errOut.String())
	}
}

func TestCheckRunJSON(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir)
	writeInstalledHook(t, dir)
	io, _, out, _ := iostreams.Test()
	r := &stubRunner{version: "pre-commit 3.7.0\n"}
	opts := &CheckOptions{
		IO:      io,
		GitRoot: func() (string, error) { return dir, nil },
		Runner:  r,
		JSON:    true,
	}
	if err := checkRun(opts); err != nil {
		t.Fatalf("checkRun() error = %v", err)
	}
	var got precommit.Result
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v; raw=%q", err, out.String())
	}
	if !got.OK || !got.ConfigFound || !got.ToolInstalled || !got.HookInstalled {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestCheckRunNotInGitRepo(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &CheckOptions{
		IO:      io,
		GitRoot: func() (string, error) { return "", os.ErrNotExist },
		Runner:  &stubRunner{},
	}
	err := checkRun(opts)
	if err == nil {
		t.Fatal("expected error when not in a git repo")
	}
}

func TestCheckRunNotInGitRepoJSON(t *testing.T) {
	io, _, out, _ := iostreams.Test()
	opts := &CheckOptions{
		IO:      io,
		GitRoot: func() (string, error) { return "", os.ErrNotExist },
		Runner:  &stubRunner{},
		JSON:    true,
	}
	err := checkRun(opts)
	if err == nil {
		t.Fatal("expected error when not in a git repo")
	}
	var got precommit.Result
	if jerr := json.Unmarshal(out.Bytes(), &got); jerr != nil {
		t.Fatalf("output is not valid JSON: %v; raw=%q", jerr, out.String())
	}
	if got.OK || got.Reason != precommit.ReasonNotInRepo {
		t.Fatalf("want ok=false reason=%q, got %+v", precommit.ReasonNotInRepo, got)
	}
}

func TestCheckRunFlagFailureExits1(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir)
	writeInstalledHook(t, dir)
	io, _, _, errOut := iostreams.Test()
	r := &stubRunner{version: "pre-commit 3.7.0\n", runErr: true}
	opts := &CheckOptions{
		IO:      io,
		GitRoot: func() (string, error) { return dir, nil },
		Runner:  r,
		Run:     true,
	}
	err := checkRun(opts)
	if err == nil {
		t.Fatal("expected error when pre-commit run fails")
	}
	if got := cmdutil.ExitCode(err); got != cmdutil.ExitError {
		t.Fatalf("ExitCode = %d, want %d", got, cmdutil.ExitError)
	}
	// A ready environment whose hooks failed is a check failure, not an
	// unready environment; the messaging must not mislead.
	stderr := errOut.String()
	if !strings.Contains(stderr, "pre-commit checks failed") {
		t.Fatalf("expected 'pre-commit checks failed', got %q", stderr)
	}
	if strings.Contains(stderr, "Environment not ready") {
		t.Fatalf("must not report 'Environment not ready' on a hook failure: %q", stderr)
	}
	if !strings.Contains(stderr, "hook failed") {
		t.Fatalf("expected captured run output in stderr, got %q", stderr)
	}
}

func TestCheckRunFlagFailureJSON(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir)
	writeInstalledHook(t, dir)
	io, _, out, _ := iostreams.Test()
	r := &stubRunner{version: "pre-commit 3.7.0\n", runErr: true}
	opts := &CheckOptions{
		IO:      io,
		GitRoot: func() (string, error) { return dir, nil },
		Runner:  r,
		Run:     true,
		JSON:    true,
	}
	err := checkRun(opts)
	if err == nil {
		t.Fatal("expected error when pre-commit run fails")
	}
	var got precommit.Result
	if jerr := json.Unmarshal(out.Bytes(), &got); jerr != nil {
		t.Fatalf("output is not valid JSON: %v; raw=%q", jerr, out.String())
	}
	if got.OK || got.RunResult != "failed" {
		t.Fatalf("expected ok=false, run_result=failed, got %+v", got)
	}
	if got.RunOutput == "" {
		t.Fatalf("expected run_output to carry failure detail, got %+v", got)
	}
}

func TestCheckMutuallyExclusiveFlags(t *testing.T) {
	f := cmdutil.TestFactory()
	cmd := NewCmdCheck(f, func(*CheckOptions) error { return nil })
	cmd.SetArgs([]string{"--no-install", "--yes"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when --no-install and --yes are combined")
	}
}

func TestCheckRunFlagSuccess(t *testing.T) {
	dir := t.TempDir()
	writeConfig(t, dir)
	writeInstalledHook(t, dir)
	io, _, out, _ := iostreams.Test()
	r := &stubRunner{version: "pre-commit 3.7.0\n"}
	opts := &CheckOptions{
		IO:      io,
		GitRoot: func() (string, error) { return dir, nil },
		Runner:  r,
		Run:     true,
	}
	if err := checkRun(opts); err != nil {
		t.Fatalf("checkRun() error = %v", err)
	}
	if !strings.Contains(out.String(), "pre-commit run passed") {
		t.Fatalf("expected run passed output, got %q", out.String())
	}
}

// --- test helpers ---

func writeConfig(t *testing.T, root string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, ".pre-commit-config.yaml"), []byte("repos: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeInstalledHook(t *testing.T, root string) {
	t.Helper()
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hooks, "pre-commit"), []byte("# File generated by pre-commit\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}

type stubRunner struct {
	version string // output for `pre-commit --version`; empty => error
	runErr  bool   // if true, `pre-commit run --all-files` returns an error
}

func (s *stubRunner) Look(string) bool { return false }

func (s *stubRunner) Run(dir string, name string, args ...string) (string, error) {
	if name == "pre-commit" && len(args) >= 1 && args[0] == "run" {
		if s.runErr {
			return "hook failed", errStub{}
		}
		return "all passed", nil
	}
	return s.RunStdout(dir, name, args...)
}

func (s *stubRunner) RunStdout(_ string, name string, args ...string) (string, error) {
	if name == "pre-commit" && len(args) == 1 && args[0] == "--version" {
		if s.version == "" {
			return "", errStub{}
		}
		return s.version, nil
	}
	return "", nil
}

type errStub struct{}

func (errStub) Error() string { return "not found" }

var _ precommit.CommandRunner = (*stubRunner)(nil)
