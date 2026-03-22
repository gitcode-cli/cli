// Package git provides Git operations for gitcode-cli
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// IsRepo returns true if current directory is a git repository
func IsRepo() bool {
	_, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	return err == nil
}

// RootDir returns the root directory of the git repository
func RootDir() (string, error) {
	output, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository")
	}
	return strings.TrimSpace(string(output)), nil
}

// CurrentBranch returns the current branch name
func CurrentBranch() (string, error) {
	output, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// HasLocalChanges returns true if there are uncommitted changes
func HasLocalChanges() (bool, error) {
	output, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return false, err
	}
	return len(output) > 0, nil
}

// Remotes returns the list of remote names
func Remotes() ([]string, error) {
	output, err := exec.Command("git", "remote").Output()
	if err != nil {
		return nil, err
	}
	remotes := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(remotes) == 1 && remotes[0] == "" {
		return []string{}, nil
	}
	return remotes, nil
}

// RemoteURL returns the URL of a remote
func RemoteURL(name string) (string, error) {
	output, err := exec.Command("git", "remote", "get-url", name).Output()
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// DefaultRemote returns the default remote (origin or first available)
func DefaultRemote() (string, error) {
	remotes, err := Remotes()
	if err != nil {
		return "", err
	}
	if len(remotes) == 0 {
		return "", fmt.Errorf("no remotes found")
	}
	// Prefer origin
	for _, r := range remotes {
		if r == "origin" {
			return r, nil
		}
	}
	return remotes[0], nil
}

// Run executes a git command and returns the output
func Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, string(output))
	}
	return string(output), nil
}

// RunInDir executes a git command in a specific directory
func RunInDir(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, string(output))
	}
	return string(output), nil
}