# `gc precommit check` Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `gc precommit check` command that detects pre-commit config, verifies/auto-installs the pre-commit tool, ensures the git hook is initialized, and optionally runs the checks — cross-platform.

**Architecture:** A pure-logic package `pkg/precommit/` does detection/install/orchestration behind an injectable `CommandRunner` interface (so tests never touch real `pre-commit`/`python`). A thin command layer `pkg/cmd/precommit/` wires it to the Cobra `Factory` + `IOStreams`, following the existing `auth/status` pattern. Mutating the environment is gated: allowed in a TTY, or with `--yes` in non-TTY, and never with `--no-install`.

**Tech Stack:** Go, Cobra, `gitcode.com/gitcode-cli/cli` internal packages (`pkg/cmdutil`, `pkg/iostreams`, `git`).

---

## File Structure

- Create `pkg/precommit/runner.go` — `CommandRunner` interface + `execRunner` (wraps `os/exec`).
- Create `pkg/precommit/detect.go` — locate config file, query tool version, check hook installation.
- Create `pkg/precommit/install.go` — ensure the tool is installed; install the git hook.
- Create `pkg/precommit/check.go` — orchestrate the pipeline into a `Result`.
- Create `pkg/precommit/fake_runner_test.go` — shared test fake for `CommandRunner`.
- Create `pkg/precommit/*_test.go` — unit tests per file above.
- Create `pkg/cmd/precommit/precommit.go` — `precommit` command group.
- Create `pkg/cmd/precommit/check/check.go` — `check` subcommand.
- Create `pkg/cmd/precommit/check/check_test.go` — command-layer tests.
- Modify `pkg/cmd/root/root.go` — register the `precommit` command.
- Modify `docs/COMMANDS.md` and the gitcode-cli skill cheat sheets — doc sync.

---

## Task 1: CommandRunner abstraction

**Files:**
- Create: `pkg/precommit/runner.go`
- Create: `pkg/precommit/fake_runner_test.go`

- [ ] **Step 1: Write `runner.go`**

```go
// Package precommit detects and manages pre-commit configuration and environment.
package precommit

import "os/exec"

// CommandRunner runs external commands. It is abstracted so tests can avoid
// invoking real pre-commit / python binaries.
type CommandRunner interface {
	// Look reports whether an executable named name exists on PATH.
	Look(name string) bool
	// Run executes name with args. If dir != "", it is the working directory.
	// It returns the combined stdout+stderr output and any execution error.
	Run(dir, name string, args ...string) (string, error)
}

type execRunner struct{}

// NewExecRunner returns a CommandRunner backed by os/exec.
func NewExecRunner() CommandRunner { return execRunner{} }

func (execRunner) Look(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func (execRunner) Run(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}
```

- [ ] **Step 2: Write the shared test fake `fake_runner_test.go`**

```go
package precommit

import "strings"

// fakeRunner is a deterministic CommandRunner for tests.
type fakeRunner struct {
	look      map[string]bool   // executables present on PATH
	responses map[string]fakeResp // keyed by "name arg1 arg2"
	calls     []string          // recorded "name arg1 arg2" invocations
}

type fakeResp struct {
	out string
	err error
}

func newFakeRunner() *fakeRunner {
	return &fakeRunner{
		look:      map[string]bool{},
		responses: map[string]fakeResp{},
	}
}

func key(name string, args ...string) string {
	if len(args) == 0 {
		return name
	}
	return name + " " + strings.Join(args, " ")
}

func (f *fakeRunner) Look(name string) bool { return f.look[name] }

func (f *fakeRunner) Run(_ string, name string, args ...string) (string, error) {
	k := key(name, args...)
	f.calls = append(f.calls, k)
	r := f.responses[k]
	return r.out, r.err
}

func (f *fakeRunner) called(name string, args ...string) bool {
	want := key(name, args...)
	for _, c := range f.calls {
		if c == want {
			return true
		}
	}
	return false
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./pkg/precommit/...`
Expected: builds with no error (no tests yet).

- [ ] **Step 4: Commit**

```bash
git add pkg/precommit/runner.go pkg/precommit/fake_runner_test.go
git commit -m "feat(precommit): add CommandRunner abstraction and test fake"
```

---

## Task 2: Detection (config file, tool version, hook installed)

**Files:**
- Create: `pkg/precommit/detect.go`
- Test: `pkg/precommit/detect_test.go`

- [ ] **Step 1: Write the failing tests `detect_test.go`**

```go
package precommit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigFile(t *testing.T) {
	dir := t.TempDir()
	if _, ok := ConfigFile(dir); ok {
		t.Fatal("expected no config in empty dir")
	}

	yaml := filepath.Join(dir, ".pre-commit-config.yaml")
	if err := os.WriteFile(yaml, []byte("repos: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, ok := ConfigFile(dir)
	if !ok || got != yaml {
		t.Fatalf("ConfigFile() = %q, %v; want %q, true", got, ok, yaml)
	}
}

func TestConfigFileYmlExtension(t *testing.T) {
	dir := t.TempDir()
	yml := filepath.Join(dir, ".pre-commit-config.yml")
	if err := os.WriteFile(yml, []byte("repos: []\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, ok := ConfigFile(dir)
	if !ok || got != yml {
		t.Fatalf("ConfigFile() = %q, %v; want %q, true", got, ok, yml)
	}
}

func TestToolVersion(t *testing.T) {
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	v, ok := ToolVersion(r)
	if !ok || v != "3.7.0" {
		t.Fatalf("ToolVersion() = %q, %v; want 3.7.0, true", v, ok)
	}
}

func TestToolVersionMissing(t *testing.T) {
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{err: os.ErrNotExist}
	if v, ok := ToolVersion(r); ok || v != "" {
		t.Fatalf("ToolVersion() = %q, %v; want \"\", false", v, ok)
	}
}

func TestHookInstalled(t *testing.T) {
	root := t.TempDir()
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	if HookInstalled(root) {
		t.Fatal("expected hook not installed")
	}
	hook := filepath.Join(hooks, "pre-commit")
	if err := os.WriteFile(hook, []byte("#!/bin/sh\n# File generated by pre-commit: https://pre-commit.com\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !HookInstalled(root) {
		t.Fatal("expected hook installed")
	}
}

func TestHookInstalledIgnoresForeignHook(t *testing.T) {
	root := t.TempDir()
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	hook := filepath.Join(hooks, "pre-commit")
	if err := os.WriteFile(hook, []byte("#!/bin/sh\necho custom\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if HookInstalled(root) {
		t.Fatal("expected foreign hook to be treated as not pre-commit-installed")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./pkg/precommit/ -run 'Config|ToolVersion|Hook' -v`
Expected: FAIL — `undefined: ConfigFile`, `ToolVersion`, `HookInstalled`.

- [ ] **Step 3: Write `detect.go`**

```go
package precommit

import (
	"os"
	"path/filepath"
	"strings"
)

// configNames lists accepted pre-commit config filenames, in priority order.
var configNames = []string{".pre-commit-config.yaml", ".pre-commit-config.yml"}

// ConfigFile returns the path to the pre-commit config in root, and whether one exists.
func ConfigFile(root string) (string, bool) {
	for _, name := range configNames {
		p := filepath.Join(root, name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p, true
		}
	}
	return "", false
}

// ToolVersion returns the installed pre-commit version (e.g. "3.7.0") and whether
// the tool is available.
func ToolVersion(r CommandRunner) (string, bool) {
	out, err := r.Run("", "pre-commit", "--version")
	if err != nil {
		return "", false
	}
	// Output looks like "pre-commit 3.7.0".
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) >= 2 {
		return fields[len(fields)-1], true
	}
	if len(fields) == 1 {
		return fields[0], true
	}
	return "", true
}

// HookInstalled reports whether root/.git/hooks/pre-commit exists and was generated
// by the pre-commit tool (detected via a signature marker).
func HookInstalled(root string) bool {
	hook := filepath.Join(root, ".git", "hooks", "pre-commit")
	data, err := os.ReadFile(hook)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "pre-commit")
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./pkg/precommit/ -run 'Config|ToolVersion|Hook' -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add pkg/precommit/detect.go pkg/precommit/detect_test.go
git commit -m "feat(precommit): detect config file, tool version, hook installation"
```

---

## Task 3: Install (ensure tool, install hook)

**Files:**
- Create: `pkg/precommit/install.go`
- Test: `pkg/precommit/install_test.go`

- [ ] **Step 1: Write the failing tests `install_test.go`**

```go
package precommit

import (
	"errors"
	"testing"
)

func TestEnsureToolAlreadyPresent(t *testing.T) {
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	action, err := EnsureTool(r)
	if err != nil {
		t.Fatalf("EnsureTool() error = %v", err)
	}
	if action != "" {
		t.Fatalf("action = %q, want empty (no install needed)", action)
	}
}

func TestEnsureToolViaPipx(t *testing.T) {
	r := newFakeRunner()
	// Missing first, present after install.
	calls := 0
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errors.New("not found")}
	r.look["pipx"] = true
	r.responses[key("pipx", "install", "pre-commit")] = fakeResp{out: "installed"}

	// Make the version check succeed after pipx install by swapping the response
	// on the second call via a small wrapper.
	wrapped := &versionAfterInstall{fakeRunner: r, succeedAfter: 1, callCount: &calls}
	action, err := EnsureTool(wrapped)
	if err != nil {
		t.Fatalf("EnsureTool() error = %v", err)
	}
	if action == "" {
		t.Fatal("expected a non-empty action describing the install")
	}
	if !r.called("pipx", "install", "pre-commit") {
		t.Fatal("expected pipx install pre-commit to be called")
	}
}

func TestEnsureToolNoPython(t *testing.T) {
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errors.New("not found")}
	// no pipx, python3, or python on PATH
	_, err := EnsureTool(r)
	if err == nil {
		t.Fatal("expected error when no installer is available")
	}
}

func TestInstallHook(t *testing.T) {
	r := newFakeRunner()
	r.responses[key("pre-commit", "install")] = fakeResp{out: "pre-commit installed at .git/hooks/pre-commit"}
	action, err := InstallHook(r, "/repo")
	if err != nil {
		t.Fatalf("InstallHook() error = %v", err)
	}
	if action == "" {
		t.Fatal("expected non-empty action")
	}
	if !r.called("pre-commit", "install") {
		t.Fatal("expected pre-commit install to be called")
	}
}

// versionAfterInstall makes pre-commit --version fail until succeedAfter calls have
// happened, simulating a tool that appears after installation.
type versionAfterInstall struct {
	*fakeRunner
	succeedAfter int
	callCount    *int
}

func (w *versionAfterInstall) Run(dir, name string, args ...string) (string, error) {
	if name == "pre-commit" && len(args) == 1 && args[0] == "--version" {
		*w.callCount++
		if *w.callCount > w.succeedAfter {
			return "pre-commit 3.7.0\n", nil
		}
		return "", errors.New("not found")
	}
	return w.fakeRunner.Run(dir, name, args...)
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./pkg/precommit/ -run 'EnsureTool|InstallHook' -v`
Expected: FAIL — `undefined: EnsureTool`, `InstallHook`.

- [ ] **Step 3: Write `install.go`**

```go
package precommit

import (
	"fmt"
	"runtime"
)

// EnsureTool makes sure the pre-commit tool is available, installing it when
// necessary. It returns a human-readable description of any action taken (empty
// if the tool was already present) and an error if installation was impossible
// or failed.
func EnsureTool(r CommandRunner) (string, error) {
	if _, ok := ToolVersion(r); ok {
		return "", nil
	}

	var action string
	switch {
	case r.Look("pipx"):
		if out, err := r.Run("", "pipx", "install", "pre-commit"); err != nil {
			return "", fmt.Errorf("pipx install pre-commit failed: %w: %s", err, out)
		}
		action = "installed pre-commit via pipx"
	case r.Look("python3"):
		if out, err := r.Run("", "python3", "-m", "pip", "install", "--user", "pre-commit"); err != nil {
			return "", fmt.Errorf("python3 -m pip install pre-commit failed: %w: %s", err, out)
		}
		action = "installed pre-commit via python3 -m pip --user"
	case r.Look("python"):
		if out, err := r.Run("", "python", "-m", "pip", "install", "--user", "pre-commit"); err != nil {
			return "", fmt.Errorf("python -m pip install pre-commit failed: %w: %s", err, out)
		}
		action = "installed pre-commit via python -m pip --user"
	default:
		return "", fmt.Errorf("cannot auto-install pre-commit: no pipx/python3/python found. %s", manualInstallHint())
	}

	if _, ok := ToolVersion(r); !ok {
		return "", fmt.Errorf("pre-commit still not found after install attempt; ensure the install location is on PATH. %s", manualInstallHint())
	}
	return action, nil
}

// InstallHook runs `pre-commit install` in root to set up the git hook.
func InstallHook(r CommandRunner, root string) (string, error) {
	if out, err := r.Run(root, "pre-commit", "install"); err != nil {
		return "", fmt.Errorf("pre-commit install failed: %w: %s", err, out)
	}
	return "ran pre-commit install", nil
}

func manualInstallHint() string {
	switch runtime.GOOS {
	case "darwin":
		return "Install manually, e.g.: brew install pre-commit (or: pipx install pre-commit)."
	case "windows":
		return "Install manually, e.g.: python -m pip install --user pre-commit (or: pipx install pre-commit)."
	default:
		return "Install manually, e.g.: pipx install pre-commit (or: pip install --user pre-commit; apt/yum may also provide it)."
	}
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./pkg/precommit/ -run 'EnsureTool|InstallHook' -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add pkg/precommit/install.go pkg/precommit/install_test.go
git commit -m "feat(precommit): ensure tool installed and install git hook"
```

---

## Task 4: Check orchestration

**Files:**
- Create: `pkg/precommit/check.go`
- Test: `pkg/precommit/check_test.go`

- [ ] **Step 1: Write the failing tests `check_test.go`**

```go
package precommit

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestCheckNoConfigSkips(t *testing.T) {
	root := t.TempDir()
	r := newFakeRunner()
	res, err := Check(r, Options{Root: root})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if res.ConfigFound {
		t.Fatal("ConfigFound should be false")
	}
	if !res.OK {
		t.Fatal("no config should be OK (skipped)")
	}
}

func TestCheckMissingToolNoInstall(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	r := newFakeRunner() // pre-commit --version returns err by default
	res, err := Check(r, Options{Root: root, AllowInstall: false})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if res.ToolInstalled {
		t.Fatal("ToolInstalled should be false")
	}
	if res.OK {
		t.Fatal("missing tool without install should not be OK")
	}
	if len(res.ActionsTaken) != 0 {
		t.Fatalf("no actions should be taken, got %v", res.ActionsTaken)
	}
}

func TestCheckAllReady(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	writeInstalledHook(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	res, err := Check(r, Options{Root: root})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !res.OK || !res.ConfigFound || !res.ToolInstalled || !res.HookInstalled {
		t.Fatalf("expected fully ready result, got %+v", res)
	}
	if res.ToolVersion != "3.7.0" {
		t.Fatalf("ToolVersion = %q", res.ToolVersion)
	}
}

func TestCheckInstallsHookWhenAllowed(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	// Simulate `pre-commit install` writing the hook so the post-install re-check passes.
	r.responses[key("pre-commit", "install")] = fakeResp{out: "ok"}
	hookWriter := &hookInstallingRunner{fakeRunner: r, root: root}

	res, err := Check(hookWriter, Options{Root: root, AllowInstall: true})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !res.HookInstalled || !res.OK {
		t.Fatalf("expected hook installed + OK, got %+v", res)
	}
	if len(res.ActionsTaken) == 0 {
		t.Fatal("expected an action recorded for hook install")
	}
}

func TestCheckRunFails(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	writeInstalledHook(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	r.responses[key("pre-commit", "run", "--all-files")] = fakeResp{out: "black....Failed", err: errExit{}}

	res, err := Check(r, Options{Root: root, Run: true})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if res.RunResult != "failed" {
		t.Fatalf("RunResult = %q, want failed", res.RunResult)
	}
	if res.OK {
		t.Fatal("failed run should make result not OK")
	}
}

type errExit struct{}

func (errExit) Error() string { return "exit status 1" }

// hookInstallingRunner writes a pre-commit hook file when `pre-commit install` runs,
// simulating the real tool's side effect.
type hookInstallingRunner struct {
	*fakeRunner
	root string
}

func (h *hookInstallingRunner) Run(dir, name string, args ...string) (string, error) {
	out, err := h.fakeRunner.Run(dir, name, args...)
	if name == "pre-commit" && len(args) == 1 && args[0] == "install" {
		writeInstalledHookNoT(h.root)
	}
	return out, err
}

func writeInstalledHookNoT(root string) {
	hooks := filepath.Join(root, ".git", "hooks")
	_ = os.MkdirAll(hooks, 0o755)
	_ = os.WriteFile(filepath.Join(hooks, "pre-commit"), []byte("# File generated by pre-commit\n"), 0o755)
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./pkg/precommit/ -run TestCheck -v`
Expected: FAIL — `undefined: Check`, `Options`, `Result`.

- [ ] **Step 3: Write `check.go`**

```go
package precommit

// Options controls a Check run.
type Options struct {
	// Root is the git repository root directory.
	Root string
	// AllowInstall permits mutating the environment (installing the tool and/or
	// the git hook) when something is missing.
	AllowInstall bool
	// Run, when true, executes `pre-commit run --all-files` after the environment
	// is confirmed ready.
	Run bool
}

// Result is the structured outcome of a Check. JSON tags match docs/COMMANDS.md.
type Result struct {
	ConfigFound   bool     `json:"config_found"`
	ToolInstalled bool     `json:"tool_installed"`
	ToolVersion   string   `json:"tool_version,omitempty"`
	HookInstalled bool     `json:"hook_installed"`
	ActionsTaken  []string `json:"actions_taken"`
	RunResult     string   `json:"run_result,omitempty"` // "passed" | "failed" | ""
	OK            bool     `json:"ok"`
}

// Check runs the detection/remediation pipeline and returns a structured Result.
// It returns a non-nil error only for hard failures (e.g. an install attempt
// failed). "Not ready" states are reported via Result.OK == false with no error.
func Check(r CommandRunner, opts Options) (Result, error) {
	res := Result{ActionsTaken: []string{}}

	// 1. Config detection — absence is a clean skip, not an error.
	if _, found := ConfigFile(opts.Root); !found {
		res.OK = true
		return res, nil
	}
	res.ConfigFound = true

	// 2. Tool detection + optional install.
	version, ok := ToolVersion(r)
	if !ok && opts.AllowInstall {
		action, err := EnsureTool(r)
		if err != nil {
			return res, err
		}
		if action != "" {
			res.ActionsTaken = append(res.ActionsTaken, action)
		}
		version, ok = ToolVersion(r)
	}
	res.ToolInstalled = ok
	res.ToolVersion = version
	if !ok {
		return res, nil // not ready; OK stays false
	}

	// 3. Hook detection + optional install.
	hookOK := HookInstalled(opts.Root)
	if !hookOK && opts.AllowInstall {
		action, err := InstallHook(r, opts.Root)
		if err != nil {
			return res, err
		}
		if action != "" {
			res.ActionsTaken = append(res.ActionsTaken, action)
		}
		hookOK = HookInstalled(opts.Root)
	}
	res.HookInstalled = hookOK
	if !hookOK {
		return res, nil
	}

	// 4. Environment is ready.
	res.OK = true

	// 5. Optional: actually run the checks.
	if opts.Run {
		if _, err := r.Run(opts.Root, "pre-commit", "run", "--all-files"); err != nil {
			res.RunResult = "failed"
			res.OK = false
		} else {
			res.RunResult = "passed"
		}
	}

	return res, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./pkg/precommit/ -v`
Expected: PASS (all precommit package tests).

- [ ] **Step 5: Commit**

```bash
git add pkg/precommit/check.go pkg/precommit/check_test.go
git commit -m "feat(precommit): orchestrate detection/install/run pipeline"
```

---

## Task 5: `precommit check` command

**Files:**
- Create: `pkg/cmd/precommit/check/check.go`
- Test: `pkg/cmd/precommit/check/check_test.go`

- [ ] **Step 1: Write the failing tests `check_test.go`**

```go
package check

import (
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
	r := &stubRunner{} // pre-commit --version errors by default
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
	io, _, _, _ := iostreams.Test()
	r := &stubRunner{} // tool missing
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
	if !strings.Contains(out.String(), `"ok": true`) && !strings.Contains(out.String(), `"ok":true`) {
		t.Fatalf("json output = %q", out.String())
	}
}

// --- test helpers ---

type stubRunner struct {
	version string // output for `pre-commit --version`; empty => error
}

func (s *stubRunner) Look(string) bool { return false }

func (s *stubRunner) Run(_ string, name string, args ...string) (string, error) {
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
```

Also add the `writeConfig` / `writeInstalledHook` helpers to this test file (copy the bodies from `pkg/precommit/check_test.go` Step 1, adjusted to this package):

```go
import (
	"os"
	"path/filepath"
)

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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./pkg/cmd/precommit/check/ -v`
Expected: FAIL — `undefined: NewCmdCheck`, `CheckOptions`, `checkRun`.

- [ ] **Step 3: Write `check.go`**

```go
// Package check implements the `precommit check` command.
package check

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	gitpkg "gitcode.com/gitcode-cli/cli/git"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
	"gitcode.com/gitcode-cli/cli/pkg/precommit"
)

// CheckOptions holds dependencies and flags for the check command.
type CheckOptions struct {
	IO      *iostreams.IOStreams
	GitRoot func() (string, error)
	Runner  precommit.CommandRunner

	Run       bool
	NoInstall bool
	Yes       bool
	JSON      bool
}

// NewCmdCheck creates the `precommit check` command.
func NewCmdCheck(f *cmdutil.Factory, runF func(*CheckOptions) error) *cobra.Command {
	opts := &CheckOptions{
		IO:      f.IOStreams,
		GitRoot: gitpkg.RootDir,
		Runner:  precommit.NewExecRunner(),
	}

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check pre-commit configuration and local environment before committing",
		Long: heredoc.Doc(`
			Check whether the current repository configures pre-commit and whether the
			local environment is ready to run it before committing code.

			The command:
			  1. Detects a .pre-commit-config.yaml (or .yml) in the repository root.
			  2. Verifies the pre-commit tool is installed.
			  3. Verifies the git pre-commit hook is initialized.
			  4. Optionally runs the hooks with --run.

			When something is missing it auto-installs/initializes in an interactive
			terminal. In a non-interactive (non-TTY) environment, pass --yes to allow
			environment changes, or --no-install to only diagnose. Cross-platform
			(Windows, Linux x86/arm, macOS).
		`),
		Example: heredoc.Doc(`
			# Verify the environment is ready
			$ gc precommit check

			# Verify and actually run the hooks
			$ gc precommit check --run

			# Only diagnose, never modify the environment
			$ gc precommit check --no-install

			# Machine-readable output
			$ gc precommit check --json
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return checkRun(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.Run, "run", false, "Run pre-commit hooks (pre-commit run --all-files) after verifying")
	cmd.Flags().BoolVar(&opts.NoInstall, "no-install", false, "Only diagnose; never install the tool or hook")
	cmd.Flags().BoolVarP(&opts.Yes, "yes", "y", false, "Allow environment changes (install/init) in non-interactive mode")
	cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output result as JSON")

	return cmd
}

func checkRun(opts *CheckOptions) error {
	root, err := opts.GitRoot()
	if err != nil {
		return cmdutil.NewCLIError(cmdutil.ExitError, "not in a git repository", err)
	}

	allowInstall := !opts.NoInstall && (opts.IO.IsStdoutTTY() || opts.Yes)

	res, err := precommit.Check(opts.Runner, precommit.Options{
		Root:         root,
		AllowInstall: allowInstall,
		Run:          opts.Run,
	})
	if err != nil {
		return cmdutil.NewCLIError(cmdutil.ExitError, "pre-commit check failed", err)
	}

	if opts.JSON {
		if writeErr := cmdutil.WriteJSON(opts.IO.Out, res); writeErr != nil {
			return writeErr
		}
	} else {
		printResult(opts, res, allowInstall)
	}

	if !res.OK {
		return cmdutil.NewCLIError(cmdutil.ExitError, "pre-commit environment is not ready", nil)
	}
	return nil
}

func printResult(opts *CheckOptions, res precommit.Result, allowInstall bool) {
	cs := opts.IO.ColorScheme()
	out := opts.IO.Out

	if !res.ConfigFound {
		fmt.Fprintf(out, "%s No pre-commit configuration found; nothing to check.\n", cs.Green("✓"))
		return
	}

	mark := func(ok bool) string {
		if ok {
			return cs.Green("✓")
		}
		return cs.Red("✗")
	}

	fmt.Fprintf(out, "%s pre-commit configuration found\n", mark(true))
	if res.ToolInstalled {
		fmt.Fprintf(out, "%s pre-commit tool installed (%s)\n", mark(true), res.ToolVersion)
	} else {
		fmt.Fprintf(out, "%s pre-commit tool not installed\n", mark(false))
	}
	fmt.Fprintf(out, "%s git hook initialized\n", mark(res.HookInstalled))

	for _, a := range res.ActionsTaken {
		fmt.Fprintf(out, "  - %s\n", a)
	}

	switch res.RunResult {
	case "passed":
		fmt.Fprintf(out, "%s pre-commit run passed\n", mark(true))
	case "failed":
		fmt.Fprintf(out, "%s pre-commit run failed\n", mark(false))
	}

	if !res.OK {
		fmt.Fprintf(opts.IO.ErrOut, "\nEnvironment not ready.\n")
		if !allowInstall && (!res.ToolInstalled || !res.HookInstalled) {
			fmt.Fprintf(opts.IO.ErrOut, "Re-run in a terminal, or pass --yes to auto-install/initialize.\n")
		}
		if !res.ToolInstalled {
			fmt.Fprintf(opts.IO.ErrOut, "Install pre-commit, e.g.: pipx install pre-commit (or pip install --user pre-commit).\n")
		} else if !res.HookInstalled {
			fmt.Fprintf(opts.IO.ErrOut, "Initialize hooks: pre-commit install\n")
		}
	}
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./pkg/cmd/precommit/check/ -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add pkg/cmd/precommit/check/check.go pkg/cmd/precommit/check/check_test.go
git commit -m "feat(precommit): add precommit check command"
```

---

## Task 6: Command group + root registration

**Files:**
- Create: `pkg/cmd/precommit/precommit.go`
- Modify: `pkg/cmd/root/root.go` (imports block + `AddCommand` list around lines 13-26 and 50-60)

- [ ] **Step 1: Write `precommit.go`**

```go
// Package precommit implements the precommit command group.
package precommit

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/pkg/cmd/precommit/check"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

// NewCmdPrecommit creates the precommit command group.
func NewCmdPrecommit(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "precommit <command>",
		Short: "Work with pre-commit configuration and hooks",
		Long: heredoc.Doc(`
			Detect and verify pre-commit configuration and the local pre-commit
			environment before committing code.
		`),
		Example: heredoc.Doc(`
			# Verify the environment is ready
			$ gc precommit check

			# Verify and run the hooks
			$ gc precommit check --run
		`),
		Annotations: map[string]string{
			cmdutil.TopicAnnotation: "precommit",
		},
	}

	cmd.AddCommand(check.NewCmdCheck(f, nil))

	return cmd
}
```

> Note: confirm the constant name by grepping `TopicAnnotation` in `pkg/cmdutil/`. If it does not exist, drop the `Annotations` field entirely (it is non-essential). `commit.go` uses `cmdutil.TopicAnnotation`, so it should exist.

- [ ] **Step 2: Register in `root.go`**

Add to the import block (after the `pr` import, keeping alphabetical-ish grouping):

```go
	precommitcmd "gitcode.com/gitcode-cli/cli/pkg/cmd/precommit"
```

Add to the `AddCommand` list (after the `commitcmd` line):

```go
	cmd.AddCommand(precommitcmd.NewCmdPrecommit(f))
```

- [ ] **Step 3: Build and verify the command is wired**

Run: `go build -o ./gc.exe ./cmd/gc; ./gc.exe precommit check --help`
Expected: build succeeds; help text for `precommit check` prints with `--run`, `--no-install`, `--yes`, `--json` flags.

- [ ] **Step 4: Run the full test suite**

Run: `go test ./...`
Expected: PASS (all packages).

- [ ] **Step 5: Commit**

```bash
git add pkg/cmd/precommit/precommit.go pkg/cmd/root/root.go
git commit -m "feat(precommit): register precommit command group on root"
```

---

## Task 7: Documentation sync

**Files:**
- Modify: `docs/COMMANDS.md` (add a `precommit check` entry following the existing command-doc style)
- Modify: `.ai/skills/gitcode-cli/SKILL.md` and `.claude/skills/gitcode-cli/SKILL.md` cheat-sheet tables (add a `gc precommit check` row)

- [ ] **Step 1: Read the current docs to match style**

Run: `Read docs/COMMANDS.md` (find an existing command section, e.g. `commit`, to mirror heading depth and table format).

- [ ] **Step 2: Add the `precommit check` documentation**

Add a section documenting:
- Purpose: detect config, verify/auto-install tool, verify/init hook, optional run.
- Flags: `--run`, `--no-install`, `--yes`, `--json`.
- Non-TTY behavior: requires `--yes` to mutate the environment.
- Exit codes: `0` ready/skipped, `1` not ready, `2` usage.
- Note `--json` output fields: `config_found`, `tool_installed`, `tool_version`, `hook_installed`, `actions_taken`, `run_result`, `ok`.

- [ ] **Step 3: Add a cheat-sheet row to both skill SKILL.md files**

Add to the command tables: `| \`gc precommit check\` | 提交前检查 pre-commit 配置与本地环境 |`

- [ ] **Step 4: Verify docs build/links (manual scan)**

Confirm no broken markdown tables; the new command name matches `cmd.Use` exactly.

- [ ] **Step 5: Commit**

```bash
git add docs/COMMANDS.md .ai/skills/gitcode-cli/SKILL.md .claude/skills/gitcode-cli/SKILL.md
git commit -m "docs(precommit): document gc precommit check"
```

---

## Task 8: Final verification

- [ ] **Step 1: Build**

Run: `go build -o ./gc.exe ./cmd/gc`
Expected: no errors.

- [ ] **Step 2: Full test + vet**

Run: `go test ./... ; go vet ./...`
Expected: all PASS; vet clean.

- [ ] **Step 3: Manual smoke (this repo has a real .pre-commit-config.yaml)**

Run: `./gc.exe precommit check --no-install`
Expected: prints config found; tool/hook status reflect the local machine; exit code reflects readiness. (Use `--no-install` so the smoke test never mutates the environment.)

Run: `./gc.exe precommit check --json --no-install`
Expected: valid JSON on stdout with the documented fields.

- [ ] **Step 4: Confirm exit codes**

Run (PowerShell): `./gc.exe precommit check --no-install; echo $LASTEXITCODE`
Expected: `0` when ready, `1` when not ready — matching the printed status.

---

## Notes for the implementer

- Follow TDD strictly: write the test, watch it fail, implement, watch it pass, commit.
- Never invoke real `pre-commit`/`python` in unit tests — always go through `CommandRunner` fakes.
- Keep `pkg/precommit/` free of Cobra/IOStreams imports; all I/O lives in the command layer.
- The repository already has a `.pre-commit-config.yaml`, which is useful for the Task 8 manual smoke test (always with `--no-install`).
- Exit-code mapping is handled by `cmdutil.ExitCode` in `cmd/gc/main.go`; return `cmdutil.NewCLIError(...)` to control it.
