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
	for _, ins := range installers() {
		if !r.Look(ins.name) {
			continue
		}
		if out, err := r.Run("", ins.name, ins.args...); err != nil {
			attempts = append(attempts, fmt.Sprintf("%s: %v: %s", ins.desc, err, strings.TrimSpace(out)))
			continue
		}
		if _, ok := ToolVersion(r); !ok {
			attempts = append(attempts, fmt.Sprintf("%s ran but pre-commit was still not found on PATH", ins.desc))
			continue
		}
		return "installed pre-commit via " + ins.desc, nil
	}

	if len(attempts) == 0 {
		return "", fmt.Errorf("cannot auto-install pre-commit: no pipx/python3/python found. %s", manualInstallHint())
	}
	return "", fmt.Errorf("failed to auto-install pre-commit (%s). %s", strings.Join(attempts, "; "), manualInstallHint())
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
