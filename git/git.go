// Package git provides Git operations for gitcode-cli
package git

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"regexp"
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
	return runWithEnv("", nil, args...)
}

// RunInDir executes a git command in a specific directory
func RunInDir(dir string, args ...string) (string, error) {
	return runWithEnv(dir, nil, args...)
}

// RunWithEnv executes a git command with extra environment variables.
func RunWithEnv(env map[string]string, args ...string) (string, error) {
	return runWithEnv("", env, args...)
}

// RunInDirWithEnv executes a git command in a specific directory with extra environment variables.
func RunInDirWithEnv(dir string, env map[string]string, args ...string) (string, error) {
	return runWithEnv(dir, env, args...)
}

func runWithEnv(dir string, env map[string]string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if len(env) > 0 {
		cmd.Env = os.Environ()
		for key, value := range env {
			cmd.Env = append(cmd.Env, key+"="+value)
		}
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, string(output))
	}
	return string(output), nil
}

// refPattern matches valid git ref names. Ref names must not:
// - Start with "-" (would be interpreted as a git option)
// - Contain control characters, spaces, or shell metacharacters
var refPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._/\-]*$`)

// ErrInvalidRef is returned when a git ref fails validation.
var ErrInvalidRef = fmt.Errorf("invalid git ref")

// ValidateRef validates that a git ref (branch, tag, or remote ref) is safe
// for use as a git command argument. It rejects:
//   - Empty strings
//   - Strings starting with "-" (prevent option injection)
//   - Strings with control characters, spaces, or shell metacharacters
//   - Strings that don't match the expected ref pattern
func ValidateRef(ref string) error {
	if ref == "" {
		return fmt.Errorf("%w: ref must not be empty", ErrInvalidRef)
	}
	if strings.HasPrefix(ref, "-") {
		return fmt.Errorf("%w: ref must not start with '-': %q", ErrInvalidRef, ref)
	}
	if !refPattern.MatchString(ref) {
		return fmt.Errorf("%w: ref contains invalid characters: %q", ErrInvalidRef, ref)
	}
	return nil
}

// ValidateFetchURL validates a git fetch/push URL. It rejects:
//   - Empty strings
//   - URLs starting with "-" (prevent option injection via dash-prefixed host)
//   - URLs with invalid or unexpected schemes
func ValidateFetchURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("%w: URL must not be empty", ErrInvalidRef)
	}
	if strings.HasPrefix(rawURL, "-") {
		return fmt.Errorf("%w: URL must not start with '-': %q", ErrInvalidRef, rawURL)
	}

	// For SSH-style URLs (git@host:path), validate the host portion
	if strings.Contains(rawURL, "@") && !strings.HasPrefix(rawURL, "http") {
		parts := strings.SplitN(rawURL, ":", 2)
		if len(parts) == 2 {
			hostPart := parts[0]
			if idx := strings.LastIndex(hostPart, "@"); idx >= 0 {
				host := hostPart[idx+1:]
				if strings.HasPrefix(host, "-") {
					return fmt.Errorf("%w: host must not start with '-': %q", ErrInvalidRef, host)
				}
			}
		}
		return nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("%w: invalid URL: %w", ErrInvalidRef, err)
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return fmt.Errorf("%w: unsupported URL scheme: %q", ErrInvalidRef, u.Scheme)
	}
	if strings.HasPrefix(u.Host, "-") {
		return fmt.Errorf("%w: host must not start with '-': %q", ErrInvalidRef, u.Host)
	}
	return nil
}

// SafeCheckout runs "git checkout <branch>" after validating the branch name.
// This prevents option-injection attacks where a branch name starts with "-".
func SafeCheckout(branch string) error {
	if err := ValidateRef(branch); err != nil {
		return err
	}
	_, err := runWithEnv("", nil, "checkout", branch)
	return err
}

// SafeCheckoutWithOutput runs "git checkout <branch>" after validating the
// branch name, streaming stdout/stderr to the provided writers.
func SafeCheckoutWithOutput(stdout, stderr io.Writer, dir string, branch string) error {
	if err := ValidateRef(branch); err != nil {
		return err
	}
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}

// SafeFetch runs "git fetch <remote> <ref>:<localBranch>" after validating
// both the remote ref and the local branch name. All parameters including
// "remote" are validated to prevent option-injection.
func SafeFetch(remote, ref, localBranch string) error {
	if err := ValidateRef(remote); err != nil {
		return fmt.Errorf("invalid remote name for fetch: %w", err)
	}
	if err := ValidateRef(localBranch); err != nil {
		return fmt.Errorf("invalid local branch for fetch: %w", err)
	}
	if err := ValidateRef(ref); err != nil {
		return fmt.Errorf("invalid remote ref for fetch: %w", err)
	}
	refspec := ref + ":" + localBranch
	_, err := runWithEnv("", nil, "fetch", remote, refspec)
	return err
}

// SafeFetchWithOutput runs a validated git fetch, streaming output.
// All parameters including "remote" are validated to prevent option-injection.
func SafeFetchWithOutput(stdout, stderr io.Writer, dir, remote, ref, localBranch string) error {
	if err := ValidateRef(remote); err != nil {
		return fmt.Errorf("invalid remote name for fetch: %w", err)
	}
	if err := ValidateRef(localBranch); err != nil {
		return fmt.Errorf("invalid local branch for fetch: %w", err)
	}
	if err := ValidateRef(ref); err != nil {
		return fmt.Errorf("invalid remote ref for fetch: %w", err)
	}
	refspec := ref + ":" + localBranch
	cmd := exec.Command("git", "fetch", remote, refspec)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}

// SafeFetchFromURL runs "git fetch <fetchURL> <ref>:<localBranch>" after
// validating both the URL and the ref names. This is used when fetching from
// a fork or external repository.
func SafeFetchFromURL(fetchURL, ref, localBranch string) error {
	if err := ValidateFetchURL(fetchURL); err != nil {
		return fmt.Errorf("invalid fetch URL: %w", err)
	}
	if err := ValidateRef(localBranch); err != nil {
		return fmt.Errorf("invalid local branch for fetch: %w", err)
	}
	if err := ValidateRef(ref); err != nil {
		return fmt.Errorf("invalid remote ref for fetch: %w", err)
	}
	refspec := ref + ":" + localBranch
	_, err := runWithEnv("", nil, "fetch", fetchURL, refspec)
	return err
}

// SafeFetchFromURLWithOutput runs a validated git fetch from a URL, streaming output.
func SafeFetchFromURLWithOutput(stdout, stderr io.Writer, dir, fetchURL, ref, localBranch string) error {
	if err := ValidateFetchURL(fetchURL); err != nil {
		return fmt.Errorf("invalid fetch URL: %w", err)
	}
	if err := ValidateRef(localBranch); err != nil {
		return fmt.Errorf("invalid local branch for fetch: %w", err)
	}
	if err := ValidateRef(ref); err != nil {
		return fmt.Errorf("invalid remote ref for fetch: %w", err)
	}
	refspec := ref + ":" + localBranch
	cmd := exec.Command("git", "fetch", fetchURL, refspec)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}
