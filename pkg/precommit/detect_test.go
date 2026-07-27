package precommit

import (
	"os"
	"path/filepath"
	"runtime"
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

func TestToolVersionNameOnly(t *testing.T) {
	r := newFakeRunner()
	// Degenerate output: only the program name, no version token.
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit\n"}
	v, ok := ToolVersion(r)
	if !ok {
		t.Fatal("expected tool to be reported installed")
	}
	if v != "" {
		t.Fatalf("ToolVersion() = %q, want empty version (name must not be used as version)", v)
	}
}

func TestToolVersionIgnoresStderrWarning(t *testing.T) {
	r := newFakeRunner()
	// A warning on stderr must not leak into the parsed version: ToolVersion
	// reads stdout only.
	r.responses[key("pre-commit", "--version")] = fakeResp{
		out:    "pre-commit 3.7.0\n",
		stderr: "WARNING: deprecated flag\n",
	}
	v, ok := ToolVersion(r)
	if !ok || v != "3.7.0" {
		t.Fatalf("ToolVersion() = %q, %v; want 3.7.0, true", v, ok)
	}
}

func TestConfigFileYamlPriority(t *testing.T) {
	dir := t.TempDir()
	yaml := filepath.Join(dir, ".pre-commit-config.yaml")
	yml := filepath.Join(dir, ".pre-commit-config.yml")
	for _, p := range []string{yaml, yml} {
		if err := os.WriteFile(p, []byte("repos: []\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, ok := ConfigFile(dir)
	if !ok || got != yaml {
		t.Fatalf("ConfigFile() = %q, %v; want %q (.yaml takes priority), true", got, ok, yaml)
	}
}

func TestHookInstalled(t *testing.T) {
	root := t.TempDir()
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	if HookInstalled(newFakeRunner(), root) {
		t.Fatal("expected hook not installed")
	}
	hook := filepath.Join(hooks, "pre-commit")
	if err := os.WriteFile(hook, managedHookContent(HookTypePreCommit), 0o755); err != nil {
		t.Fatal(err)
	}
	if !HookInstalled(newFakeRunner(), root) {
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
	if HookInstalled(newFakeRunner(), root) {
		t.Fatal("expected foreign hook to be treated as not pre-commit-installed")
	}
}

func TestHookInstalledResolvesGitPath(t *testing.T) {
	root := t.TempDir()
	// Simulate a worktree where the real hook lives elsewhere.
	realHooks := filepath.Join(t.TempDir(), "commondir", "hooks")
	if err := os.MkdirAll(realHooks, 0o755); err != nil {
		t.Fatal(err)
	}
	hook := filepath.Join(realHooks, "pre-commit")
	if err := os.WriteFile(hook, managedHookContent(HookTypePreCommit), 0o755); err != nil {
		t.Fatal(err)
	}
	r := newFakeRunner()
	r.responses[key("git", "rev-parse", "--git-path", "hooks/pre-commit")] = fakeResp{out: hook + "\n"}
	if !HookInstalled(r, root) {
		t.Fatal("expected hook resolved via git rev-parse to be detected")
	}
}

func TestHookTypeInstalled(t *testing.T) {
	root := t.TempDir()
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hooks, "pre-push"), managedHookContent(HookTypePrePush), 0o755); err != nil {
		t.Fatal(err)
	}

	r := newFakeRunner()
	if HookTypeInstalled(r, root, HookTypePreCommit) {
		t.Fatal("pre-commit hook should be missing")
	}
	if !HookTypeInstalled(r, root, HookTypePrePush) {
		t.Fatal("pre-push hook should be installed")
	}
}

func TestHookTypeInstalledRejectsUnknownType(t *testing.T) {
	if HookTypeInstalled(newFakeRunner(), t.TempDir(), HookType("post-checkout")) {
		t.Fatal("unsupported hook type must not be treated as installed")
	}
}

func TestHookTypeInstalledRejectsMismatchedType(t *testing.T) {
	root := t.TempDir()
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(hooks, "pre-push"),
		managedHookContent(HookTypePreCommit),
		0o755,
	); err != nil {
		t.Fatal(err)
	}
	if HookTypeInstalled(newFakeRunner(), root, HookTypePrePush) {
		t.Fatal("copied pre-commit hook must not satisfy pre-push detection")
	}
}

func TestHookTypeInstalledRequiresExecutableOnPOSIX(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not expose POSIX executable mode semantics")
	}
	root := t.TempDir()
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(hooks, "pre-push"),
		managedHookContent(HookTypePrePush),
		0o644,
	); err != nil {
		t.Fatal(err)
	}
	if HookTypeInstalled(newFakeRunner(), root, HookTypePrePush) {
		t.Fatal("non-executable hook must not be treated as installed")
	}
}

func TestHookTypeInstalledRejectsOtherOnlyExecutableOnPOSIX(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not expose POSIX executable mode semantics")
	}
	root := t.TempDir()
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	hook := filepath.Join(hooks, "pre-push")
	if err := os.WriteFile(hook, managedHookContent(HookTypePrePush), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(hook, 0o001); err != nil {
		t.Fatal(err)
	}
	if HookTypeInstalled(newFakeRunner(), root, HookTypePrePush) {
		t.Fatal("other-only executable hook must not be treated as executable by its owner")
	}
}

func managedHookContent(hookType HookType) []byte {
	return []byte("# File generated by pre-commit\nARGS=(hook-impl --hook-type=" + string(hookType) + ")\n")
}
