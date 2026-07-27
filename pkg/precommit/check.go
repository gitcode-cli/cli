package precommit

import (
	"errors"
	"fmt"
	"strings"
)

// Reason values are stable, machine-readable classifications of a Check outcome,
// emitted in Result.Reason so scripts/agents can branch without parsing prose.
// Reason is empty when the environment is fully ready (and any requested run
// passed). The values are part of the --json contract; see docs/COMMANDS.md.
const (
	// ReasonNoConfig: the repository has no pre-commit config; nothing to check.
	// Paired with OK == true (a clean skip, not a failure).
	ReasonNoConfig = "no_config"
	// ReasonToolMissing: the pre-commit tool is not installed (and was not, or
	// could not be, installed).
	ReasonToolMissing = "tool_missing"
	// ReasonConfigInvalid: pre-commit rejected the repository configuration.
	ReasonConfigInvalid = "config_invalid"
	// ReasonHookMissing: a required git hook is not initialized.
	ReasonHookMissing = "hook_missing"
	// ReasonRunFailed: the environment is ready but `pre-commit run` failed.
	ReasonRunFailed = "run_failed"
	// ReasonInstallFailed: an auto-install attempt was made but failed. The
	// machine-readable failure categories are carried in
	// Result.InstallFailureCategories. Paired with a non-nil Check error.
	ReasonInstallFailed = "install_failed"
	// ReasonNotInRepo: the working directory is not inside a git repository.
	// Set by the command layer, which never reaches Check.
	ReasonNotInRepo = "not_in_repo"
)

// Options controls a Check run.
type Options struct {
	// Root is the git repository root directory.
	Root string
	// AllowInstall permits mutating the environment (installing the tool and/or
	// the required git hooks) when something is missing.
	AllowInstall bool
	// Run, when true, executes every required pre-commit hook stage across all
	// files after the environment is confirmed ready.
	Run bool
}

// Result is the structured outcome of a Check. JSON tags match docs/COMMANDS.md.
//
// OK means "nothing blocks committing or pushing": it is true when the
// repository has no pre-commit config (nothing to check), or when the config is
// present and the tool, hooks, and any requested --run all succeeded. Always
// read OK together with ConfigFound to distinguish "ready" from "no config,
// skipped".
type Result struct {
	ConfigFound   bool   `json:"config_found"`
	ToolInstalled bool   `json:"tool_installed"`
	ToolVersion   string `json:"tool_version,omitempty"`
	// HookInstalled preserves the original pre-commit-specific JSON contract.
	HookInstalled          bool     `json:"hook_installed"`
	HooksInstalled         bool     `json:"hooks_installed"`
	PreCommitHookInstalled bool     `json:"pre_commit_hook_installed"`
	PrePushHookInstalled   bool     `json:"pre_push_hook_installed"`
	ActionsTaken           []string `json:"actions_taken"`
	RunResult              string   `json:"run_result,omitempty"` // aggregate: "passed" | "failed" | ""
	RunOutput              string   `json:"run_output,omitempty"` // failed stage output
	PreCommitRunResult     string   `json:"pre_commit_run_result,omitempty"`
	PrePushRunResult       string   `json:"pre_push_run_result,omitempty"`
	OK                     bool     `json:"ok"`
	// Reason is a stable, machine-readable classification of the outcome (one of
	// the Reason* constants), or "" when the environment is fully ready.
	Reason string `json:"reason,omitempty"`
	// InstallFailureCategories carries the distinct auto-install failure
	// categories ("permission" | "network" | "toolchain"), in first-seen order,
	// when Reason == ReasonInstallFailed. Empty otherwise.
	InstallFailureCategories []string `json:"install_failure_categories,omitempty"`
}

// Check runs the detection/remediation pipeline and returns a structured Result.
// It returns a non-nil error only for hard failures (e.g. an install attempt
// failed). "Not ready" states are reported via Result.OK == false with no error.
func Check(r CommandRunner, opts Options) (Result, error) {
	res := Result{ActionsTaken: []string{}}

	// 1. Config detection — absence is a clean skip, not an error.
	configPath, found := ConfigFile(opts.Root)
	if !found {
		res.OK = true
		res.Reason = ReasonNoConfig
		return res, nil
	}
	res.ConfigFound = true

	// 2. Tool detection + optional install.
	version, ok := ToolVersion(r)
	if !ok && opts.AllowInstall {
		action, err := EnsureTool(r)
		if err != nil {
			// A hard install failure: classify it so --json consumers get a
			// machine-readable reason and categories, not just stderr prose.
			res.Reason = ReasonInstallFailed
			var ie *InstallError
			if errors.As(err, &ie) {
				res.InstallFailureCategories = ie.CategoryNames()
			}
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
		res.Reason = ReasonToolMissing
		return res, nil // not ready; OK stays false
	}

	// 3. Configuration validation.
	if out, err := r.Run(opts.Root, "pre-commit", "validate-config", configPath); err != nil {
		res.Reason = ReasonConfigInvalid
		return res, fmt.Errorf("pre-commit config validation failed: %w: %s", err, strings.TrimSpace(out))
	}

	// 4. Hook detection + optional install.
	setHookStatus(r, opts.Root, &res)
	if !res.HooksInstalled && opts.AllowInstall {
		action, err := InstallHook(r, opts.Root)
		setHookStatus(r, opts.Root, &res)
		if err != nil {
			res.Reason = ReasonInstallFailed
			return res, err
		}
		if !res.HooksInstalled {
			res.Reason = ReasonInstallFailed
			return res, errors.New("pre-commit install completed but required hooks are still missing")
		}
		if action != "" {
			res.ActionsTaken = append(res.ActionsTaken, action)
		}
	}
	if !res.HooksInstalled {
		res.Reason = ReasonHookMissing
		return res, nil
	}

	// 5. Environment is ready.
	res.OK = true

	// 6. Optional: run every configured git-hook stage.
	if opts.Run {
		outputs, runErrors := runRequiredHookStages(r, opts.Root)
		res.PreCommitRunResult = runResult(runErrors[HookTypePreCommit])
		res.PrePushRunResult = runResult(runErrors[HookTypePrePush])
		if hasRunFailure(runErrors) {
			res.RunResult = "failed"
			res.RunOutput = failedRunOutput(outputs, runErrors)
			res.OK = false
			res.Reason = ReasonRunFailed
		} else {
			res.RunResult = "passed"
		}
	}

	return res, nil
}

func setHookStatus(r CommandRunner, root string, res *Result) {
	statuses := make(map[HookType]bool, len(requiredHookTypes))
	res.HooksInstalled = true
	for _, hookType := range requiredHookTypes {
		statuses[hookType] = HookTypeInstalled(r, root, hookType)
		res.HooksInstalled = res.HooksInstalled && statuses[hookType]
	}
	res.PreCommitHookInstalled = statuses[HookTypePreCommit]
	res.PrePushHookInstalled = statuses[HookTypePrePush]
	res.HookInstalled = res.PreCommitHookInstalled
}

func runHookStage(r CommandRunner, root string, hookType HookType) (string, error) {
	return r.Run(
		root,
		"pre-commit",
		"run", "--all-files", "--hook-stage", string(hookType),
	)
}

func runRequiredHookStages(r CommandRunner, root string) (map[HookType]string, map[HookType]error) {
	outputs := make(map[HookType]string, len(requiredHookTypes))
	runErrors := make(map[HookType]error, len(requiredHookTypes))
	for _, hookType := range requiredHookTypes {
		outputs[hookType], runErrors[hookType] = runHookStage(r, root, hookType)
	}
	return outputs, runErrors
}

func hasRunFailure(runErrors map[HookType]error) bool {
	for _, hookType := range requiredHookTypes {
		if runErrors[hookType] != nil {
			return true
		}
	}
	return false
}

func runResult(err error) string {
	if err != nil {
		return "failed"
	}
	return "passed"
}

func failedRunOutput(outputs map[HookType]string, runErrors map[HookType]error) string {
	var failures []string
	for _, hookType := range requiredHookTypes {
		if runErrors[hookType] != nil {
			failures = append(failures, string(hookType)+":\n"+strings.TrimSpace(outputs[hookType]))
		}
	}
	return strings.Join(failures, "\n\n")
}
