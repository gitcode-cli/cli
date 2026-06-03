package precommit

import (
	"fmt"
	"runtime"
	"strings"
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
			return "", fmt.Errorf("pipx install pre-commit failed: %w: %s", err, strings.TrimSpace(out))
		}
		action = "installed pre-commit via pipx"
	case r.Look("python3"):
		if out, err := r.Run("", "python3", "-m", "pip", "install", "--user", "pre-commit"); err != nil {
			return "", fmt.Errorf("python3 -m pip install pre-commit failed: %w: %s", err, strings.TrimSpace(out))
		}
		action = "installed pre-commit via python3 -m pip --user"
	case r.Look("python"):
		if out, err := r.Run("", "python", "-m", "pip", "install", "--user", "pre-commit"); err != nil {
			return "", fmt.Errorf("python -m pip install pre-commit failed: %w: %s", err, strings.TrimSpace(out))
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
