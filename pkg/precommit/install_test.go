package precommit

import (
	"errors"
	"strings"
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

func TestEnsureToolPipxFailsFallsBackToPython3(t *testing.T) {
	r := newFakeRunner()
	calls := 0
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errors.New("not found")}
	r.look["pipx"] = true
	r.look["python3"] = true
	// pipx install fails; the cascade must fall through to python3.
	r.responses[key("pipx", "install", "pre-commit")] = fakeResp{err: errors.New("pipx boom"), out: "boom"}
	r.responses[key("python3", "-m", "pip", "install", "--user", "pre-commit")] = fakeResp{out: "ok"}

	wrapped := &versionAfterInstall{fakeRunner: r, succeedAfter: 1, callCount: &calls}
	action, err := EnsureTool(wrapped)
	if err != nil {
		t.Fatalf("EnsureTool() error = %v", err)
	}
	if !r.called("pipx", "install", "pre-commit") {
		t.Fatal("expected pipx install to be attempted first")
	}
	if !r.called("python3", "-m", "pip", "install", "--user", "pre-commit") {
		t.Fatal("expected fallthrough to python3 install after pipx failed")
	}
	if action == "" {
		t.Fatal("expected non-empty action after successful fallback")
	}
}

func TestClassifyInstallFailure(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		output string
		want   installFailureCategory
	}{
		{"permission", errors.New("exit status 1"), "ERROR: Could not install packages due to an OSError: [Errno 13] Permission denied", failPermission},
		{"windows permission", errors.New("exit 1"), "Access is denied", failPermission},
		{"network dns", errors.New("exit 1"), "Could not resolve host: pypi.org", failNetwork},
		{"network retries", errors.New("exit 1"), "Max retries exceeded with url", failNetwork},
		{"toolchain", errors.New("exit 1"), "No module named pip", failToolchain},
		{"other", errors.New("exit 1"), "some unrelated failure", failOther},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := classifyInstallFailure(tc.err, tc.output); got != tc.want {
				t.Fatalf("classifyInstallFailure(%q) = %v, want %v", tc.output, got, tc.want)
			}
		})
	}
}

func TestEnsureToolPermissionGuidance(t *testing.T) {
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errors.New("not found")}
	r.look["pipx"] = true
	r.responses[key("pipx", "install", "pre-commit")] = fakeResp{
		err: errors.New("exit status 1"),
		out: "ERROR: [Errno 13] Permission denied",
	}
	_, err := EnsureTool(r)
	if err == nil {
		t.Fatal("expected error when install fails")
	}
	if !strings.Contains(err.Error(), "Permission denied:") {
		t.Fatalf("expected targeted permission guidance, got %q", err.Error())
	}
}

func TestEnsureToolNetworkGuidance(t *testing.T) {
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errors.New("not found")}
	r.look["pipx"] = true
	r.responses[key("pipx", "install", "pre-commit")] = fakeResp{
		err: errors.New("exit status 1"),
		out: "Could not resolve host: pypi.org",
	}
	_, err := EnsureTool(r)
	if err == nil {
		t.Fatal("expected error when install fails")
	}
	if !strings.Contains(err.Error(), "Network failure:") {
		t.Fatalf("expected targeted network guidance, got %q", err.Error())
	}
}

func TestEnsureToolDedupesGuidance(t *testing.T) {
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errors.New("not found")}
	r.look["pipx"] = true
	r.look["python3"] = true
	// Both installers fail with the same category: guidance must appear once.
	r.responses[key("pipx", "install", "pre-commit")] = fakeResp{err: errors.New("e"), out: "Permission denied"}
	r.responses[key("python3", "-m", "pip", "install", "--user", "pre-commit")] = fakeResp{err: errors.New("e"), out: "Access is denied"}
	_, err := EnsureTool(r)
	if err == nil {
		t.Fatal("expected error")
	}
	if n := strings.Count(err.Error(), "Permission denied:"); n != 1 {
		t.Fatalf("expected permission guidance exactly once, got %d in %q", n, err.Error())
	}
}

// versionAfterInstall makes pre-commit --version fail until succeedAfter calls have
// happened, simulating a tool that appears after installation.
type versionAfterInstall struct {
	*fakeRunner
	succeedAfter int
	callCount    *int
}

func (w *versionAfterInstall) RunStdout(dir, name string, args ...string) (string, error) {
	if name == "pre-commit" && len(args) == 1 && args[0] == "--version" {
		*w.callCount++
		if *w.callCount > w.succeedAfter {
			return "pre-commit 3.7.0\n", nil
		}
		return "", errors.New("not found")
	}
	return w.fakeRunner.RunStdout(dir, name, args...)
}
