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
