package precommit

import (
	"errors"
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
	// ReasonHookMissing: the git pre-commit hook is not initialized.
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
	// the git hook) when something is missing.
	AllowInstall bool
	// Run, when true, executes `pre-commit run --all-files` after the environment
	// is confirmed ready.
	Run bool
}

// Result is the structured outcome of a Check. JSON tags match docs/COMMANDS.md.
//
// OK means "nothing blocks committing": it is true when the repository has no
// pre-commit config (nothing to check), or when the config is present and the
// tool, hook, and any requested --run all succeeded. Always read OK together
// with ConfigFound to distinguish "ready" from "no config, skipped".
type Result struct {
	ConfigFound   bool     `json:"config_found"`
	ToolInstalled bool     `json:"tool_installed"`
	ToolVersion   string   `json:"tool_version,omitempty"`
	HookInstalled bool     `json:"hook_installed"`
	ActionsTaken  []string `json:"actions_taken"`
	RunResult     string   `json:"run_result,omitempty"` // "passed" | "failed" | ""
	RunOutput     string   `json:"run_output,omitempty"` // pre-commit run output when RunResult == "failed"
	OK            bool     `json:"ok"`
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
	if _, found := ConfigFile(opts.Root); !found {
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

	// 3. Hook detection + optional install.
	hookOK := HookInstalled(r, opts.Root)
	if !hookOK && opts.AllowInstall {
		action, err := InstallHook(r, opts.Root)
		if err != nil {
			return res, err
		}
		if action != "" {
			res.ActionsTaken = append(res.ActionsTaken, action)
		}
		hookOK = HookInstalled(r, opts.Root)
	}
	res.HookInstalled = hookOK
	if !hookOK {
		res.Reason = ReasonHookMissing
		return res, nil
	}

	// 4. Environment is ready.
	res.OK = true

	// 5. Optional: actually run the checks. Capture the output so a failure
	// surfaces why, instead of only reporting "failed".
	if opts.Run {
		out, err := r.Run(opts.Root, "pre-commit", "run", "--all-files")
		if err != nil {
			res.RunResult = "failed"
			res.RunOutput = strings.TrimSpace(out)
			res.OK = false
			res.Reason = ReasonRunFailed
		} else {
			res.RunResult = "passed"
		}
	}

	return res, nil
}
