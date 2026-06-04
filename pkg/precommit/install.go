package precommit

import (
	"fmt"
	"runtime"
	"strings"
)

// installer describes one way to install pre-commit.
type installer struct {
	name string   // executable to look up on PATH
	args []string // arguments to run it with
	desc string   // human-readable description for the action/error message
}

// installers returns the ordered list of install strategies to try for the
// current platform. Each is attempted in turn; a failure falls through to the
// next rather than aborting, so e.g. a failed pipx install still tries pip.
func installers() []installer {
	list := []installer{
		{"pipx", []string{"install", "pre-commit"}, "pipx"},
		{"python3", []string{"-m", "pip", "install", "--user", "pre-commit"}, "python3 -m pip --user"},
		{"python", []string{"-m", "pip", "install", "--user", "pre-commit"}, "python -m pip --user"},
	}
	if runtime.GOOS == "windows" {
		// The Windows Python launcher is often present when "python"/"python3"
		// are not on PATH.
		list = append(list, installer{"py", []string{"-m", "pip", "install", "--user", "pre-commit"}, "py -m pip --user"})
	}
	return list
}

// EnsureTool makes sure the pre-commit tool is available, installing it when
// necessary. It tries each available installer in order, falling through on
// failure. It returns a human-readable description of any action taken (empty
// if the tool was already present) and an error if installation was impossible
// or every attempt failed.
func EnsureTool(r CommandRunner) (string, error) {
	if _, ok := ToolVersion(r); ok {
		return "", nil
	}

	var attempts []string
	cats := newCategorySet()
	for _, ins := range installers() {
		if !r.Look(ins.name) {
			continue
		}
		if out, err := r.Run("", ins.name, ins.args...); err != nil {
			attempts = append(attempts, fmt.Sprintf("%s: %v: %s", ins.desc, err, strings.TrimSpace(out)))
			cats.add(classifyInstallFailure(err, out))
			continue
		}
		if _, ok := ToolVersion(r); !ok {
			attempts = append(attempts, fmt.Sprintf("%s ran but pre-commit was still not found on PATH", ins.desc))
			continue
		}
		return "installed pre-commit via " + ins.desc, nil
	}

	if len(attempts) == 0 {
		// No usable installer at all: the host lacks a Python/pip toolchain.
		return "", &InstallError{
			msg:        fmt.Sprintf("cannot auto-install pre-commit: no pipx/python3/python found. %s", manualInstallHint()),
			categories: []installFailureCategory{failToolchain},
		}
	}
	msg := fmt.Sprintf("failed to auto-install pre-commit (%s).", strings.Join(attempts, "; "))
	if guidance := cats.guidance(); guidance != "" {
		msg += " " + guidance
	}
	return "", &InstallError{
		msg:        fmt.Sprintf("%s %s", msg, manualInstallHint()),
		categories: cats.order,
	}
}

// InstallError reports a failure to auto-install pre-commit. Beyond the
// human-readable message it carries the distinct failure categories (in
// first-seen order) so callers can surface them in structured output (e.g.
// --json) instead of only emitting prose to stderr.
type InstallError struct {
	msg        string
	categories []installFailureCategory
}

func (e *InstallError) Error() string { return e.msg }

// CategoryNames returns the failure categories as stable, machine-readable
// identifiers ("permission" | "network" | "toolchain"), in first-seen order.
// Unclassified failures are omitted, so the slice may be empty.
func (e *InstallError) CategoryNames() []string {
	names := make([]string, 0, len(e.categories))
	for _, c := range e.categories {
		if n := c.name(); n != "" {
			names = append(names, n)
		}
	}
	return names
}

// installFailureCategory classifies an install failure so the aggregated error
// can carry targeted remediation guidance instead of only raw tool output.
type installFailureCategory int

const (
	failOther installFailureCategory = iota
	failPermission
	failNetwork
	failToolchain
)

// classifyInstallFailure infers the failure category from the installer's error
// and output. It is best-effort and intentionally string-based: the underlying
// installers (pip/pipx) report these conditions in their messages, and matching
// keeps the logic testable through the fake CommandRunner.
func classifyInstallFailure(err error, output string) installFailureCategory {
	s := strings.ToLower(output)
	if err != nil {
		s += " " + strings.ToLower(err.Error())
	}
	switch {
	case containsAny(s, "permission denied", "access is denied", "errno 13", "not permitted", "operation not permitted"):
		return failPermission
	case containsAny(s,
		"could not resolve", "temporary failure in name resolution", "network is unreachable",
		"connection refused", "connection reset", "connection timed out", "timed out", "timeout",
		"max retries exceeded", "failed to establish a new connection", "proxy", "ssl"):
		return failNetwork
	case containsAny(s, "no module named", "command not found", "not recognized", "no such file"):
		return failToolchain
	default:
		return failOther
	}
}

// name is the stable, machine-readable identifier for a category, emitted in
// the --json contract (Result.InstallFailureCategories). failOther has no name.
func (c installFailureCategory) name() string {
	switch c {
	case failPermission:
		return "permission"
	case failNetwork:
		return "network"
	case failToolchain:
		return "toolchain"
	default:
		return ""
	}
}

func (c installFailureCategory) hint() string {
	switch c {
	case failPermission:
		return "Permission denied: avoid a privileged global install — prefer pipx, or install into a --user/virtualenv location you own."
	case failNetwork:
		return "Network failure: check your connection/proxy and retry, or install pre-commit from an offline package."
	case failToolchain:
		return "Toolchain missing: pip/python was not usable — install Python (with pip) first, then retry."
	default:
		return ""
	}
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// categorySet collects distinct failure categories in first-seen order so the
// guidance is deterministic and free of duplicates when several installers fail.
type categorySet struct {
	seen  map[installFailureCategory]bool
	order []installFailureCategory
}

func newCategorySet() *categorySet {
	return &categorySet{seen: map[installFailureCategory]bool{}}
}

func (cs *categorySet) add(c installFailureCategory) {
	if c == failOther || cs.seen[c] {
		return
	}
	cs.seen[c] = true
	cs.order = append(cs.order, c)
}

func (cs *categorySet) guidance() string {
	var hints []string
	for _, c := range cs.order {
		if h := c.hint(); h != "" {
			hints = append(hints, h)
		}
	}
	return strings.Join(hints, " ")
}

// InstallHook runs `pre-commit install` in root to set up the git hook.
func InstallHook(r CommandRunner, root string) (string, error) {
	if out, err := r.Run(root, "pre-commit", "install"); err != nil {
		return "", fmt.Errorf("pre-commit install failed: %w: %s", err, strings.TrimSpace(out))
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
