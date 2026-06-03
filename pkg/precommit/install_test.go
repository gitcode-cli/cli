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
	calls := 0
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errors.New("not found")}
	r.look["pipx"] = true
	r.responses[key("pipx", "install", "pre-commit")] = fakeResp{out: "installed"}

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

func TestEnsureToolViaPython3(t *testing.T) {
	r := newFakeRunner()
	calls := 0
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errors.New("not found")}
	r.look["python3"] = true
	r.responses[key("python3", "-m", "pip", "install", "--user", "pre-commit")] = fakeResp{out: "ok"}
	wrapped := &versionAfterInstall{fakeRunner: r, succeedAfter: 1, callCount: &calls}
	action, err := EnsureTool(wrapped)
	if err != nil {
		t.Fatalf("EnsureTool() error = %v", err)
	}
	if action == "" {
		t.Fatal("expected non-empty action")
	}
	if !r.called("python3", "-m", "pip", "install", "--user", "pre-commit") {
		t.Fatal("expected python3 pip install to be called")
	}
}

func TestEnsureToolViaPython(t *testing.T) {
	r := newFakeRunner()
	calls := 0
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errors.New("not found")}
	r.look["python"] = true
	r.responses[key("python", "-m", "pip", "install", "--user", "pre-commit")] = fakeResp{out: "ok"}
	wrapped := &versionAfterInstall{fakeRunner: r, succeedAfter: 1, callCount: &calls}
	_, err := EnsureTool(wrapped)
	if err != nil {
		t.Fatalf("EnsureTool() error = %v", err)
	}
	if !r.called("python", "-m", "pip", "install", "--user", "pre-commit") {
		t.Fatal("expected python pip install to be called")
	}
}

func TestEnsureToolStillMissingAfterInstall(t *testing.T) {
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errors.New("not found")}
	r.look["pipx"] = true
	r.responses[key("pipx", "install", "pre-commit")] = fakeResp{out: "ok"}
	// pre-commit --version keeps failing even after install -> EnsureTool must error.
	_, err := EnsureTool(r)
	if err == nil {
		t.Fatal("expected error when tool is still missing after install attempt")
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
