package precommit

import (
	"os"
	"path/filepath"
	"strings"
)

// configNames lists accepted pre-commit config filenames, in priority order.
var configNames = []string{".pre-commit-config.yaml", ".pre-commit-config.yml"}

// ConfigFile returns the path to the pre-commit config in root, and whether one exists.
func ConfigFile(root string) (string, bool) {
	for _, name := range configNames {
		p := filepath.Join(root, name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p, true
		}
	}
	return "", false
}

// ToolVersion returns the installed pre-commit version (e.g. "3.7.0") and whether
// the tool is available.
func ToolVersion(r CommandRunner) (string, bool) {
	out, err := r.Run("", "pre-commit", "--version")
	if err != nil {
		return "", false
	}
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) >= 2 {
		return fields[len(fields)-1], true
	}
	if len(fields) == 1 {
		return fields[0], true
	}
	return "", false
}

// HookInstalled reports whether root/.git/hooks/pre-commit exists and was generated
// by the pre-commit tool (detected via a signature marker).
func HookInstalled(root string) bool {
	hook := filepath.Join(root, ".git", "hooks", "pre-commit")
	data, err := os.ReadFile(hook)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "pre-commit")
}
