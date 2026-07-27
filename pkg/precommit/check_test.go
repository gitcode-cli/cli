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
	writeInstalledHookType(t, root, HookTypePreCommit)
	writeInstalledHookType(t, root, HookTypePrePush)
}

func writeInstalledHookType(t *testing.T, root string, hookType HookType) {
	t.Helper()
	hooks := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hooks, string(hookType)), managedHookContent(hookType), 0o755); err != nil {
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
	r := newFakeRunner()
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
	if !res.OK || !res.ConfigFound || !res.ToolInstalled || !res.HookInstalled || !res.HooksInstalled ||
		!res.PreCommitHookInstalled || !res.PrePushHookInstalled {
		t.Fatalf("expected fully ready result, got %+v", res)
	}
	if res.ToolVersion != "3.7.0" {
		t.Fatalf("ToolVersion = %q", res.ToolVersion)
	}
}

func TestCheckRequiresBothHookTypes(t *testing.T) {
	tests := []struct {
		name      string
		installed HookType
	}{
		{name: "pre-commit only", installed: HookTypePreCommit},
		{name: "pre-push only", installed: HookTypePrePush},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeConfig(t, root)
			writeInstalledHookType(t, root, tt.installed)
			r := newFakeRunner()
			r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}

			res, err := Check(r, Options{Root: root})
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			if res.OK || res.HooksInstalled || res.Reason != ReasonHookMissing {
				t.Fatalf("expected missing aggregate hook state, got %+v", res)
			}
			if res.HookInstalled != (tt.installed == HookTypePreCommit) {
				t.Fatalf("legacy hook state should match pre-commit, got %+v", res)
			}
			if res.PreCommitHookInstalled != (tt.installed == HookTypePreCommit) {
				t.Fatalf("unexpected pre-commit hook state: %+v", res)
			}
			if res.PrePushHookInstalled != (tt.installed == HookTypePrePush) {
				t.Fatalf("unexpected pre-push hook state: %+v", res)
			}
		})
	}
}

func TestCheckInstallsHookWhenAllowed(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	r.responses[key(
		"pre-commit", "install",
		"--hook-type", "pre-commit",
		"--hook-type", "pre-push",
	)] = fakeResp{out: "ok"}
	hookWriter := &hookInstallingRunner{fakeRunner: r, root: root}

	res, err := Check(hookWriter, Options{Root: root, AllowInstall: true})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !res.HookInstalled || !res.HooksInstalled ||
		!res.PreCommitHookInstalled || !res.PrePushHookInstalled || !res.OK {
		t.Fatalf("expected hook installed + OK, got %+v", res)
	}
	if len(res.ActionsTaken) == 0 {
		t.Fatal("expected an action recorded for hook install")
	}
}

func TestCheckHookInstallFailureReason(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	r.responses[key(
		"pre-commit", "install",
		"--hook-type", "pre-commit",
		"--hook-type", "pre-push",
	)] = fakeResp{out: "permission denied", err: errExit{}}
	partialWriter := &partialHookInstallingRunner{fakeRunner: r, root: root}

	res, err := Check(partialWriter, Options{Root: root, AllowInstall: true})
	if err == nil {
		t.Fatal("expected a hard error when hook installation fails")
	}
	if res.OK || res.Reason != ReasonInstallFailed {
		t.Fatalf("want not-OK + reason=%q, got %+v", ReasonInstallFailed, res)
	}
	if !res.HookInstalled || res.HooksInstalled ||
		!res.PreCommitHookInstalled || res.PrePushHookInstalled {
		t.Fatalf("hook status should reflect the partially installed final state, got %+v", res)
	}
}

func TestCheckHookInstallIncompleteReason(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	r.responses[key(
		"pre-commit", "install",
		"--hook-type", "pre-commit",
		"--hook-type", "pre-push",
	)] = fakeResp{out: "ok"}

	res, err := Check(r, Options{Root: root, AllowInstall: true})
	if err == nil {
		t.Fatal("expected a hard error when hooks remain missing after installation")
	}
	if res.OK || res.Reason != ReasonInstallFailed || res.HooksInstalled {
		t.Fatalf("want incomplete installation classified as %q, got %+v", ReasonInstallFailed, res)
	}
	if len(res.ActionsTaken) != 0 {
		t.Fatalf("unverified installation must not be reported as completed, got %v", res.ActionsTaken)
	}
}

func TestCheckRunFails(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	writeInstalledHook(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	r.responses[key(
		"pre-commit", "run", "--all-files", "--hook-stage", "pre-commit",
	)] = fakeResp{out: "black....Failed", err: errExit{}}
	r.responses[key(
		"pre-commit", "run", "--all-files", "--hook-stage", "pre-push",
	)] = fakeResp{out: "security passed"}

	res, err := Check(r, Options{Root: root, Run: true})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if res.RunResult != "failed" {
		t.Fatalf("RunResult = %q, want failed", res.RunResult)
	}
	if res.PreCommitRunResult != "failed" || res.PrePushRunResult != "passed" {
		t.Fatalf("unexpected stage results: %+v", res)
	}
	if res.OK {
		t.Fatal("failed run should make result not OK")
	}
	if res.RunOutput != "pre-commit:\nblack....Failed" {
		t.Fatalf("RunOutput = %q, want the captured stage failure output", res.RunOutput)
	}
	if !r.called("pre-commit", "run", "--all-files", "--hook-stage", "pre-push") {
		t.Fatalf("pre-push stage should run even after pre-commit fails, calls=%v", r.calls)
	}
}

func TestCheckRunPrePushFails(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	writeInstalledHook(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	r.responses[key(
		"pre-commit", "run", "--all-files", "--hook-stage", "pre-commit",
	)] = fakeResp{out: "checks passed"}
	r.responses[key(
		"pre-commit", "run", "--all-files", "--hook-stage", "pre-push",
	)] = fakeResp{out: "history scan failed", err: errExit{}}

	res, err := Check(r, Options{Root: root, Run: true})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if res.PreCommitRunResult != "passed" || res.PrePushRunResult != "failed" {
		t.Fatalf("unexpected stage results: %+v", res)
	}
	if res.RunOutput != "pre-push:\nhistory scan failed" {
		t.Fatalf("RunOutput = %q, want the captured pre-push failure", res.RunOutput)
	}
}

func TestCheckRunBothStagesFail(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	writeInstalledHook(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	r.responses[key(
		"pre-commit", "run", "--all-files", "--hook-stage", "pre-commit",
	)] = fakeResp{out: "format failed", err: errExit{}}
	r.responses[key(
		"pre-commit", "run", "--all-files", "--hook-stage", "pre-push",
	)] = fakeResp{out: "history failed", err: errExit{}}

	res, err := Check(r, Options{Root: root, Run: true})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	wantOutput := "pre-commit:\nformat failed\n\npre-push:\nhistory failed"
	if res.PreCommitRunResult != "failed" || res.PrePushRunResult != "failed" {
		t.Fatalf("unexpected stage results: %+v", res)
	}
	if res.RunOutput != wantOutput {
		t.Fatalf("RunOutput = %q, want %q", res.RunOutput, wantOutput)
	}
}

func TestCheckRunPasses(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	writeInstalledHook(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	r.responses[key(
		"pre-commit", "run", "--all-files", "--hook-stage", "pre-commit",
	)] = fakeResp{out: "all passed"}
	r.responses[key(
		"pre-commit", "run", "--all-files", "--hook-stage", "pre-push",
	)] = fakeResp{out: "all passed"}
	res, err := Check(r, Options{Root: root, Run: true})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if res.RunResult != "passed" {
		t.Fatalf("RunResult = %q, want passed", res.RunResult)
	}
	if res.PreCommitRunResult != "passed" || res.PrePushRunResult != "passed" {
		t.Fatalf("unexpected stage results: %+v", res)
	}
	if !r.called("pre-commit", "run", "--all-files", "--hook-stage", "pre-commit") ||
		!r.called("pre-commit", "run", "--all-files", "--hook-stage", "pre-push") {
		t.Fatalf("expected both stages to run, calls=%v", r.calls)
	}
	if !res.OK {
		t.Fatal("expected OK after passing run")
	}
}

func TestCheckReasonNoConfig(t *testing.T) {
	res, err := Check(newFakeRunner(), Options{Root: t.TempDir()})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !res.OK || res.Reason != ReasonNoConfig {
		t.Fatalf("want OK + reason=%q, got OK=%v reason=%q", ReasonNoConfig, res.OK, res.Reason)
	}
}

func TestCheckReasonToolMissing(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	res, err := Check(newFakeRunner(), Options{Root: root})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if res.OK || res.Reason != ReasonToolMissing {
		t.Fatalf("want not-OK + reason=%q, got OK=%v reason=%q", ReasonToolMissing, res.OK, res.Reason)
	}
}

func TestCheckReasonConfigInvalid(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.1.0\n"}
	r.responses[key(
		"pre-commit", "validate-config", filepath.Join(root, ".pre-commit-config.yaml"),
	)] = fakeResp{out: "pre-commit version 3.2.0 is required", err: errExit{}}

	res, err := Check(r, Options{Root: root})
	if err == nil {
		t.Fatal("expected invalid configuration to return a hard error")
	}
	if res.OK || res.Reason != ReasonConfigInvalid {
		t.Fatalf("want not-OK + reason=%q, got %+v", ReasonConfigInvalid, res)
	}
}

func TestCheckReasonHookMissing(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	res, err := Check(r, Options{Root: root}) // AllowInstall false: hook stays missing
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if res.OK || res.Reason != ReasonHookMissing {
		t.Fatalf("want not-OK + reason=%q, got OK=%v reason=%q", ReasonHookMissing, res.OK, res.Reason)
	}
}

func TestCheckReasonRunFailed(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	writeInstalledHook(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	r.responses[key(
		"pre-commit", "run", "--all-files", "--hook-stage", "pre-commit",
	)] = fakeResp{out: "boom", err: errExit{}}
	r.responses[key(
		"pre-commit", "run", "--all-files", "--hook-stage", "pre-push",
	)] = fakeResp{out: "ok"}
	res, err := Check(r, Options{Root: root, Run: true})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if res.OK || res.Reason != ReasonRunFailed {
		t.Fatalf("want not-OK + reason=%q, got OK=%v reason=%q", ReasonRunFailed, res.OK, res.Reason)
	}
}

func TestCheckReasonEmptyWhenReady(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	writeInstalledHook(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{out: "pre-commit 3.7.0\n"}
	res, err := Check(r, Options{Root: root})
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !res.OK || res.Reason != "" {
		t.Fatalf("want OK + empty reason, got OK=%v reason=%q", res.OK, res.Reason)
	}
}

func TestCheckReasonInstallFailed(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root)
	r := newFakeRunner()
	r.responses[key("pre-commit", "--version")] = fakeResp{err: errExit{}} // tool missing
	r.look["pipx"] = true
	r.responses[key("pipx", "install", "pre-commit")] = fakeResp{
		err: errExit{},
		out: "ERROR: [Errno 13] Permission denied",
	}
	res, err := Check(r, Options{Root: root, AllowInstall: true})
	if err == nil {
		t.Fatal("expected a hard error when auto-install fails")
	}
	if res.OK || res.Reason != ReasonInstallFailed {
		t.Fatalf("want not-OK + reason=%q, got OK=%v reason=%q", ReasonInstallFailed, res.OK, res.Reason)
	}
	if len(res.InstallFailureCategories) != 1 || res.InstallFailureCategories[0] != "permission" {
		t.Fatalf("want categories=[permission], got %v", res.InstallFailureCategories)
	}
}

type errExit struct{}

func (errExit) Error() string { return "exit status 1" }

// hookInstallingRunner writes both required hook files when installation runs.
type hookInstallingRunner struct {
	*fakeRunner
	root string
}

// partialHookInstallingRunner simulates a hook install that writes pre-commit
// before failing to create pre-push.
type partialHookInstallingRunner struct {
	*fakeRunner
	root string
}

func (h *partialHookInstallingRunner) Run(dir, name string, args ...string) (string, error) {
	out, err := h.fakeRunner.Run(dir, name, args...)
	if name == "pre-commit" && len(args) > 0 && args[0] == "install" {
		writeInstalledHookNoT(h.root, HookTypePreCommit)
	}
	return out, err
}

func (h *hookInstallingRunner) Run(dir, name string, args ...string) (string, error) {
	out, err := h.fakeRunner.Run(dir, name, args...)
	if name == "pre-commit" && len(args) > 0 && args[0] == "install" {
		writeInstalledHookNoT(h.root, HookTypePreCommit)
		writeInstalledHookNoT(h.root, HookTypePrePush)
	}
	return out, err
}

func writeInstalledHookNoT(root string, hookType HookType) {
	hooks := filepath.Join(root, ".git", "hooks")
	_ = os.MkdirAll(hooks, 0o755)
	_ = os.WriteFile(filepath.Join(hooks, string(hookType)), managedHookContent(hookType), 0o755)
}
